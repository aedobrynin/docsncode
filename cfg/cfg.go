package cfg

import (
	"errors"
	"log"
)

// TODO: унести всё в yaml-конфиг
var (
	EXTENSION_TO_LANGUAGE_MAPPING            = map[string]string{".go": "Golang"}
	LANGUAGE_TO_ONE_LINE_COMMENT_START_TOKEN = map[string]string{"Golang": "//"}
	LANGUAGE_TO_MULTILINE_COMMENT_INFO       = map[string]*MultilineCommentInfo{"Golang": &MultilineCommentInfo{
		StartToken: "/*",
		EndToken:   "*/",
	}}

	// TODO: нет ли проблем с тем, что эти токены совпадают?
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
	Language                 string
	OneLineCommentStartToken string
	MultilineCommentInfo     *MultilineCommentInfo
}

func GetLanguageInfo(file_extension string) (*LanguageInfo, error) {
	language, isPresent := EXTENSION_TO_LANGUAGE_MAPPING[file_extension]
	if !isPresent {
		return nil, ErrExtensionNotSupported
	}

	oneLineCommentStartToken, isPresent := LANGUAGE_TO_ONE_LINE_COMMENT_START_TOKEN[language]
	if !isPresent {
		log.Fatalf("Language is listed in extensions mapping, but not found in one line comment start token mapping")
	}

	multilineCommentInfo := LANGUAGE_TO_MULTILINE_COMMENT_INFO[language]
	return &LanguageInfo{
		Language:                 language,
		OneLineCommentStartToken: oneLineCommentStartToken,
		MultilineCommentInfo:     multilineCommentInfo,
	}, nil
}
