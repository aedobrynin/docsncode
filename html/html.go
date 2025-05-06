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

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/mermaid"

	"docsncode/cfg"
)

// TODO: порефакторить код с парсингом блоков

// TODO: перестать использовать числовые константы в шаблонах (Code и Comment вместо 0 и 1)
// TODO: поправить отступы в HTML
var HTML_TEMPLATE = template.Must(template.New("docsncode").Parse(`<!DOCTYPE html>
<html>
<head>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles/default.min.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/highlight.min.js"></script>
</head>
<body>
    {{range .}}
        {{if eq .Type 0}}
			<pre><code>{{.Content}}</code></pre>
        {{else if eq .Type 1}}
			<p>
				<pre>{{.Content}}</pre>
			</p>
		{{end}}
	{{end}}
	<script>hljs.highlightAll();</script>
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

func isSingleLineCommentBlockStart(line string, languageInfo cfg.LanguageInfo) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, languageInfo.SingleLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, languageInfo.SingleLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_START_TOKEN)
}

func isSingleLineCommentBlockEnd(line string, languageInfo cfg.LanguageInfo) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, languageInfo.SingleLineCommentStartToken) {
		return false
	}
	line = strings.TrimPrefix(line, languageInfo.SingleLineCommentStartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN)
}

func isMultilineCommentBlockStart(line string, languageInfo cfg.LanguageInfo) bool {
	if languageInfo.MultilineCommentInfo == nil {
		return false
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, languageInfo.MultilineCommentInfo.StartToken) {
		return false
	}
	line = strings.TrimPrefix(line, languageInfo.MultilineCommentInfo.StartToken)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, cfg.COMMENT_BLOCK_START_TOKEN)
}

func isMultilineCommentBlockEnd(line string, languageInfo cfg.LanguageInfo) bool {
	if languageInfo.MultilineCommentInfo == nil {
		return false
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN) {
		return false
	}
	line = strings.TrimPrefix(line, cfg.COMMENT_BLOCK_END_TOKEN)
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, languageInfo.MultilineCommentInfo.EndToken)
}

func parseSingleLineCommentBlock(scanner *bufio.Scanner, languageInfo cfg.LanguageInfo) ([]byte, error) {
	log.Println("Start parsing single line comment block")
	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isSingleLineCommentBlockEnd(line, languageInfo) {
			log.Println("Found comment block end, stop parsing comment block raw content")
			return content, nil
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, languageInfo.SingleLineCommentStartToken) {
			log.Println("Line doesn't have comment start token even though we're inside comments block")
		}
		line = strings.TrimPrefix(line, languageInfo.SingleLineCommentStartToken)
		// TODO: trim leading spaces?
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, ErrCommentBlockEndNotFound
}

func parseMultilineCommentBlock(scanner *bufio.Scanner, languageInfo cfg.LanguageInfo) ([]byte, error) {
	log.Println("Start parsing multiline comment block")
	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isMultilineCommentBlockEnd(line, languageInfo) {
			log.Println("Found multiline comment block end, stop parsing comment block raw content")
			return content, nil
		}

		// TODO: поддержка отступов
		line = strings.TrimSpace(line)
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, ErrCommentBlockEndNotFound
}

func convertMarkdownToHTML(md []byte, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string) ([]byte, error) {
	// TODO: поддержка ссылок с абсолютным путём относительно корня проекта
	// TODO: убрать необходимость добавления .html для ссылки

	// TODO: не создавать новый конвертер на каждый файл
	converter := goldmark.New(
		// TODO: подумать, не нужен ли RenderModeServer?
		goldmark.WithExtensions(&mermaid.Extender{RenderMode: mermaid.RenderModeClient}),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(util.Prioritized(&linksResolverTransformer{
				absPathToProjectRoot: absPathToProjectRoot,
				absPathToCurrentFile: absPathToCurrentFile,
				absPathToResultDir:   absPathToResultDir,
				absPathToResultFile:  absPathToResultFile,
			}, 0)),
		),
	)

	var buf bytes.Buffer
	if err := converter.Convert(md, &buf); err != nil {
		return nil, fmt.Errorf("error on converting markdown to HTML: %w", err)
	}
	return buf.Bytes(), nil
}

// Assumes that comment block start line is already parsed
func parseAndBuildCommentBlock(scanner *bufio.Scanner, languageInfo cfg.LanguageInfo, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string, isMultiline bool) (*Block, error) {
	// TODO: учитывать отступ всего блока с комментарием
	log.Printf("Start parsing and building comment block, isMultiline=%t", isMultiline)

	var rawContent []byte
	var err error

	if isMultiline {
		rawContent, err = parseMultilineCommentBlock(scanner, languageInfo)
	} else {
		rawContent, err = parseSingleLineCommentBlock(scanner, languageInfo)
	}

	if err != nil {
		return nil, fmt.Errorf("error on parsing comment block: %w", err)
	}

	htmlContent, err := convertMarkdownToHTML(rawContent, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile)
	if err != nil {
		return nil, err
	}
	return &Block{
		Type:    Comment,
		Content: string(htmlContent),
	}, nil
}

func EscapeHTMLInCodeBlocks(blocks []Block) {
	for i, _ := range blocks {
		if blocks[i].Type != Code {
			continue
		}
		blocks[i].Content = template.HTMLEscapeString(blocks[i].Content)
	}
}

// TODO: нужен ли тут bytes.Buffer или достаточно []byte?
func BuildHTML(file *os.File, languageInfo cfg.LanguageInfo, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string) (*bytes.Buffer, error) {
	blocks := []Block{}
	scanner := bufio.NewScanner(file)

	var current_code_block_content []byte

	for scanner.Scan() {
		line := scanner.Text()

		isSingleLineCommentBlockStart := isSingleLineCommentBlockStart(line, languageInfo)
		isMultilineCommentBlockStart := isMultilineCommentBlockStart(line, languageInfo)

		if isSingleLineCommentBlockStart || isMultilineCommentBlockStart {
			log.Println("Found single line comment block start")
			if current_code_block_content != nil {
				log.Println("Append current code block")
				blocks = append(blocks, Block{
					Type: Code,
					// TODO: use unsafe?
					Content: string(current_code_block_content),
				})
				current_code_block_content = nil
			}

			block, err := parseAndBuildCommentBlock(scanner, languageInfo, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile, isMultilineCommentBlockStart)
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

	EscapeHTMLInCodeBlocks(blocks)

	resultBuf := bytes.NewBuffer([]byte{})
	err := HTML_TEMPLATE.Execute(resultBuf, blocks)

	if err != nil {
		return nil, fmt.Errorf("error on filling HTML template: %w", err)
	}
	return resultBuf, nil
}
