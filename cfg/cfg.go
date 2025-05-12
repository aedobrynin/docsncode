package cfg

type Language string

const (
	Ada        Language = "Ada"
	C          Language = "C"
	CSharp     Language = "C#"
	Cpp        Language = "C++"
	D          Language = "D"
	Go         Language = "Go"
	Java       Language = "Java"
	JavaScript Language = "JavaScript"
	Lua        Language = "Lua"
	ObjectiveC Language = "Objective-C"
	Perl       Language = "Perl"
	PHP        Language = "PHP"
	Rust       Language = "Rust"
	Scala      Language = "Scala"
	Swift      Language = "Swift"
	Text       Language = "Text"
	TypeScript Language = "TypeScript"
)

type CommentType string

const (
	CStyle CommentType = "C-style"
)

// TODO: унести всё в yaml-конфиг
// TODO: валидация конфига
var (
	EXTENSION_TO_LANGUAGE_MAPPING = map[string]Language{
		".adb":   Ada,
		".ads":   Ada,
		".c":     C,
		".h":     C,
		".cs":    CSharp,
		".cpp":   Cpp,
		".hpp":   Cpp,
		".d":     D,
		".go":    Go,
		".java":  Java,
		".js":    JavaScript,
		".lua":   Lua,
		".m":     ObjectiveC,
		".pl":    Perl,
		".pm":    Perl,
		".php":   PHP,
		".rs":    Rust,
		".scala": Scala,
		".swift": Swift,
		".txt":   Text,
		".ts":    TypeScript,
	}

	LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME = map[Language]string{
		Ada:        "ada",
		C:          "c",
		CSharp:     "csharp",
		Cpp:        "c++",
		D:          "d",
		Go:         "golang",
		Java:       "java",
		JavaScript: "js",
		Lua:        "lua",
		ObjectiveC: "objectivec",
		Perl:       "perl",
		PHP:        "php",
		Rust:       "rust",
		Scala:      "scala",
		Swift:      "swift",
		TypeScript: "ts",
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

func GetLanguageCommentsType(language Language) CommentType {
	// TODO(important): support other comment styles
	return CStyle
}
