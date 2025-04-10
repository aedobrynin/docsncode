package main

// TODO: унести всё в yaml-конфиг
var (
	SUPPORTED_EXTENSIONS                     = []string{".go"}
	EXTENSION_TO_LANGUAGE_MAPPING            = map[string]string{".go": "Golang"}
	LANGUAGE_TO_ONE_LINE_COMMENT_START_TOKEN = map[string]string{"Golang": "//"}

	COMMENT_BLOCK_START_TOKEN = "@docsncode_comment_block_start"
	COMMENT_BLOCK_END_TOKEN   = "@docsncode_comment_block_end"
)

func IsExtensionSupported(extension string) bool {
	for _, ext := range SUPPORTED_EXTENSIONS {
		if ext == extension {
			return true
		}
	}
	return false
}
