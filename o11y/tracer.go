package o11y

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	tracelib "go.opentelemetry.io/otel/trace"
)

var _ trace.SpanProcessor = new(spanProcessor)

type spanProcessor struct {
	log         logger.Log
	spans       map[tracelib.TraceID]map[tracelib.SpanID]trace.ReadOnlySpan
	shuttedDown bool
}

func newSpanProcessor(log logger.Log) *spanProcessor {
	return &spanProcessor{
		log:   log,
		spans: make(map[tracelib.TraceID]map[tracelib.SpanID]trace.ReadOnlySpan),
	}
}

func (s *spanProcessor) OnStart(_ context.Context, span trace.ReadWriteSpan) {
	if s.shuttedDown {
		return
	}

	traceID := span.SpanContext().TraceID()
	if _, exists := s.spans[traceID]; !exists {
		s.spans[traceID] = make(map[tracelib.SpanID]trace.ReadOnlySpan)
	}
}

func (s *spanProcessor) OnEnd(span trace.ReadOnlySpan) {
	if s.shuttedDown {
		return
	}

	traceID := span.SpanContext().TraceID()
	spanID := span.SpanContext().SpanID()
	s.spans[traceID][spanID] = span

	if isRoot := !span.Parent().SpanID().IsValid(); !isRoot {
		return
	}

	s.print(span.SpanContext().TraceID())
	delete(s.spans, traceID)
}

func (s *spanProcessor) Shutdown(_ context.Context) error {
	s.shuttedDown = true
	s.spans = nil

	return nil
}

func (s *spanProcessor) ForceFlush(ctx context.Context) error {
	return s.Shutdown(ctx)
}

type spanTree map[tracelib.SpanID]*spanNode

type spanNode struct {
	trace.ReadOnlySpan
	children children
}

var _ sort.Interface = children{}

type children []trace.ReadOnlySpan

func (c children) Len() int {
	return len(c)
}

func (c children) Less(i, j int) bool {
	return c[i].StartTime().Before(c[j].StartTime())
}

func (c children) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (s *spanProcessor) print(traceID tracelib.TraceID) {
	var (
		rootSpanID           tracelib.SpanID
		traceStart, traceEnd time.Time
	)

	tree := make(spanTree)
	for _, span := range s.spans[traceID] {
		spanID := span.SpanContext().SpanID()

		node, exists := tree[spanID]
		if !exists {
			node = &spanNode{}
		}
		node.ReadOnlySpan = span
		tree[spanID] = node

		if isRoot := !span.Parent().SpanID().IsValid(); isRoot {
			rootSpanID = spanID
			traceStart, traceEnd = span.StartTime(), span.EndTime()

			continue
		}

		parentID := span.Parent().SpanID()
		parent, exists := tree[parentID]
		if !exists {
			parent = &spanNode{}
		}
		parent.children = append(parent.children, span)
		tree[parentID] = parent
	}

	s.printSpan(tree, rootSpanID, traceStart, traceEnd)
}

func (s *spanProcessor) printSpan(tree spanTree, id tracelib.SpanID, traceStart, traceEnd time.Time) {
	traceDuration := traceEnd.Sub(traceStart).Nanoseconds()
	durToBins := func(d time.Duration) int {
		bins := float64(d.Nanoseconds()) / float64(traceDuration) * 10
		return int(math.Max(1, math.Round(bins)))
	}

	node := tree[id]
	ascii := []string{
		strings.Repeat(" ", durToBins(node.StartTime().Sub(traceStart))),
		strings.Repeat("-", durToBins(node.EndTime().Sub(node.StartTime()))),
		strings.Repeat(" ", durToBins(traceEnd.Sub(node.EndTime()))),
		" ",
		node.Name(),
	}

	logAttrs := []interface{}{
		strings.Join(ascii, ""),
		kv.New("duration", node.EndTime().Sub(node.StartTime())),
	}
	if status := node.Status(); status.Code == codes.Error {
		logAttrs[0] = errors.New(logAttrs[0].(string))
		logAttrs = append(logAttrs, kv.New("error", status.Description))
	}

	appendAttributes := func(attrs []attribute.KeyValue) {
		for _, attr := range attrs {
			var val interface{}

			switch attr.Value.Type() {
			case attribute.BOOL:
				val = attr.Value.AsBool()

			case attribute.INT64:
				val = attr.Value.AsInt64()

			case attribute.FLOAT64:
				val = attr.Value.AsFloat64()

			case attribute.STRING:
				val = attr.Value.AsString()

			default:
				continue
			}

			logAttrs = append(logAttrs, kv.New(string(attr.Key), val))
		}
	}

	appendEvents := func(evs []trace.Event) {
		events := make([]string, 0, len(evs))
		for _, ev := range evs {
			elapsed := ev.Time.Sub(node.StartTime())
			events = append(events, fmt.Sprintf("@%s %s", elapsed, ev.Name))
			appendAttributes(ev.Attributes)
		}

		logAttrs = append(logAttrs, kv.New("events", strings.Join(events, "; ")))
	}

	appendAttributes(node.Attributes())
	appendEvents(node.Events())

	s.log(logAttrs...)

	sort.Sort(node.children)
	for _, child := range node.children {
		s.printSpan(tree, child.SpanContext().SpanID(), traceStart, traceEnd)
	}
}
