package main

import "strings"

func writeStringln(sb *strings.Builder, s string) {
	if sb.Len() > 0 {
		sb.WriteString(NewLine)
	}
	sb.WriteString(s)
}
