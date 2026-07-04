.PHONY: build preview

# Build the site and highlight gsx blocks (production-like output).
build:
	hugo --minify --buildFuture
	go run ./cmd/gsxhl ./public

# Build, highlight, and serve the static output locally so gsx colours show.
preview: build
	hugo server --renderToDisk --disableFastRender --buildFuture & \
	sleep 2 && go run ./cmd/gsxhl ./public && \
	echo "Serving highlighted static build; re-run 'make build' after edits." && \
	python3 -m http.server -d public 1313
