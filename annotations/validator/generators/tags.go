package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/validator/model"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

const tagName = "gav"

func buildTag(lit *ast.BasicLit) *model.Tag {
	if lit.Kind != token.STRING {
		return nil
	}
	line := splitTagLine(lit.Value)
	if line == "" {
		return nil
	}
	line = strings.Trim(line, `"`)
	items := strings.Split(line, ",")
	var tag model.Tag
	for _, item := range items {
		pair := strings.Split(item, "=")
		if len(pair) != 2 {
			continue
		}
		switch pair[0] {
		case "ignore":
			tag.Ignore = strings.ToLower(pair[1]) == "true"
		case "validator":
			tag.FuncName = pair[1]
		case "range":
			tag.Range = getRange(pair[1])
		case "enum":

		}
	}
	return &tag
}

func getRange(value string) model.Range {
	pair := strings.Split(value, "..")
	if len(pair) != 2 {
		panic(fmt.Sprintf("inavlid pair for range %s", value))
	}
	left, err := strconv.Atoi(pair[0])
	if err != nil {
		panic(err)
	}
	right, err := strconv.Atoi(pair[1])
	if err != nil {
		panic(err)
	}
	// TODO improve range validation
	return model.Range{
		Left:  left,
		Right: right,
	}
}

func splitTagLine(tag string) string {
	tag = strings.Trim(tag, "`")
	reg, err := regexp.Compile("gav:*\".*?\"")
	if err != nil {
		panic(err)
	}
	b := reg.Find([]byte(tag))

	str := string(b)
	pair := strings.Split(str, ":")
	if len(pair) != 2 {
		return ""
	}
	return pair[1]
}
