# Vendored tree-sitter-gsx grammar

C sources and headers here are copied verbatim from `../gsxhq/tree-sitter-gsx`.
`gsx` has an external scanner (`scanner.c`), so both `parser.c` and `scanner.c`
are required.

## Regenerate after a grammar change

```bash
# in ../gsxhq/tree-sitter-gsx
npx tree-sitter generate
# back in this repo
cp ../gsxhq/tree-sitter-gsx/src/parser.c        internal/gsxgrammar/parser.c
cp ../gsxhq/tree-sitter-gsx/src/scanner.c       internal/gsxgrammar/scanner.c
cp ../gsxhq/tree-sitter-gsx/src/tree_sitter/*.h internal/gsxgrammar/tree_sitter/
cp ../gsxhq/tree-sitter-gsx/queries/*.scm       internal/gsxgrammar/queries/gsx/
```

Then run `go test ./...`.
