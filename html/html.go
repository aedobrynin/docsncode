package html

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/mermaid"

	"docsncode/cfg"
	"docsncode/parsers"
	"docsncode/pathsignorer"
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

func convertMarkdownToHTML(md []byte, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string, pathsIgnorer pathsignorer.PathsIgnorer) ([]byte, error) {
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
				pathsIgnorer:         pathsIgnorer,
			}, 0)),
		),
	)

	var buf bytes.Buffer
	if err := converter.Convert(md, &buf); err != nil {
		return nil, fmt.Errorf("error on converting markdown to HTML: %w", err)
	}
	return buf.Bytes(), nil
}

func escapeHTMLInCodeBlocks(blocks []Block) {
	for i, _ := range blocks {
		if blocks[i].Type != Code {
			continue
		}
		blocks[i].Content = template.HTMLEscapeString(blocks[i].Content)
	}
}

func buildParsersByLanguage(language cfg.Language) []parsers.CommentParser {
	commentType := cfg.GetLanguageCommentsType(language)
	switch commentType {
	case cfg.CStyle:
		return []parsers.CommentParser{
			parsers.NewCStyleSingleLineCommentBlockParser(),
			parsers.NewCStyleMultilineCommentBlockParser(),
		}
	case cfg.PythonStyle:
		return []parsers.CommentParser{
			parsers.NewPythonStyleSingleLineCommentBlockParser(),
		}
	}
	// TODO: make log error
	log.Printf("Unexpected commentType=%s", commentType)
	return []parsers.CommentParser{}
}

// TODO(important): игнорировать блоки с кодом, в которых только пробельные символы (см. tests/links/link_to_website)
// TODO: нужен ли тут bytes.Buffer или достаточно []byte?
func BuildHTML(file *os.File, language cfg.Language, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile string, pathsIgnorer pathsignorer.PathsIgnorer) (*bytes.Buffer, error) {
	blocks := []Block{}
	scanner := bufio.NewScanner(file)

	var current_code_block_content []byte

	parsers := buildParsersByLanguage(language)

	for scanner.Scan() {
		line := scanner.Text()

		anyParserTriggered := false
		for _, parser := range parsers {
			if !parser.Trigger(line) {
				continue
			}
			// TODO: add parser name to log
			fmt.Println("Some parser triggered")
			anyParserTriggered = true

			if current_code_block_content != nil {
				log.Println("Append current code block")
				blocks = append(blocks, Block{
					Type: Code,
					// TODO: use unsafe?
					Content: string(current_code_block_content),
				})
				current_code_block_content = nil
			}

			parsingResult, err := parser.Parse(line, scanner)
			if err != nil {
				log.Printf("error on parsing: %s", err)
				anyParserTriggered = false
				continue
			}

			htmlContent, err := convertMarkdownToHTML(parsingResult.Content, absPathToProjectRoot, absPathToCurrentFile, absPathToResultDir, absPathToResultFile, pathsIgnorer)
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, Block{
				Type:            Comment,
				Content:         string(htmlContent),
				IndentSpacesCnt: parsingResult.BlockIndent,
			})
		}

		if anyParserTriggered {
			continue
		}

		if current_code_block_content == nil {
			current_code_block_content = []byte(line)
		} else {
			current_code_block_content = append(current_code_block_content, '\n')
			current_code_block_content = append(current_code_block_content, line...)
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

	escapeHTMLInCodeBlocks(blocks)

	resultBuf := bytes.NewBuffer([]byte{})

	err := htmlTemplate.Execute(resultBuf, htmlTemplateData{Blocks: blocks, HighlightJsLanguageName: cfg.GetHighlightJSLanguageName(language)})

	if err != nil {
		return nil, fmt.Errorf("error on filling HTML template: %w", err)
	}
	return resultBuf, nil
}
