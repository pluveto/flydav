package util

import (
	"fmt"
	"math"
	"path"
	"path/filepath"
	"strings"
)

func FormatSize(size int64, base int, space string) string {
	k := 1000
	if base == 2 {
		k = 1024
	}

	if size < 0 {
		return "-" + FormatSize(-size, base, space)
	}
	if size < int64(k) {
		return fmt.Sprintf("%d%sB", size, space)
	}
	if size < int64(math.Pow(float64(k), 2)) {
		return fmt.Sprintf("%.2f%sKB", float64(size)/float64(k), space)
	}
	if size < int64(math.Pow(float64(k), 3)) {
		return fmt.Sprintf("%.2f%sMB", float64(size)/math.Pow(float64(k), 2), space)
	}
	if size < int64(math.Pow(float64(k), 4)) {
		return fmt.Sprintf("%.2f%sGB", float64(size)/math.Pow(float64(k), 3), space)
	}
	if size < int64(math.Pow(float64(k), 5)) {
		return fmt.Sprintf("%.2f%sTB", float64(size)/math.Pow(float64(k), 4), space)
	}
	if size < int64(math.Pow(float64(k), 6)) {
		return fmt.Sprintf("%.2f%sPB", float64(size)/math.Pow(float64(k), 5), space)
	}
	return fmt.Sprintf("%.2f%sEB", float64(size)/math.Pow(float64(k), 6), space)
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	path, err := filepath.Abs(fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/")))
	if err != nil {
		return ""
	}
	return path
}
