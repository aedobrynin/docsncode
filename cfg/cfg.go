package cfg

import (
	"errors"
	"log"
)

// TODO: унести всё в yaml-конфиг
// TODO: валидация конфига
var (
	EXTENSION_TO_LANGUAGE_MAPPING = map[string]string{
		".adb":   "ADA",
		".ads":   "ADA",
		".c":     "C",
		".h":     "C",
		".cs":    "C#",
		".cpp":   "C++",
		".hpp":   "C++",
		".d":     "D",
		".go":    "Golang",
		".java":  "Java",
		".js":    "JavaScript",
		".lua":   "Lua",
		".m":     "Objective-C",
		".pl":    "Perl",
		".pm":    "Perl",
		".php":   "PHP",
		".rs":    "Rust",
		".scala": "Scala",
		".swift": "Swift",
		".txt":   "Text",
		".ts":    "TypeScript",
	}
	LANGUAGE_TO_SINGLE_LINE_COMMENT_START_TOKEN = map[string]string{
		"ADA":         "//",
		"C":           "//",
		"C#":          "//",
		"C++":         "//",
		"Golang":      "//",
		"Java":        "//",
		"JavaScript":  "//",
		"Lua":         "//",
		"Objective-C": "//",
		"Perl":        "//",
		"PHP":         "//",
		"Rust":        "//",
		"Scala":       "//",
		"Swift":       "//",
		"TypeScript":  "//",
		"Text":        "//",
	}

	C_STYLE_MULTILINE_COMMENT_INFO = MultilineCommentInfo{
		StartToken: "/*",
		EndToken:   "*/",
	}

	LANGUAGE_TO_MULTILINE_COMMENT_INFO = map[string]*MultilineCommentInfo{
		"ADA":         &C_STYLE_MULTILINE_COMMENT_INFO,
		"C":           &C_STYLE_MULTILINE_COMMENT_INFO,
		"C#":          &C_STYLE_MULTILINE_COMMENT_INFO,
		"C++":         &C_STYLE_MULTILINE_COMMENT_INFO,
		"D":           &C_STYLE_MULTILINE_COMMENT_INFO,
		"Golang":      &C_STYLE_MULTILINE_COMMENT_INFO,
		"Java":        &C_STYLE_MULTILINE_COMMENT_INFO,
		"JavaScript":  &C_STYLE_MULTILINE_COMMENT_INFO,
		"Lua":         &C_STYLE_MULTILINE_COMMENT_INFO,
		"Objective-C": &C_STYLE_MULTILINE_COMMENT_INFO,
		"Perl":        &C_STYLE_MULTILINE_COMMENT_INFO,
		"PHP":         &C_STYLE_MULTILINE_COMMENT_INFO,
		"Rust":        &C_STYLE_MULTILINE_COMMENT_INFO,
		"Scala":       &C_STYLE_MULTILINE_COMMENT_INFO,
		"Swift":       &C_STYLE_MULTILINE_COMMENT_INFO,
		"TypeScript":  &C_STYLE_MULTILINE_COMMENT_INFO,
	}

	LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME = map[string]string{
		"ADA":         "ada",
		"C":           "c",
		"C#":          "csharp",
		"C++":         "c++",
		"D":           "d",
		"Golang":      "golang",
		"Java":        "java",
		"JavaScript":  "js",
		"Lua":         "lua",
		"Objective-C": "objectivec",
		"Perl":        "perl",
		"PHP":         "php",
		"Rust":        "rust",
		"Scala":       "scala",
		"Swift":       "swift",
		"TypeScript":  "ts",
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
