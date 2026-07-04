package gsxhl

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ProcessDir highlights every *.html under root in place, returning the number
// of files changed.
func (h *Highlighter) ProcessDir(root string) (int, error) {
	changed := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		in := string(b)
		if !strings.Contains(in, "gsx-hl") {
			return nil
		}
		out, err := h.ProcessHTML(in)
		if err != nil {
			return err
		}
		if out == in {
			return nil
		}
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			return err
		}
		changed++
		return nil
	})
	return changed, err
}
