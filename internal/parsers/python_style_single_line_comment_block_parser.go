package parsers

var (
	defaultPythonStyleSingleLineCommentBlockParser = newBaseSingleLineCommentBlockParser("#")
)

func NewPythonStyleSingleLineCommentBlockParser() CommentParser {
	return defaultPythonStyleSingleLineCommentBlockParser
}
