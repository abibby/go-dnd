package blade

import (
	"fmt"
	"regexp"
)

type Parser struct {
	directives map[string]func(args []string) string
}

func New() *Parser {
	return &Parser{
		directives: map[string]func(args []string) string{},
	}
}

func (b *Parser) Directive(name string, callback func(args []string) string) {
	b.directives[name] = callback
}

func (b *Parser) Parse(doc string) string {
	for name, callback := range b.directives {
		re := regexp.MustCompile(fmt.Sprintf(`\@%s\(([^\)]*)\)`, regexp.QuoteMeta(name)))
		doc = re.ReplaceAllStringFunc(doc, func(s string) string {
			args := []string{}
			arg := ""
			quoted := rune(0)
			for _, c := range re.FindStringSubmatch(s)[1] {
				if c == '"' || c == '\'' {
					if quoted == 0 {
						quoted = c
					} else {
						quoted = 0
					}
					continue
				} else if c == ' ' {
					if quoted == 0 {
						args = append(args, arg)
						arg = ""
						continue
					}
				}
				arg += string(c)
			}

			args = append(args, arg)
			return callback(args)
		})

	}
	return doc
}
