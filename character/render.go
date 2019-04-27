package character

import (
	"html/template"
	"io"

	"github.com/gomarkdown/markdown"
)

//go:generate go-bindata -pkg character character.html

type tpl struct {
	*Character

	Extra template.HTML
}

func (c *Character) Render(wr io.Writer) error {
	b := MustAsset("character.html")
	t, err := template.New("character").Parse(string(b))
	if err != nil {
		return err
	}

	return t.Execute(wr, tpl{
		Character: c,
		Extra:     template.HTML(markdown.ToHTML([]byte(c.rawMD), nil, nil)),
	})
}
