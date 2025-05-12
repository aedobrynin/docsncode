package parsers

import (
	"bufio"
	"docsncode/cfg"
	"log"
	"strings"
	"unicode"
)

type cStyleMultilineCommentBlockParser struct {
}

var (
	defaultcStyleMultilineCommentBlockParser = cStyleMultilineCommentBlockParser{}

	multilineCommentStartToken = "/*"
	multilineCommentEndToken   = "*/"
)

func NewCStyleMultilineCommentBlockParser() CommentParser {
	return &defaultcStyleMultilineCommentBlockParser
}

func (cStyleMultilineCommentBlockParser) Trigger(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, multilineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, multilineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_START_TOKEN)
}

// TODO: remove Fatalf
func extractIndentFromMultilineCommentBlock(line string) string {
	indx := strings.Index(line, multilineCommentStartToken)
	if indx == -1 {
		log.Fatalf("The line should be start of comment block, but it isn't")
	}

	for _, r := range line[:indx] {
		if !unicode.IsSpace(r) {
			log.Fatalf("The line should be start of comment block, but it isn't")
		}
	}
	return line[:indx]
}

func isMultilineCommentBlockEnd(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN) {
		return false
	}
	line = strings.TrimPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, multilineCommentEndToken)
}

func (cStyleMultilineCommentBlockParser) Parse(startLine string, scanner *bufio.Scanner) (*ParsingResult, error) {
	log.Println("Start parsing multiline comment block")

	indent := extractIndentFromMultilineCommentBlock(startLine)
	indentSize := calculateIndentSpacesCnt(indent)

	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isMultilineCommentBlockEnd(line) {
			log.Println("Found multiline comment block end, stop parsing comment block raw content")
			return &ParsingResult{
				Content:     content,
				BlockIndent: indentSize,
			}, nil
		}

		line = strings.TrimSpace(line)
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, ErrCommentBlockEndNotFound
}
