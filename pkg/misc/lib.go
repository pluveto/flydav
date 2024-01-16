package misc

import (
	"fmt"
	"path/filepath"
	"strings"
)

func MustGetFileExt(path string) (string, error) {
    ext := filepath.Ext(path)
    if ext == "" || ext == "." {
        return "", fmt.Errorf("invalid file path " + path)
    }
    return strings.TrimPrefix(ext, "."), nil
}

