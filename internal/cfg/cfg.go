package cfg

type Language string

const (
	Ada          Language = "Ada"
	Bash         Language = "Bash"
	C            Language = "C"
	CoffeeScript Language = "CoffeeScript"
	CSharp       Language = "C#"
	Cpp          Language = "C++"
	D            Language = "D"
	FSharp       Language = "F#"
	Go           Language = "Go"
	Java         Language = "Java"
	JavaScript   Language = "JavaScript"
	Lua          Language = "Lua"
	ObjectiveC   Language = "Objective-C"
	Perl         Language = "Perl"
	PHP          Language = "PHP"
	Python       Language = "Python"
	Ruby         Language = "Ruby"
	Rust         Language = "Rust"
	Scala        Language = "Scala"
	Swift        Language = "Swift"
	Text         Language = "Text"
	TypeScript   Language = "TypeScript"
)

type CommentType string

const (
	CStyle      CommentType = "C-style"
	PythonStyle CommentType = "Python-style"
)

// TODO: унести всё в yaml-конфиг
// TODO: валидация конфига
var (
	EXTENSION_TO_LANGUAGE_MAPPING = map[string]Language{
		".adb":    Ada,
		".ads":    Ada,
		".sh":     Bash,
		".c":      C,
		".h":      C,
		".coffee": CoffeeScript,
		".cpp":    Cpp,
		".hpp":    Cpp,
		".cs":     CSharp,
		".d":      D,
		".fs":     FSharp,
		".go":     Go,
		".java":   Java,
		".js":     JavaScript,
		".lua":    Lua,
		".m":      ObjectiveC,
		".pl":     Perl,
		".pm":     Perl,
		".php":    PHP,
		".py":     Python,
		".rb":     Ruby,
		".rs":     Rust,
		".scala":  Scala,
		".swift":  Swift,
		".txt":    Text,
		".ts":     TypeScript,
	}

	LANGUAGE_TO_HIGHLIGHT_JS_LANGUAGE_NAME = map[Language]string{
		Ada:          "ada",
		Bash:         "bash",
		C:            "c",
		CoffeeScript: "coffeescript",
		Cpp:          "c++",
		CSharp:       "csharp",
		D:            "d",
		Go:           "golang",
		Java:         "java",
		JavaScript:   "js",
		Lua:          "lua",
		ObjectiveC:   "objectivec",
		Perl:         "perl",
		PHP:          "php",
		Python:       "python",
		Ruby:         "ruby",
		Rust:         "rust",
		Scala:        "scala",
		Swift:        "swift",
		TypeScript:   "ts",
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
	switch language {
	case Bash, CoffeeScript, FSharp, Python, Ruby:
		return PythonStyle
	default:
		return CStyle
	}
}
