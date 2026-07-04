// Package gsxhl highlights gsx source to HTML using tree-sitter.
package gsxhl

import (
	"bytes"
	"context"
	"fmt"
	"strings"
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
func (h *Highlighter) HighlightHTML(source []byte) (string, error) {
	events := h.inner.Highlight(context.Background(), *h.gsx, source, h.inject)

	render := highlight.NewHTMLRender()

	var buf bytes.Buffer
	attr := func(hl highlight.Highlight, _ string) []byte {
		if hl == highlight.DefaultHighlight {
			return nil
		}
		i := int(hl)
		if i < 0 || i >= len(recognizedNames) {
			return nil
		}
		return []byte(fmt.Sprintf(`class=%q`, className(recognizedNames[i])))
	}
	if err := render.Render(&buf, events, source, attr); err != nil {
		return "", fmt.Errorf("render: %w", err)
	}
	return buf.String(), nil
}
