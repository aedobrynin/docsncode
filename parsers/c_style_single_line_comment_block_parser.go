package parsers

import (
	"bufio"
	"docsncode/cfg"
	"log"
	"strings"
	"unicode"
)

type cStyleSingleLineCommentBlockParser struct {
}

var (
	defaultcStyleSingleLineCommentBlockParser = cStyleSingleLineCommentBlockParser{}

	singleLineCommentStartToken = "//"
)

func NewCStyleSingleLineCommentBlockParser() CommentParser {
	return &defaultcStyleSingleLineCommentBlockParser
}

func (cStyleSingleLineCommentBlockParser) Trigger(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, singleLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, singleLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_START_TOKEN)
}

// TODO: remove Fatalf
func extractIndentFromSingleLineCommentBlock(line string) string {
	indx := strings.Index(line, singleLineCommentStartToken)
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

func isSingleLineCommentBlockEnd(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, singleLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, singleLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN)
}

func (cStyleSingleLineCommentBlockParser) Parse(startLine string, scanner *bufio.Scanner) (*ParsingResult, error) {
	log.Println("Start parsing single line comment block")

	indent := extractIndentFromSingleLineCommentBlock(startLine)
	indentSize := calculateIndentSpacesCnt(indent)

	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isSingleLineCommentBlockEnd(line) {
			log.Println("Found comment block end, stop parsing comment block raw content")
			return &ParsingResult{Content: content, BlockIndent: indentSize}, nil
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, singleLineCommentStartToken) {
			// TODO: make waring
			log.Println("Line doesn't have comment start token even though we're inside comments block")
		}
		line = strings.TrimPrefix(line, singleLineCommentStartToken)
		// TODO: trim leading spaces?
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, ErrCommentBlockEndNotFound
}
