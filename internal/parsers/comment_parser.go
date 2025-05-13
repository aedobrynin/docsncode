package parsers

import (
	"bufio"
	"errors"
)

var (
	ErrCommentBlockEndNotFound = errors.New("didn't see comment block end")
)

type CommentParser interface {
	// If Trigger returns true, we should execute the parser
	Trigger(line string) bool

	// Trigger(startLine) must be true
	Parse(startLine string, scanner *bufio.Scanner) (*ParsingResult, error)
}

type ParsingResult struct {
	Content     []byte
	BlockIndent int
}
