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
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/mermaid"

	"docsncode/cfg"
)

// TODO(important): порефакторить код с парсингом блоков

// TODO: перестать использовать числовые константы в шаблонах (Code и Comment вместо 0 и 1)
// TODO: не подключать highlight.js, если в файле не будет блоков с кодом
// TODO: вынести настройку tab-size в конфиг
var htmlTemplate = template.Must(template.New("docsncode").Parse(`<!DOCTYPE html>
<html>
<head>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles/default.min.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/highlight.min.js"></script>
	<style>pre {tab-size: 4ch;}</style>
</head>
<body>
    {{range .Blocks}}
        {{if eq .Type 0}}
			{{if $.HighlightJsLanguageName }}
				<pre><code class="language-{{$.HighlightJsLanguageName}}">{{.Content}}</code></pre>
			{{else}}
				<pre><code>{{.Content}}</code></pre>
			{{end}}
        {{else if eq .Type 1}}
			<div style="padding-left: calc({{.IndentSpacesCnt}}ch + 1em); font-size:12px;">{{.Content}}</div>
		{{end}}
	{{end}}
	<script>hljs.highlightAll();</script>
</body>
</html>
`))

type htmlTemplateData struct {
	Blocks                  []Block
	HighlightJsLanguageName *string
}

// TODO: make internal
type BlockType int

const (
	Code BlockType = iota
	Comment
)

// TODO: растащить на две структуры
type Block struct {
	Type            BlockType
	Content         string
	IndentSpacesCnt int
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

type Indent = string

// TODO: remove Fatalf
func extractIndentFromSingleLineCommentBlockStart(line, singleLineCommentStartToken string) Indent {
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

func parseSingleLineCommentBlock(scanner *bufio.Scanner, startLine string, languageInfo cfg.LanguageInfo) ([]byte, Indent, error) {
	log.Println("Start parsing single line comment block")

	indent := extractIndentFromSingleLineCommentBlockStart(startLine, languageInfo.SingleLineCommentStartToken)

	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isSingleLineCommentBlockEnd(line, languageInfo) {
			log.Println("Found comment block end, stop parsing comment block raw content")
			return content, indent, nil
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

	return nil, "", ErrCommentBlockEndNotFound
}

// TODO: remove Fatalf
func extractIndentFromMultilineCommentBlock(line, multilineCommentStartToken string) Indent {
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

func parseMultilineCommentBlock(scanner *bufio.Scanner, startLine string, languageInfo cfg.LanguageInfo) ([]byte, Indent, error) {
	log.Println("Start parsing multiline comment block")

	indent := extractIndentFromMultilineCommentBlock(startLine, languageInfo.MultilineCommentInfo.StartToken)

	var content []byte
	for scanner.Scan() {
		line := scanner.Text()

		if isMultilineCommentBlockEnd(line, languageInfo) {
			log.Println("Found multiline comment block end, stop parsing comment block raw content")
			return content, indent, nil
		}

		line = strings.TrimSpace(line)
		if len(content) != 0 {
			content = append(content, '\n')
		}
		content = append(content, line...)
	}

	return nil, "", ErrCommentBlockEndNotFound
}

func convertMarkdownToHTML(md []byte, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string) ([]byte, error) {
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

func calculateIndentSpacesCnt(indent Indent) int {
	cnt := 0
	for _, r := range indent {
		if r == '\t' {
			cnt += 4 // TODO: move tab size to config
		} else {
			cnt++
		}
	}
	return cnt
}

// Assumes that comment block start line is already parsed
func parseAndBuildCommentBlock(scanner *bufio.Scanner, startLine string, languageInfo cfg.LanguageInfo, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string, isMultiline bool) (*Block, error) {
	log.Printf("Start parsing and building comment block, isMultiline=%t", isMultiline)

	var rawContent []byte
	var indent Indent
	var err error

	if isMultiline {
		rawContent, indent, err = parseMultilineCommentBlock(scanner, startLine, languageInfo)
	} else {
		rawContent, indent, err = parseSingleLineCommentBlock(scanner, startLine, languageInfo)
	}

	if err != nil {
		return nil, fmt.Errorf("error on parsing comment block: %w", err)
	}

	htmlContent, err := convertMarkdownToHTML(rawContent, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile)
	if err != nil {
		return nil, err
	}
	return &Block{
		Type:            Comment,
		Content:         string(htmlContent),
		IndentSpacesCnt: calculateIndentSpacesCnt(indent),
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

// TODO(important): игнорировать блоки с кодом, в которых только пробельные символы (см. tests/links/link_to_website)
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

			block, err := parseAndBuildCommentBlock(scanner, line, languageInfo, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile, isMultilineCommentBlockStart)
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
			Content:         string(current_code_block_content),
			IndentSpacesCnt: 0,
		})
		current_code_block_content = nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error on scanning file: %w", err)
	}

	EscapeHTMLInCodeBlocks(blocks)

	resultBuf := bytes.NewBuffer([]byte{})

	err := htmlTemplate.Execute(resultBuf, htmlTemplateData{Blocks: blocks, HighlightJsLanguageName: languageInfo.HighlightJsLanguageName})

	if err != nil {
		return nil, fmt.Errorf("error on filling HTML template: %w", err)
	}
	return resultBuf, nil
}
