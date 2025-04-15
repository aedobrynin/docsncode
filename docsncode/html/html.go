package html

import (
	"bufio"
	"bytes"
	"docsncode/cfg"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// TODO: перестать использовать числовые константы в шаблонах (Code и Comment вместо 0 и 1)
var HTML_TEMPLATE = template.Must(template.New("docsncode").Parse(`<!DOCTYPE html>
<html>
<head>
    <style>
        .code_block { font-size: 16px; color: black; }
        .comment { font-size: 16px; color: red; font-weight: bold; }
    </style>
</head>
<body>
    {{range .}}
        {{if eq .Type 0}}
            <p class="code_block">
				<pre><code>{{.Content}}</code></pre>
			</p>
        {{else if eq .Type 1}}
            <p>
				<pre class="comment">{{.Content}}</pre>
			</p>
        {{end}}
    {{end}}
</body>
</html>
`))

type BlockType int

const (
	Code BlockType = iota
	Comment
)

type Block struct {
	Type    BlockType
	Content bytes.Buffer
}

var (
	ErrCommentBlockEndNotFound = errors.New("didn't see comment block end")
)

func isCommentBlockStart(line string, languageInfo cfg.LanguageInfo) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, languageInfo.OneLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, languageInfo.OneLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_START_TOKEN)
}

func isCommentBlockEnd(line string, languageInfo cfg.LanguageInfo) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, languageInfo.OneLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, languageInfo.OneLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN)
}

// Assumes that comment block start line is already parsed
func parseCommentBlock(scanner *bufio.Scanner, languageInfo cfg.LanguageInfo) (*Block, error) {
	log.Println("Started parsing comment block")
	buf := bytes.NewBuffer([]byte{})
	for scanner.Scan() {
		line := scanner.Text()

		if isCommentBlockEnd(line, languageInfo) {
			log.Println("Found comment block end, stop parsing comment block")
			return &Block{
				Type:    Comment,
				Content: *buf,
			}, nil
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, languageInfo.OneLineCommentStartToken) {
			log.Println("Line doesn't have comment start token even though we're inside comments block")
		}
		line = strings.TrimPrefix(line, languageInfo.OneLineCommentStartToken)
		// TODO: trim leading spaces?
		if buf.Len() != 0 {
			// TODO: check error
			buf.WriteRune('\n')
		}
		// TODO: check error
		buf.WriteString(line)
	}

	return nil, ErrCommentBlockEndNotFound
}

func BuildHTML(file *os.File, languageInfo cfg.LanguageInfo) (*bytes.Buffer, error) {
	blocks := []Block{}
	scanner := bufio.NewScanner(file)

	var current_code_block *Block

	for scanner.Scan() {
		line := scanner.Text()
		if isCommentBlockStart(line, languageInfo) {
			log.Println("Found comment block start")
			if current_code_block != nil {
				log.Println("Append current code block")
				blocks = append(blocks, *current_code_block)
				current_code_block = nil
			}

			block, err := parseCommentBlock(scanner, languageInfo)
			if err != nil {
				return nil, fmt.Errorf("error on parsing comment block: %w", err)
			}
			blocks = append(blocks, *block)
		} else {
			if current_code_block == nil {
				current_code_block = &Block{
					Type:    Code,
					Content: *bytes.NewBufferString(line),
				}
			} else {
				// TODO: check error
				current_code_block.Content.WriteRune('\n')
				// TODO: check error
				current_code_block.Content.WriteString(line)
			}
		}
	}

	if current_code_block != nil {
		log.Println("Append final code block")
		blocks = append(blocks, *current_code_block)
		current_code_block = nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error on scanning file: %w", err)
	}

	resultBuf := bytes.NewBuffer([]byte{})
	err := HTML_TEMPLATE.Execute(resultBuf, blocks)

	if err != nil {
		return nil, fmt.Errorf("error on filling HTML template: %w", err)
	}
	return resultBuf, nil
}
