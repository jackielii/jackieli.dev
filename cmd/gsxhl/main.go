// Command gsxhl highlights gsx code blocks in a directory of built HTML,
// in place. Intended to run after `hugo` over ./public.
package main

import (
	"fmt"
	"os"

	"github.com/jackielii/jackieli.dev/internal/gsxhl"
)

func main() {
	root := "./public"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	h, err := gsxhl.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gsxhl:", err)
		os.Exit(1)
	}
	changed, err := h.ProcessDir(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gsxhl:", err)
		os.Exit(1)
	}
	fmt.Printf("gsxhl: highlighted gsx in %d file(s) under %s\n", changed, root)
}
