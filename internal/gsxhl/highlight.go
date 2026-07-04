// Package gsxhl highlights gsx source to HTML using tree-sitter.
package gsxhl

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
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

// New builds a Highlighter for gsx.
func New() (*Highlighter, error) {
	gsxLang := tree_sitter.NewLanguage(gsxgrammar.Language())
	cfg, err := highlight.NewConfiguration(
		gsxLang, "gsx",
		gsxgrammar.GSXHighlights, gsxgrammar.GSXInjections, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("gsx configuration: %w", err)
	}
	cfg.Configure(recognizedNames)
	return &Highlighter{
		inner:  highlight.New(),
		gsx:    cfg,
		inject: func(string) *highlight.Configuration { return nil },
	}, nil
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
