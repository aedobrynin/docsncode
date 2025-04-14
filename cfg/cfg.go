package cfg

import (
	"errors"
	"log"
)

// TODO: унести всё в yaml-конфиг
var (
	EXTENSION_TO_LANGUAGE_MAPPING            = map[string]string{".go": "Golang"}
	LANGUAGE_TO_ONE_LINE_COMMENT_START_TOKEN = map[string]string{"Golang": "//"}

	COMMENT_BLOCK_START_TOKEN = "@docsncode_comment_block_start"
	COMMENT_BLOCK_END_TOKEN   = "@docsncode_comment_block_end"
)

var (
	ErrExtensionNotSupported = errors.New("provided file extension is not supported")
)

type LanguageInfo struct {
	Language                 string
	OneLineCommentStartToken string
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

	return &LanguageInfo{
		Language:                 language,
		OneLineCommentStartToken: oneLineCommentStartToken,
	}, nil
}
