package cursor

type Cursor struct {
	line, column int
	getHeight    func() int
	getLineWidth func(int) int

	contexts []string
}

func New(line, column int, height func() int, lineWidth func(int) int) Cursor {
	return Cursor{
		line:         line,
		column:       column,
		getHeight:    height,
		getLineWidth: lineWidth,
		contexts:     []string{},
	}
}

func (c *Cursor) Line() int {
	return c.line
}

func (c *Cursor) Column() int {
	return c.column
}

func (c *Cursor) CursorUp() {
	if c.Line() == 0 {
		c.column = 0
		return
	}

	c.line--

	var lenOfLine = c.getLineWidth(c.line) // c.CurrentLine().Length()
	if lenOfLine < c.column {
		c.column = lenOfLine
	}
}

func (c *Cursor) CursorDown() {
	if c.line == c.getHeight()-1 {
		c.column = c.getLineWidth(c.getHeight() - 1) //line.Length()
		return
	}

	c.line++

	var lenOfLine = c.getLineWidth(c.line) //c.CurrentLine().Length()
	if lenOfLine < c.column {
		c.column = lenOfLine
	}
}

func (c *Cursor) CursorLeft() {
	if c.column > 0 {
		c.column--
		return
	}

	if c.Line() == 0 {
		return
	}
	c.CursorUp()
	c.column = c.getLineWidth(c.line) //c.CurrentLine().Length()
}

func (c *Cursor) CursorRight() {
	if c.column < c.getLineWidth(c.line) { //c.CurrentLine().Length() {
		c.column++
		return
	}

	if c.Line() == c.getHeight()-1 {
		return
	}
	c.CursorDown()
	c.column = 0
}

func (c *Cursor) UpdateCursor(line, column int) {

}
