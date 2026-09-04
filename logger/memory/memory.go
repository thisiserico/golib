// Package memory is a io.Writer implementation to use when testing log lines
// being produced.
package memory

import "encoding/json"

var emptyLine = Line{}

// Writer implements io.Writer and provides a way to fetch the log lines that
// were produced.
type Writer struct {
	lines []Line
}

// Line encapsulates the different elements that were logged.
type Line struct {
	// Fields contains the log tags.
	Fields map[string]interface{} `json:"fields"`

	// Level indicates the log level.
	Level string `json:"level"`

	// Message contains the actual message string.
	Message string `json:"message"`
}

// New returns a new Writer.
func New() *Writer {
	return &Writer{
		lines: make([]Line, 0),
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	var l Line
	if err := json.Unmarshal(p, &l); err != nil {
		return 0, err
	}

	w.lines = append(w.lines, l)
	return len(p), nil
}

// Line fetches the indicated log line. It also returns a boolean indicating
// whether the requested log line was produced.
func (w *Writer) Line(index int) (Line, bool) {
	if index < len(w.lines) {
		return w.lines[index], true
	}

	return emptyLine, false
}
