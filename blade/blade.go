package blade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

type Parser struct {
	directives map[string]func(args Args) string
}

func New() *Parser {
	return &Parser{
		directives: map[string]func(args Args) string{},
	}
}

func (b *Parser) Directive(name string, callback func(args Args) string) {
	b.directives[name] = callback
}

func (b *Parser) Parse(doc string) string {
	bDoc := []byte(doc)
	for name, callback := range b.directives {
		re := regexp.MustCompile(fmt.Sprintf(`\@%s\(([^\)]*)\)`, regexp.QuoteMeta(name)))
		bDoc = re.ReplaceAllFunc(bDoc, func(s []byte) []byte {
			return []byte(callback(Args(re.FindSubmatch(s)[1])))
		})

	}
	return string(bDoc)
}

type Args []byte

func (a Args) Unmarshal(args ...interface{}) error {
	d := json.NewDecoder(bytes.NewReader([]byte(a)))
	for _, arg := range args {
		err := d.Decode(arg)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}
