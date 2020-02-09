package memory

import "encoding/json"

var emptyLine = Line{}

type Writer struct {
	lines []Line
}

type Line struct {
	Fields  map[string]interface{} `json:"fields"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
}

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

func (w *Writer) Line(index int) (Line, bool) {
	if index < len(w.lines) {
		return w.lines[index], true
	}

	return emptyLine, false
}
