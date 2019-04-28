package character

import (
	"bytes"
	"fmt"
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
	t, err := template.New("character").Funcs(template.FuncMap{
		"sign": sign,
	}).Parse(string(b))
	if err != nil {
		return err
	}

	t, err = t.New("extra").Parse(string(markdown.ToHTML([]byte(c.rawMD), nil, nil)))
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	t.ExecuteTemplate(buf, "extra", c)

	return t.ExecuteTemplate(wr, "character", tpl{
		Character: c,
		Extra:     template.HTML(buf.Bytes()),
	})
}

func (c *Character) SkillTable() template.HTML {
	table := tag("table")()
	tr := tag("tr")()
	td := tag("td")

	tbody := ""
	for _, s := range c.Skills() {
		bonusSign := "zero"
		if s.Bonus > 0 {
			bonusSign = "positive"
		} else if s.Bonus < 0 {
			bonusSign = "negative"
		}
		bonusClass := fmt.Sprintf(`class="bonus %s"`, bonusSign)
		tbody += tr(
			td(`class="prof"`)(s.Prof),
			td(`class="mod"`)(s.Mod),
			td(`class="skill"`)(s.Skill),
			td(bonusClass)(sign(s.Bonus)),
		)
	}

	return template.HTML(table(tbody))
}

func tag(name string) func(attributes ...string) func(children ...interface{}) string {
	return func(attributes ...string) func(children ...interface{}) string {
		return func(children ...interface{}) string {
			attrs := ""
			for _, attr := range attributes {
				attrs += " " + attr
			}
			return fmt.Sprintf("<%s%s>%s</%s>", name, attrs, fmt.Sprint(children...), name)
		}
	}
}

func sign(i int) template.HTML {
	s := ""
	span := tag("span")
	spanSign := span(`class="sign"`)
	if i >= 0 {
		s = spanSign("+")
	} else {
		i = i * -1
		s = spanSign("-")
	}
	return template.HTML(fmt.Sprintf("%s%s", s, span(`class="value"`)(i)))
}
