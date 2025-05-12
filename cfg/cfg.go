package cfg

type Language string

// TODO: унести всё в yaml-конфиг
// TODO: валидация конфига
var (
	EXTENSION_TO_LANGUAGE_MAPPING = map[string]Language{
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

	LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME = map[Language]string{
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

func GetLanguageNameIfSupported(fileExtension string) *Language {
	lang, isPresent := EXTENSION_TO_LANGUAGE_MAPPING[fileExtension]
	if !isPresent {
		return nil
	}
	return &lang
}

func GetHighlightJSLanguageName(language Language) *string {
	name, isPresent := LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME[language]
	if !isPresent {
		return nil
	}
	return &name
}
