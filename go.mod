module github.com/jackielii/jackieli.dev

go 1.26.1

require (
	github.com/tree-sitter/go-tree-sitter v0.25.0
	github.com/tree-sitter/tree-sitter-css v0.25.0
	github.com/tree-sitter/tree-sitter-go v0.25.0
	github.com/tree-sitter/tree-sitter-javascript v0.25.0
	go.gopad.dev/go-tree-sitter-highlight v0.0.0-20241203223050-3ffb64c3a650
)

require github.com/mattn/go-pointer v0.0.1 // indirect

// Local fork of go.gopad.dev/go-tree-sitter-highlight patching several
// injection-handling bugs (property-only `#set! injection.language` never
// resolving, inverted layer sort-order comparisons, and a queue-dequeue bug
// that hung forever on <script>/<style> injections) that otherwise prevent
// gsx's go/javascript/css injections from working at all. See the PATCHED
// comments in internal/thirdparty/gopadhighlight/*.go for details; no
// upstream fix exists as of the pinned commit.
replace go.gopad.dev/go-tree-sitter-highlight => ./internal/thirdparty/gopadhighlight
