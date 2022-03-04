package strutil

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

// SuffixMultiplier struct
type SuffixMultiplier struct {
	Suffix     string
	Multiplier int64
}

var sizeSuffixesParse []SuffixMultiplier = []SuffixMultiplier{
	{"kib", 1024}, {"kb", 1024}, {"ki", 1024}, {"k", 1024},
	{"mib", 1024 * 1024}, {"mb", 1024 * 1024}, {"mi", 1024 * 1024}, {"m", 1024 * 1024},
	{"gib", 1024 * 1024 * 1024}, {"gb", 1024 * 1024 * 1024}, {"gi", 1024 * 1024 * 1024}, {"g", 1024 * 1024 * 1024},
	{"tib", 1024 * 1024 * 1024 * 1024}, {"tb", 1024 * 1024 * 1024 * 1024}, {"ti", 1024 * 1024 * 1024 * 1024}, {"t", 1024 * 1024 * 1024 * 1024},
	{"pib", 1024 * 1024 * 1024 * 1024 * 1024}, {"pb", 1024 * 1024 * 1024 * 1024 * 1024}, {"pi", 1024 * 1024 * 1024 * 1024 * 1024}, {"p", 1024 * 1024 * 1024 * 1024 * 1024},
}

var sizeSuffixesFormat []SuffixMultiplier = []SuffixMultiplier{
	{" PiB", 1024 * 1024 * 1024 * 1024 * 1024},
	{" TiB", 1024 * 1024 * 1024 * 1024},
	{" GiB", 1024 * 1024 * 1024},
	{" MiB", 1024 * 1024},
	{" KiB", 1024},
}

var timeSuffixesParse []SuffixMultiplier = []SuffixMultiplier{
	{"microseconds", int64(time.Microsecond)}, {"microsecond", int64(time.Microsecond)},
	{"milliseconds", int64(time.Millisecond)}, {"millisecond", int64(time.Millisecond)},
	{"minutes", int64(time.Minute)}, {"minute", int64(time.Minute)},
	{"hours", int64(time.Hour)}, {"hour", int64(time.Hour)},
	{"days", 24 * int64(time.Hour)}, {"day", 24 * int64(time.Hour)},
	{"seconds", int64(time.Second)}, {"second", int64(time.Second)},
	{"mks", int64(time.Microsecond)}, {"ms", int64(time.Millisecond)},
	{"m", int64(time.Minute)}, {"h", int64(time.Hour)}, {"d", int64(24 * time.Hour)}, {"s", int64(time.Second)},
}

var timeSuffixesFormat []SuffixMultiplier = []SuffixMultiplier{
	{" d", 24 * int64(time.Hour)},
	{" h", int64(time.Hour)},
	{" m", int64(time.Minute)},
	{" s", int64(time.Second)},
	{" ms", int64(time.Millisecond)},
	{" mks", int64(time.Microsecond)},
}

// ParseSuffixed func
func ParseSuffixed(value string, suffixes []SuffixMultiplier) (int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	value = strings.ToLower(value)

	var multiplier float64 = 1
	for _, v := range suffixes {
		if strings.HasSuffix(value, v.Suffix) {
			value = strings.TrimSpace(value[0 : len(value)-len(v.Suffix)])
			multiplier = float64(v.Multiplier)
			break
		}
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return int64(v * multiplier), nil
}

// FormatSuffixed
func FormatSuffixed(value int64, suffixes []SuffixMultiplier) string {
	if value == 0 {
		return "0"
	}
	suffix := ""
	for _, sfx := range suffixes {
		if value%sfx.Multiplier == 0 {
			value /= sfx.Multiplier
			suffix = sfx.Suffix
			break
		}
	}
	return strconv.FormatInt(value, 10) + suffix
}

// ParseSizeString func
func ParseSizeString(size string) (int64, error) {
	return ParseSuffixed(size, sizeSuffixesParse)
}

// FormatSizeString func
func FormatSizeString(value int64) string {
	return FormatSuffixed(value, sizeSuffixesFormat)
}

// ParseTimeDuration func
func ParseTimeDuration(duration string) (time.Duration, error) {
	v, err := ParseSuffixed(duration, timeSuffixesParse)
	return time.Duration(v), err
}

// FormatTimeDuration func
func FormatTimeDuration(value time.Duration) string {
	return FormatSuffixed(int64(value), timeSuffixesFormat)
}

// ParseOptions func
func ParseOptions(options string, accept func(option string) bool) bool {
	for _, option := range strings.FieldsFunc(options, func(r rune) bool { return r == ',' || r == ';' || unicode.IsSpace(r) }) {
		if option != "" {
			if !accept(option) {
				return false
			}
		}
	}
	return true
}
