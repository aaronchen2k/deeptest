package service

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	queryHelper "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/query"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"regexp"
	"strings"
)

type ParserRegxService struct {
	XPathService  *XPathService  `inject:""`
	ParserService *ParserService `inject:""`
}

func (s *ParserRegxService) ParseRegx(req *v1.ParserRequest) (ret v1.ParserResponse, err error) {
	expr, _ := s.getRegxExpr(req.DocContent, req.SelectContent,
		req.StartLine, req.StartColumn, req.EndLine, req.EndColumn)

	exprType := "regx"
	result := queryHelper.RegxQuery(req.DocContent, expr)

	fmt.Printf("%s: %v", expr, result)

	ret = v1.ParserResponse{
		SelectionType: consts.NodeText,
		Expr:          expr,
		ExprType:      exprType,
	}

	return
}

func (s *ParserRegxService) getRegxExpr(docContent, selectContent string,
	startLine, startColumn, endLine, endColumn int) (ret string, err error) {

	lines := strings.Split(docContent, "\n")

	leftChars := ""
	prefix := s.ParserService.getLeftCharsInSingleLine(lines, startLine, startColumn, 6, &leftChars)
	leftChars = regexp.QuoteMeta(leftChars)

	rightChars := ""
	postfix := s.ParserService.getRightCharsInSingleLine(lines, endLine, endColumn, 6, &rightChars)
	rightChars = regexp.QuoteMeta(rightChars)

	expr := fmt.Sprintf("%s(.+)%s", leftChars, rightChars)

	ret = prefix + expr + postfix

	return
}
