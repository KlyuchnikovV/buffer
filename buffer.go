package buffer

import (
	ctx "context"
	"log"
	"strings"

	"github.com/KlyuchnikovV/edigode/core/context"
	"github.com/KlyuchnikovV/edigode/types"
	"github.com/KlyuchnikovV/gapbuf"
	"github.com/KlyuchnikovV/linetree"
	"github.com/KlyuchnikovV/linetree/node"
	"golang.design/x/clipboard"
)

type Buffer struct {
	context.Context

	name string

	keyboardEvents chan types.KeyboardEvent

	Tree *linetree.LineTree
}

func New(ctx ctx.Context, name string, str ...byte) *Buffer {
	var (
		lines = strings.Split(string(str), "\n")
		tree  = linetree.New(
			node.New(*gapbuf.New()),
		)
		buffer = &Buffer{
			name: name,

			// TODO: length of keyboard chan
			keyboardEvents: make(chan types.KeyboardEvent, 100),
			// selectEvents:   make(chan types.SelectEvent, 100),

			Tree: tree,
			// Cursor: cursor.New(tree, 1, 0, 0),
		}
	)

	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	buffer.Context = *context.New(ctx, buffer.Init)

	if len(str) == 0 {
		return buffer
	}

	l, err := buffer.Tree.GetLine(1)
	if err != nil {
		panic(err)
	}
	l.Line.Insert(0, []byte(lines[0])...)

	for i := 1; i < len(lines); i++ {
		if err := buffer.Tree.Insert(i+1, *node.New(
			*gapbuf.New([]byte(lines[i])...),
		)); err != nil {
			panic(err)
		}
	}

	return buffer
}

func (buffer *Buffer) Start() error {
	if err := buffer.Context.Start(); err != nil {
		return err
	}

	go buffer.listenEvents()

	return nil
}

func (buffer *Buffer) Init() {
	buffer.Emit("buffer", "init", buffer.name)
}

func (buffer *Buffer) KeyboardEvents() chan types.KeyboardEvent {
	return buffer.keyboardEvents
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
		case <-buffer.Done():
			if err := buffer.Cancel(); err != nil {
				panic(err) //TODO: errors channel
			}
		}
	}
}

func (buffer *Buffer) HandleKeyboardEvent(event types.KeyboardEvent) error {
	log.Printf("TRACE: %s", event)

	var rangeDeleted bool
	if !event.Meta && (event.StartLine != event.EndLine || event.StartOffset != event.EndOffset) {
		if err := buffer.DeleteRange(event.StartLine, event.StartOffset, event.EndLine, event.EndOffset); err != nil {
			return err
		}
		rangeDeleted = true
	}

	switch event.Key {
	case "Delete":
	case "Backspace":
		if !rangeDeleted {
			buffer.Delete(event.StartLine, event.StartOffset)
		}
	case "Enter":
		buffer.NewLine(event.StartLine, event.StartOffset)
	case "Tab":
		buffer.Insert(event.StartLine, event.StartOffset, '\t')
	case "c":
		if event.Meta {
			buffer.Copy(event.StartLine, event.StartOffset, event.EndLine, event.EndOffset)
			return nil
		}
		fallthrough
	case "v":
		if event.Meta {
			buffer.Paste(event.StartLine, event.StartOffset, event.EndLine, event.EndOffset)
			break
		}
		fallthrough
	default:
		buffer.Insert(event.StartLine, event.StartOffset, []byte(event.Key)...)
	}

	buffer.Emit("buffer", "changed", event.Buffer)

	return nil
}

func (buffer *Buffer) NewLine(line, offset int, bytes ...byte) {
	gapBuffer := buffer.GetLine(line)
	line++

	if offset == -1 {
		offset = len(gapBuffer.String())
	}

	var (
		data = gapBuffer.Split(offset)
		text = append(data, bytes...)
		buf  = *gapbuf.New(text...)
		node = *node.New(buf)
	)

	if err := buffer.Tree.Insert(line, node); err != nil {
		panic(err)
	}
}

func (buffer *Buffer) Insert(line, offset int, bytes ...byte) {
	log.Printf("insert tab")
	gapBuffer := buffer.GetLine(line)
	gapBuffer.Insert(offset, bytes...)
}

func (buffer *Buffer) Delete(line, offset int) {
	var (
		gapBuffer = buffer.GetLine(line)
	)
	if offset != 0 {
		gapBuffer.Delete(offset)
		return
	}
	if line == 1 {
		return
	}

	oldLine, err := buffer.Tree.Remove(line)
	if err != nil {
		panic(err)
	}

	gapBuffer = buffer.GetLine(line - 1)

	gapBuffer.Insert(offset, []byte(oldLine.String())...)
}

func (buffer *Buffer) GetLines() ([]string, error) {
	nodes, err := buffer.Tree.GetRange(1, buffer.Tree.Size())
	if err != nil {
		return nil, err
	}
	var result = make([]string, len(nodes))
	for i, node := range nodes {
		result[i] = node.Line.String()
	}
	return result, nil
}

func (buffer *Buffer) GetLine(line int) *gapbuf.GapBuffer {
	node, err := buffer.Tree.GetLine(line)
	if err != nil {
		panic(err)
	}
	return &node.Line
}

func (buffer *Buffer) DeleteRange(startLine, startOffset, endLine, endOffset int) error {
	if startLine == endLine && startOffset == endOffset {
		return nil
	}

	var (
		line     = buffer.GetLine(startLine)
		lastLine *gapbuf.GapBuffer
		err      error
		offset   = endOffset - startOffset
	)
	if startLine != endLine {
		offset = -1
	}
	log.Printf("before: %s", line.String())
	line.DeleteRange(startOffset, offset)
	log.Printf("after: %s", line.String())

	for i := startLine + 1; i <= endLine; i++ {
		lastLine, err = buffer.Tree.Remove(startLine + 1)
		if err != nil {
			return err
		}
	}

	if startLine != endLine {
		lastLine.DeleteRange(0, endOffset)
		buffer.GetLine(startLine).Insert(startOffset, lastLine.Bytes()...)
	}

	return nil
}

func (buffer *Buffer) Copy(startLine, startOffset, endLine, endOffset int) {
	var (
		line   = buffer.GetLine(startLine)
		offset = endOffset
	)
	if startLine != endLine {
		offset = len(line.String())
	}
	var result = line.String()[startOffset:offset]
	for i := startLine + 1; i < endLine; i++ {
		line = buffer.GetLine(i)
		result += "\n" + line.String()
	}

	if startLine != endLine {
		line = buffer.GetLine(endLine)
		result += "\n" + line.String()[0:endOffset]
	}

	log.Printf("copied '%s'", result)
	clipboard.Write(clipboard.FmtText, []byte(result))
}

func (buffer *Buffer) Paste(startLine, startOffset, endLine, endOffset int) {
	if err := buffer.DeleteRange(startLine, startOffset, endLine, endOffset); err != nil {
		panic(err)
	}
	var (
		data   = clipboard.Read(clipboard.FmtText)
		line   = startLine
		offset = startOffset
	)
	for _, char := range data {
		if char == '\n' {
			line++
			buffer.NewLine(line, offset)
		} else {
			buffer.Insert(line, offset, char)
			offset++
		}
	}
	log.Printf("paste '%s'", string(data))
}
