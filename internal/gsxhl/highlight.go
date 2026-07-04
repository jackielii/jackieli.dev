// Package gsxhl highlights gsx source to HTML using tree-sitter.
package gsxhl

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
	"unsafe"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_css "github.com/tree-sitter/tree-sitter-css/bindings/go"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	highlight "go.gopad.dev/go-tree-sitter-highlight"

	"github.com/jackielii/jackieli.dev/internal/gsxgrammar"
)

// recognizedNames is the fixed capture set. Index = Highlight value returned by
// the tree-sitter highlighter; the AttributeCallback maps it back to a class.
// Covers gsx captures plus captures produced by injected go/js/css queries.
var recognizedNames = []string{
	"attribute", "comment", "constant", "constant.builtin", "constructor",
	"embedded", "function", "function.builtin", "keyword", "module", "number",
	"operator", "property", "punctuation.bracket", "punctuation.delimiter",
	"punctuation.special", "string", "string.special", "tag", "type",
	"type.builtin", "variable", "variable.builtin",
}

func className(capture string) string {
	return "ts-" + strings.ReplaceAll(capture, ".", "-")
}

// Highlighter holds the configured tree-sitter highlight Configuration(s).
type Highlighter struct {
	inner *highlight.Highlighter
	gsx   *highlight.Configuration
	// injected configs added in Task 3; nil-safe here.
	inject func(languageName string) *highlight.Configuration
}

// newConfig builds a highlight.Configuration for a language identified by its
// tree-sitter Language() pointer, name, and highlights query.
func newConfig(langPtr unsafe.Pointer, name string, highlights []byte) (*highlight.Configuration, error) {
	lang := tree_sitter.NewLanguage(langPtr)
	cfg, err := highlight.NewConfiguration(lang, name, highlights, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%s configuration: %w", name, err)
	}
	cfg.Configure(recognizedNames)
	return cfg, nil
}

// New builds a Highlighter for gsx, wired to inject go/javascript/css
// highlighting into the embedded regions gsx's injections.scm identifies.
func New() (*Highlighter, error) {
	gsxCfg, err := highlight.NewConfiguration(
		tree_sitter.NewLanguage(gsxgrammar.Language()), "gsx",
		gsxgrammar.GSXHighlights, gsxgrammar.GSXInjections, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("gsx configuration: %w", err)
	}
	gsxCfg.Configure(recognizedNames)

	goCfg, err := newConfig(tree_sitter_go.Language(), "go", gsxgrammar.GoHighlights)
	if err != nil {
		return nil, err
	}
	jsCfg, err := newConfig(tree_sitter_javascript.Language(), "javascript", gsxgrammar.JSHighlights)
	if err != nil {
		return nil, err
	}
	cssCfg, err := newConfig(tree_sitter_css.Language(), "css", gsxgrammar.CSSHighlights)
	if err != nil {
		return nil, err
	}

	inject := func(name string) *highlight.Configuration {
		switch name {
		case "go":
			return goCfg
		case "javascript":
			return jsCfg
		case "css":
			return cssCfg
		default:
			return nil
		}
	}
	return &Highlighter{inner: highlight.New(), gsx: gsxCfg, inject: inject}, nil
}

// HighlightHTML returns the highlighted inner HTML for source (spans only, no
// wrapping <pre>/<code>).
//
// It renders the highlight event stream directly rather than using the
// library's HTMLRender. That renderer closes and reopens every open span at
// each newline (a feature for line-numbered / per-line-wrapped output), which
// emits an empty `<span class="…"></span>` whenever a capture's text *starts*
// with a newline — and gsx's `punctuation.bracket` captures routinely do
// (e.g. the `\n\t<` before an element's opening `<`). gsx highlighting is
// inline and never wraps individual lines, so we deliberately do not split
// spans at newlines: a highlighted span may span a line break, producing
// clean output with no empty-span artifacts.
func (h *Highlighter) HighlightHTML(source []byte) (string, error) {
	events := h.inner.Highlight(context.Background(), *h.gsx, source, h.inject)

	var buf bytes.Buffer
	for event, err := range events {
		if err != nil {
			return "", fmt.Errorf("highlight: %w", err)
		}
		switch e := event.(type) {
		case highlight.EventCaptureStart:
			buf.WriteString("<span")
			if i := int(e.Highlight); e.Highlight != highlight.DefaultHighlight && i >= 0 && i < len(recognizedNames) {
				buf.WriteString(` class="`)
				buf.WriteString(className(recognizedNames[i]))
				buf.WriteByte('"')
			}
			buf.WriteByte('>')
		case highlight.EventCaptureEnd:
			buf.WriteString("</span>")
		case highlight.EventSource:
			writeEscapedHTML(&buf, source[e.StartByte:e.EndByte])
		}
		// EventLayerStart / EventLayerEnd carry no rendered output here; the
		// attribute we emit does not depend on the injected language name.
	}
	return buf.String(), nil
}

// writeEscapedHTML writes src to buf, HTML-escaping the characters that matter
// inside element content, dropping carriage returns and invalid UTF-8. It uses
// the same entity set and skip rules as the upstream HTMLRender, minus that
// renderer's newline span-splitting. Newlines are written through verbatim.
func writeEscapedHTML(buf *bytes.Buffer, src []byte) {
	for len(src) > 0 {
		c, size := utf8.DecodeRune(src)
		src = src[size:]
		if c == utf8.RuneError || c == '\r' {
			continue
		}
		switch c {
		case '&':
			buf.WriteString("&amp;")
		case '\'':
			buf.WriteString("&#39;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&#34;")
		default:
			buf.WriteRune(c)
		}
	}
}
