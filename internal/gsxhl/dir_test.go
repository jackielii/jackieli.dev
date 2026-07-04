package gsxhl_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackielii/jackieli.dev/internal/gsxhl"
)

func TestProcessDir(t *testing.T) {
	dir := t.TempDir()
	page := filepath.Join(dir, "post", "index.html")
	if err := os.MkdirAll(filepath.Dir(page), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page,
		[]byte(`<pre class="gsx-hl"><code>component C() {}</code></pre>`), 0o644); err != nil {
		t.Fatal(err)
	}
	// a non-gsx file must be left untouched
	other := filepath.Join(dir, "other.html")
	os.WriteFile(other, []byte(`<p>hi</p>`), 0o644)

	h, err := gsxhl.New()
	if err != nil {
		t.Fatal(err)
	}
	changed, err := h.ProcessDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if changed != 1 {
		t.Errorf("changed=%d, want 1", changed)
	}
	got, _ := os.ReadFile(page)
	if !strings.Contains(string(got), `<span class="ts-keyword">component</span>`) {
		t.Errorf("post not highlighted: %s", got)
	}
	got2, _ := os.ReadFile(other)
	if string(got2) != `<p>hi</p>` {
		t.Errorf("non-gsx file modified: %s", got2)
	}
}
