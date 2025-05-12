package parsers

var (
	defaultPythonStyleSingleLineCommentBlockParser = newBaseSingleLineCommentBlockParser("#")
)

func NewCPythonStyleSingleLineCommentBlockParser() CommentParser {
	return defaultPythonStyleSingleLineCommentBlockParser
}
