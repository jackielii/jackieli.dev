package gsxhl_test

import (
	"strings"
	"testing"

	"github.com/jackielii/jackieli.dev/internal/gsxhl"
)

func TestProcessHTML_replacesBlock(t *testing.T) {
	h, err := gsxhl.New()
	if err != nil {
		t.Fatal(err)
	}
	page := `<article><pre class="gsx-hl"><code>component Card() { &lt;h2&gt;x&lt;/h2&gt; }</code></pre></article>`
	out, err := h.ProcessHTML(page)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "gsx-hl") {
		t.Errorf("marker class not removed:\n%s", out)
	}
	if !strings.Contains(out, `<pre class="chroma">`) {
		t.Errorf("chroma wrapper missing:\n%s", out)
	}
	if !strings.Contains(out, `<span class="ts-keyword">component</span>`) {
		t.Errorf("not highlighted:\n%s", out)
	}
	// idempotent
	out2, err := h.ProcessHTML(out)
	if err != nil {
		t.Fatal(err)
	}
	if out2 != out {
		t.Errorf("ProcessHTML not idempotent")
	}
}

func TestProcessHTML_minifiedUnquoted(t *testing.T) {
	h, _ := gsxhl.New()
	page := `<pre class=gsx-hl><code>component C() {}</code></pre>`
	out, err := h.ProcessHTML(page)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `<span class="ts-keyword">component</span>`) {
		t.Errorf("unquoted-class block not highlighted:\n%s", out)
	}
}
