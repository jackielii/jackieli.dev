.PHONY: build preview

PREVIEW_PORT ?= 1314

# Build the site and highlight gsx blocks (production-like output).
build:
	hugo --minify
	go run ./cmd/gsxhl ./public

# Build with a local baseURL (so fingerprinted CSS/JS resolve on localhost),
# highlight gsx blocks, and serve the static output so gsx colours show.
# Plain `hugo server` will NOT show gsx colours — the highlighter runs on the
# built output, not the in-memory server. Re-run `make preview` after edits.
preview:
	hugo --minify --baseURL http://localhost:$(PREVIEW_PORT)/
	go run ./cmd/gsxhl ./public
	@echo "Serving highlighted static build at http://localhost:$(PREVIEW_PORT)/ (Ctrl+C to stop)"
	python3 -m http.server -d public $(PREVIEW_PORT)
