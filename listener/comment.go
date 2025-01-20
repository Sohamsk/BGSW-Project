package listener

import (
	"bosch/converter"
	"bosch/parser"
	"encoding/json"
)

func (s *TreeShapeListener) EnterComment(ctx *parser.CommentContext) {
	comment := converter.Comment{}
	comment.RuleType = "CommentRule"
	comment.CommentText = ctx.GetText()
	jsonData, err := json.Marshal(comment)
	if err != nil {
		panic(err)
	}
	s.writer.WriteString(string(jsonData) + ",")
}
