package misc

import "testing"
func TestMustGetFileExt(t *testing.T) {
    tests := []struct {
        path     string
        expected string
        err      bool
    }{
        {"file.txt", "txt", false},
        {"path/to/file.jpg", "jpg", false},
        {"no_extension", "", true},
        {"", "", true},
    }

    for _, test := range tests {
        ext, err := MustGetFileExt(test.path)
        if test.err && err == nil {
            t.Errorf("Expected error for path %s, but got nil", test.path)
        } else if !test.err && err != nil {
            t.Errorf("Unexpected error for path %s: %v", test.path, err)
        }
        if ext != test.expected {
            t.Errorf("Expected extension %s for path %s, but got %s", test.expected, test.path, ext)
        }
    }
}
