package parsers

import (
	"bufio"
	"docsncode/internal/cfg"
	"log"
	"strings"
	"unicode"
)

type baseSingleLineCommentBlockParser struct {
	singleLineCommentStartToken string
}

func newBaseSingleLineCommentBlockParser(singleLineCommentStartToken string) CommentParser {
	return &baseSingleLineCommentBlockParser{singleLineCommentStartToken: singleLineCommentStartToken}
}

func (p *baseSingleLineCommentBlockParser) Trigger(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, p.singleLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, p.singleLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_START_TOKEN)
}

// TODO: remove Fatalf
func (p *baseSingleLineCommentBlockParser) extractIndentFromSingleLineCommentBlock(line string) string {
	indx := strings.Index(line, p.singleLineCommentStartToken)
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

func (p *baseSingleLineCommentBlockParser) isSingleLineCommentBlockEnd(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, p.singleLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, p.singleLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN)
}

func (p *baseSingleLineCommentBlockParser) Parse(startLine string, scanner *bufio.Scanner) (*ParsingResult, error) {
	log.Println("Start parsing single line comment block")

	indent := p.extractIndentFromSingleLineCommentBlock(startLine)
	indentSize := calculateIndentSpacesCnt(indent)

	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if p.isSingleLineCommentBlockEnd(line) {
			log.Println("Found comment block end, stop parsing comment block raw content")
			return &ParsingResult{Content: content, BlockIndent: indentSize}, nil
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, p.singleLineCommentStartToken) {
			// TODO: make waring
			log.Println("Line doesn't have comment start token even though we're inside comments block")
		}
		line = strings.TrimPrefix(line, p.singleLineCommentStartToken)
		// TODO: trim leading spaces?
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, ErrCommentBlockEndNotFound
}
