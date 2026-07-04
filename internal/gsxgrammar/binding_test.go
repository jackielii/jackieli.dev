package gsxgrammar_test

import (
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/jackielii/jackieli.dev/internal/gsxgrammar"
)

func TestParsesGSX(t *testing.T) {
	lang := tree_sitter.NewLanguage(gsxgrammar.Language())
	if lang == nil {
		t.Fatal("nil language")
	}
	p := tree_sitter.NewParser()
	defer p.Close()
	if err := p.SetLanguage(lang); err != nil {
		t.Fatalf("SetLanguage: %v (likely tree-sitter ABI mismatch — bump go-tree-sitter or regenerate parser)", err)
	}
	src := []byte("component Card(title string) {\n\t<h2>{title}</h2>\n}\n")
	tree := p.Parse(src, nil)
	defer tree.Close()
	root := tree.RootNode()
	if root.HasError() {
		t.Fatalf("parse produced ERROR node; sexp:\n%s", root.ToSexp())
	}
	if got := root.Kind(); got != "source_file" && got != "document" {
		t.Logf("root kind = %q (informational)", got)
	}
}
