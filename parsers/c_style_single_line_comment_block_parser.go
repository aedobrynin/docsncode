package parsers

var (
	defaultCStyleSingleLineCommentBlockParser = newBaseSingleLineCommentBlockParser("//")
)

func NewCStyleSingleLineCommentBlockParser() CommentParser {
	return defaultCStyleSingleLineCommentBlockParser
}
