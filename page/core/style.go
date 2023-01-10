package core

import "strings"

var styles strings.Builder

func AppendStyles(s string) {
	styles.WriteString(s)
}

func Stylesheet() string {
	return styles.String()
}
