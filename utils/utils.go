package utils

import (
	"archive/zip"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// SuffixMultiplier struct
type SuffixMultiplier struct {
	Suffix     string
	Multiplier float64
}

var sizeSuffixes []SuffixMultiplier = []SuffixMultiplier{
	{"kib", 1024}, {"kb", 1024}, {"ki", 1024}, {"k", 1024},
	{"mib", 1024 * 1024}, {"mb", 1024 * 1024}, {"mi", 1024 * 1024}, {"m", 1024 * 1024},
	{"gib", 1024 * 1024 * 1024}, {"gb", 1024 * 1024 * 1024}, {"gi", 1024 * 1024 * 1024}, {"g", 1024 * 1024 * 1024},
	{"tib", 1024 * 1024 * 1024 * 1024}, {"tb", 1024 * 1024 * 1024 * 1024}, {"ti", 1024 * 1024 * 1024 * 1024}, {"t", 1024 * 1024 * 1024 * 1024},
	{"pib", 1024 * 1024 * 1024 * 1024 * 1024}, {"pb", 1024 * 1024 * 1024 * 1024 * 1024}, {"pi", 1024 * 1024 * 1024 * 1024 * 1024}, {"p", 1024 * 1024 * 1024 * 1024 * 1024},
}

var timeSuffixes []SuffixMultiplier = []SuffixMultiplier{
	{"microseconds", float64(time.Microsecond)}, {"microsecond", float64(time.Microsecond)},
	{"milliseconds", float64(time.Millisecond)}, {"millisecond", float64(time.Millisecond)},
	{"minutes", float64(time.Minute)}, {"minute", float64(time.Minute)},
	{"hours", float64(time.Hour)}, {"hour", float64(time.Hour)},
	{"days", float64(24 * time.Hour)}, {"day", float64(24 * time.Hour)},
	{"seconds", float64(time.Second)}, {"second", float64(time.Second)},
	{"mks", float64(time.Microsecond)}, {"ms", float64(time.Millisecond)},
	{"m", float64(time.Minute)}, {"h", float64(time.Hour)}, {"d", float64(24 * time.Hour)}, {"s", float64(time.Second)},
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
			multiplier = v.Multiplier
			break
		}
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return int64(v * multiplier), nil
}

// ParseSizeString func
func ParseSizeString(size string) (int64, error) {
	return ParseSuffixed(size, sizeSuffixes)
}

// ParseTimeDuration func
func ParseTimeDuration(duration string) (time.Duration, error) {
	v, err := ParseSuffixed(duration, timeSuffixes)
	return time.Duration(v), err
}

// FileToArchive struct
type FileToArchive struct {
	Name, Path string
}

// ZipFilesToWriter func
func ZipFilesToWriter(w *zip.Writer, files []FileToArchive) error {
	for _, file := range files {
		err := func() error {
			src, err := os.Open(file.Path)
			if err != nil {
				return err
			}
			defer src.Close()

			dest, err := w.Create(file.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(dest, src); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

// ZipFilesToFile func
func ZipFilesToFile(zipFile string, perm os.FileMode, files []FileToArchive) error {
	f, err := os.OpenFile(zipFile, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	err = func() error {
		defer f.Close()
		zw := zip.NewWriter(f)
		err := ZipFilesToWriter(zw, files)
		errClose := zw.Close()
		if err != nil {
			return err
		}
		return errClose
	}()
	if err != nil {
		os.Remove(zipFile)
	}
	return err
}

// WriteStringln func
func WriteStringln(sb *strings.Builder, s string) {
	if sb.Len() > 0 {
		sb.WriteString(NewLine)
	}
	sb.WriteString(s)
}
