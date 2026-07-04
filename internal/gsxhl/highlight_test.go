package gsxhl_test

import (
	"regexp"
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

var emptySpanRE = regexp.MustCompile(`<span[^>]*></span>`)

// TestHighlightHTML_combinedInjections covers the combined-injection path
// (<script> -> JS, <style> -> CSS), which is created up-front rather than per
// node and previously either dropped the second injection entirely or emitted
// its tokens as zero-width spans at the document end.
func TestHighlightHTML_combinedInjections(t *testing.T) {
	h, err := gsxhl.New()
	if err != nil {
		t.Fatal(err)
	}
	src := "component Card() {\n" +
		"\t<script>const n = document.querySelector(\"h2\")</script>\n" +
		"\t<style>.card { color: red; }</style>\n" +
		"}"
	out, err := h.HighlightHTML([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	// JS injected into <script>: `const` is a JS keyword.
	if !strings.Contains(out, `<span class="ts-keyword">const</span>`) {
		t.Errorf("embedded JavaScript not highlighted (script injection inactive)\n---\n%s", out)
	}
	// CSS injected into <style>: `color` is a CSS property. This is the second
	// combined injection — the one the layer-array aliasing bug used to destroy.
	if !strings.Contains(out, `<span class="ts-property">color</span>`) {
		t.Errorf("embedded CSS not highlighted (style injection inactive/dropped)\n---\n%s", out)
	}
	// Hard requirement: no empty-span artifacts anywhere in the output.
	if m := emptySpanRE.FindAllString(out, -1); len(m) > 0 {
		t.Errorf("output contains %d empty-span artifact(s): %v\n---\n%s", len(m), m, out)
	}
}

// TestHighlightHTML_noEmptySpans guards the newline-boundary case (gsx
// punctuation.bracket captures that start with a newline) that the upstream
// renderer turned into empty <span></span> pairs.
func TestHighlightHTML_noEmptySpans(t *testing.T) {
	h, err := gsxhl.New()
	if err != nil {
		t.Fatal(err)
	}
	src := "component Card() {\n\t<section>\n\t\t<h2>hi</h2>\n\t</section>\n}"
	out, err := h.HighlightHTML([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if m := emptySpanRE.FindAllString(out, -1); len(m) > 0 {
		t.Errorf("multi-line gsx output contains %d empty-span artifact(s): %v\n---\n%s", len(m), m, out)
	}
}
