package gsxhl_test

import (
	"strings"
	"testing"

	"github.com/jackielii/jackieli.dev/internal/gsxhl"
)

func TestHighlightHTML_gsxStructure(t *testing.T) {
	h, err := gsxhl.New()
	if err != nil {
		t.Fatal(err)
	}
	out, err := h.HighlightHTML([]byte(`component Card() { <h2>x</h2> }`))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<span class="ts-keyword">component</span>`,
		`<span class="ts-function">Card</span>`,
		`<span class="ts-tag">h2</span>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestHighlightHTML_injectsGo(t *testing.T) {
	h, err := gsxhl.New()
	if err != nil {
		t.Fatal(err)
	}
	// `featured` is Go inside a value-form condition; with injections on it
	// should be wrapped in some ts- span rather than left as bare text.
	out, err := h.HighlightHTML([]byte("component C(featured bool) { { if featured { <b>x</b> } } }"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `>featured</span>`) && !strings.Contains(out, `class="ts-`+`variable">featured`) {
		t.Errorf("embedded Go identifier `featured` not highlighted (injection not active)\n---\n%s", out)
	}
}
