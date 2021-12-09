package buffer

import (
	"context"
	"log"
	"strings"

	"github.com/KlyuchnikovV/buffer/cursor"
	"github.com/KlyuchnikovV/edigode/types"
	"github.com/KlyuchnikovV/gapbuf"
	"github.com/KlyuchnikovV/linetree"
	"github.com/KlyuchnikovV/linetree/node"
	"github.com/wailsapp/wails"
)

type Buffer struct {
	types.CancellableContext
	runtime *wails.Runtime

	keyboardEvents chan types.KeyboardEvent
	selectEvents   chan types.SelectEvent

	cursor.Cursors
}

func New(ctx context.Context, str ...byte) *Buffer {
	var (
		lines  = strings.Split(string(str), "\n")
		buffer = &Buffer{
			CancellableContext: *types.NewCancellableContext(ctx),

			keyboardEvents: make(chan types.KeyboardEvent, 100), //TODO: length of keyboard chan
			selectEvents:   make(chan types.SelectEvent, 100),

			Cursors: *cursor.New(
				linetree.New(
					node.New(*gapbuf.New()),
				),
				1,
			),
		}
	)

	if len(str) == 0 {
		return buffer
	}

	l, err := buffer.Tree.GetLine(1)
	if err != nil {
		panic(err)
	}
	l.Line.Insert([]byte(lines[0])...)

	for i := 1; i < len(lines); i++ {
		if err := buffer.Tree.Insert(i+1, *node.New(
			*gapbuf.New([]byte(lines[i])...),
		)); err != nil {
			panic(err)
		}
	}

	return buffer
}

func (buffer *Buffer) WailsInit(runtime *wails.Runtime) error {
	buffer.runtime = runtime

	return nil
}

func (buffer *Buffer) Start() error {
	if err := buffer.CancellableContext.Start(); err != nil {
		return err
	}

	go buffer.listenEvents()

	return nil
}

func (buffer *Buffer) KeyboardEvents() chan types.KeyboardEvent {
	return buffer.keyboardEvents
}

func (buffer *Buffer) SelectEvents() chan types.SelectEvent {
	return buffer.selectEvents
}

func (buffer *Buffer) listenEvents() {
	for {
		select {
		case event, ok := <-buffer.keyboardEvents:
			if !ok {
				panic("!ok")
			}
			if err := buffer.HandleKeyboardEvent(event); err != nil {
				panic(err)
			}
		case event, ok := <-buffer.selectEvents:
			if !ok {
				panic("!ok")
			}
			if err := buffer.HandleSelect(event); err != nil {
				panic(err)
			}
		case <-buffer.CancellableContext.Done():
			if err := buffer.CancellableContext.Cancel(); err != nil {
				panic(err) //TODO: errors channel
			}
		}
	}
}

func (buffer *Buffer) HandleKeyboardEvent(event types.KeyboardEvent) error {
	log.Printf("TRACE: %s", event)

	switch {
	case event.Key == "Delete":
	case event.Key == "Backspace":
		buffer.Delete()
	case event.Key == "Enter":
		buffer.NewLine()
	case event.Key == "Tab":
		buffer.Insert('\t')
	default:
		buffer.Insert([]byte(event.Key)...)
	}

	buffer.runtime.Events.Emit("buffer_changed" + event.Buffer)

	return nil
}

func (buffer *Buffer) HandleSelect(event types.SelectEvent) error {
	log.Printf("SELECT: %s", event)

	var line, offset = buffer.translateOffset(event.Symbol)
	log.Printf("SELECT: line %d, offset %d (symbol %d)", line, offset, event.Symbol)

	buffer.SetCursor(line, offset)

	buffer.runtime.Events.Emit("cursor_changed", event.Buffer)

	return nil
}

func (buffer *Buffer) NewLine(bytes ...byte) {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
		(*lines)[index]++

		var (
			data = gapBuffer.Split()
			text = append(data, bytes...)
			buf  = *gapbuf.New(text...)
			node = *node.New(buf)
		)

		if err := buffer.Tree.Insert((*lines)[index], node); err != nil {
			panic(err)
		}
	})
}

func (buffer *Buffer) SetCursor(line, offset int) {
	if line > 0 && line < buffer.Tree.Size() {
		buffer.Cursors.Lines[0] = line
	}

	if offset != -1 {
		buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
			gapBuffer.SetCursor(offset)
		})
	}
}

func (buffer *Buffer) Insert(bytes ...byte) {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, _ *[]int, _ int) {
		gapBuffer.Insert(bytes...)
	})
}

func (buffer *Buffer) Delete() {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
		if gapBuffer.GetCursor() != 0 {
			gapBuffer.Delete()
			return
		}
		if (*lines)[index] == 1 {
			return
		}

		oldLine, err := buffer.Tree.Remove((*lines)[index])
		if err != nil {
			panic(err)
		}
		(*lines)[index]--

		node, err := buffer.Tree.GetLine((*lines)[index])
		if err != nil {
			panic(err)
		}

		gapBuffer = &node.Line

		gapBuffer.Gap.SetCursor(gapBuffer.Size())

		gapBuffer.Insert([]byte(oldLine.String())...)
	})
}

func (buffer *Buffer) translateOffset(offset int) (int, int) {
	var (
		data      = buffer.Tree.String()[:offset]
		line      = strings.Count(data, "\n") + 1
		lastIndex = strings.LastIndex(data, "\n")
	)

	if lastIndex == -1 {
		return line, offset
	}

	return line, len(data[lastIndex+1:])
}

func (buffer *Buffer) GetLines() ([]string, error) {
	nodes, err := buffer.Tree.GetRange(1, buffer.Tree.Size())
	if err != nil {
		return nil, err
	}
	var result = make([]string, len(nodes))
	for i, node := range nodes {
		result[i] = node.Line.String()
		result[i] = strings.ReplaceAll(result[i], " ", "\u00A0")
	}
	return result, nil
}
