package messages

import "os"

type RenderMessage interface {
	GetSource() string
}

type RenderRequest struct {
	Row    int
	Column int
	Runes  []rune
	Source string
}

func (r RenderRequest) GetSource() string {
	return r.Source
}

type CursorReposition struct {
	Row    int
	Column int
	Source string
}

func (c CursorReposition) GetSource() string {
	return c.Source
}

type BufferChange struct {
	Row    int
	Source string
}

func (b BufferChange) GetSource() string {
	return b.Source
}

type InputMessage interface {
}

type ResizeSignal os.Signal
