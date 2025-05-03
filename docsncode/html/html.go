package html

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"docsncode/cfg"

	"github.com/yuin/goldmark"
)

// TODO: перестать использовать числовые константы в шаблонах (Code и Comment вместо 0 и 1)
var HTML_TEMPLATE = template.Must(template.New("docsncode").Parse(`<!DOCTYPE html>
<html>
<head>
    <style>
        .code_block { font-size: 16px; color: black; }
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
				{{.Content}}
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
	Content string
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

func parseCommentBlock(scanner *bufio.Scanner, languageInfo cfg.LanguageInfo) ([]byte, error) {
	log.Println("Start parsing comment block")
	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isCommentBlockEnd(line, languageInfo) {
			log.Println("Found comment block end, stop parsing comment block raw content")
			return content, nil
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, languageInfo.OneLineCommentStartToken) {
			log.Println("Line doesn't have comment start token even though we're inside comments block")
		}
		line = strings.TrimPrefix(line, languageInfo.OneLineCommentStartToken)
		// TODO: trim leading spaces?
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, ErrCommentBlockEndNotFound
}

func convertMarkdownToHTML(md []byte) ([]byte, error) {
	// TODO: поддержка ссылок с абсолютным путём относительно корня проекта
	// TODO: убрать необходимость добавления .html для ссылки
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return nil, fmt.Errorf("error on converting markdown to HTML: %v", err)
	}
	return buf.Bytes(), nil
}

// Assumes that comment block start line is already parsed
func parseAndBuildCommentBlock(scanner *bufio.Scanner, languageInfo cfg.LanguageInfo) (*Block, error) {
	// TODO: учитывать отступ всего блока с комментарием

	log.Println("Start parsing and building comment block")
	rawContent, err := parseCommentBlock(scanner, languageInfo)
	if err != nil {
		return nil, fmt.Errorf("error on parsing comment block: %w", err)
	}

	htmlContent, err := convertMarkdownToHTML(rawContent)
	if err != nil {
		return nil, err
	}
	return &Block{
		Type:    Comment,
		Content: string(htmlContent),
	}, nil
}

// TODO: нужен ли тут bytes.Buffer или достаточно []byte?
func BuildHTML(file *os.File, languageInfo cfg.LanguageInfo) (*bytes.Buffer, error) {
	blocks := []Block{}
	scanner := bufio.NewScanner(file)

	var current_code_block_content []byte

	for scanner.Scan() {
		line := scanner.Text()
		if isCommentBlockStart(line, languageInfo) {
			log.Println("Found comment block start")
			if current_code_block_content != nil {
				log.Println("Append current code block")
				blocks = append(blocks, Block{
					Type: Code,
					// TODO: use unsafe?
					Content: string(current_code_block_content),
				})
				current_code_block_content = nil
			}

			block, err := parseAndBuildCommentBlock(scanner, languageInfo)
			if err != nil {
				return nil, fmt.Errorf("error on parsing comment block: %w", err)
			}
			blocks = append(blocks, *block)
		} else {
			if current_code_block_content == nil {
				current_code_block_content = []byte(line)
			} else {
				current_code_block_content = append(current_code_block_content, '\n')
				current_code_block_content = append(current_code_block_content, line...)
			}
		}
	}

	if current_code_block_content != nil {
		log.Println("Append final code block")
		blocks = append(blocks, Block{
			Type: Code,
			// TODO: use unsafe?
			Content: string(current_code_block_content),
		})
		current_code_block_content = nil
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
