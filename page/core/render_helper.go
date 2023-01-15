package core

import (
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func FormatMultilineText(text string) (content []vecty.MarkupOrChild) {
	for _, line := range strings.FieldsFunc(text, func(r rune) bool { return r == '\n' }) {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			if len(content) > 0 {
				content = append(content, elem.Break())
			}
			content = append(content, vecty.Text(line))
		}
	}
	return
}
