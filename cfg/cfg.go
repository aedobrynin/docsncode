package cfg

import (
	"errors"
	"log"
)

// TODO: унести всё в yaml-конфиг
// TODO: валидация конфига
var (
	EXTENSION_TO_LANGUAGE_MAPPING = map[string]string{
		".go":   "Golang",
		".txt":  "Text",
		".cpp":  "C++",
		".hpp":  "C++",
		".c":    "C",
		".h":    "C",
		".java": "Java",
		".cs":   "C#",
	}
	LANGUAGE_TO_SINGLE_LINE_COMMENT_START_TOKEN = map[string]string{
		"Golang": "//",
		"Text":   "//",
		"C++":    "//",
		"C":      "//",
		"Java":   "//",
		"C#":     "//",
	}

	C_STYLE_MULTILINE_COMMENT_INFO = MultilineCommentInfo{
		StartToken: "/*",
		EndToken:   "*/",
	}

	LANGUAGE_TO_MULTILINE_COMMENT_INFO = map[string]*MultilineCommentInfo{
		"Golang": &C_STYLE_MULTILINE_COMMENT_INFO,
		"C++":    &C_STYLE_MULTILINE_COMMENT_INFO,
		"C":      &C_STYLE_MULTILINE_COMMENT_INFO,
		"Java":   &C_STYLE_MULTILINE_COMMENT_INFO,
		"C#":     &C_STYLE_MULTILINE_COMMENT_INFO,
	}

	LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME = map[string]string{
		"Golang": "golang",
		"C++":    "c++",
		"C":      "c",
		"Java":   "java",
		"C#":     "csharp",
	}

	COMMENT_BLOCK_START_TOKEN = "@docsncode"
	COMMENT_BLOCK_END_TOKEN   = "@docsncode"
)

var (
	ErrExtensionNotSupported = errors.New("provided file extension is not supported")
)

type MultilineCommentInfo struct {
	StartToken string
	EndToken   string
}

type LanguageInfo struct {
	Language                    string
	SingleLineCommentStartToken string
	MultilineCommentInfo        *MultilineCommentInfo
	HighlightJsLanguageName     *string
}

func GetLanguageInfo(file_extension string) (*LanguageInfo, error) {
	language, isPresent := EXTENSION_TO_LANGUAGE_MAPPING[file_extension]
	if !isPresent {
		return nil, ErrExtensionNotSupported
	}

	singleLineCommentStartToken, isPresent := LANGUAGE_TO_SINGLE_LINE_COMMENT_START_TOKEN[language]
	if !isPresent {
		log.Fatalf("Language is listed in extensions mapping, but not found in single line comment start token mapping")
	}

	multilineCommentInfo := LANGUAGE_TO_MULTILINE_COMMENT_INFO[language]

	var highlightJsLanguageName *string
	hljsLanguageName, isPresent := LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME[language]
	if isPresent {
		highlightJsLanguageName = &hljsLanguageName
	}

	return &LanguageInfo{
		Language:                    language,
		SingleLineCommentStartToken: singleLineCommentStartToken,
		MultilineCommentInfo:        multilineCommentInfo,
		HighlightJsLanguageName:     highlightJsLanguageName,
	}, nil
}
