package gsxhl

import (
	"fmt"
	"html"
	"regexp"
)

// gsxBlock matches the render hook's marker block. The class quote is optional
// because `hugo --minify` may drop quotes around the attribute value. `<pre>`
// content is preserved verbatim by the minifier, so the inner capture is exact.
var gsxBlock = regexp.MustCompile(`(?s)<pre class="?gsx-hl"?><code>(.*?)</code></pre>`)

// ProcessHTML replaces every marked gsx block in page with highlighted HTML.
func (h *Highlighter) ProcessHTML(page string) (string, error) {
	var firstErr error
	out := gsxBlock.ReplaceAllStringFunc(page, func(match string) string {
		m := gsxBlock.FindStringSubmatch(match)
		source := []byte(html.UnescapeString(m[1]))
		hl, err := h.HighlightHTML(source)
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("highlight block: %w", err)
			}
			return match // leave unchanged on error
		}
		return `<div class="highlight"><pre class="chroma"><code class="language-gsx">` +
			hl + `</code></pre></div>`
	})
	if firstErr != nil {
		return "", firstErr
	}
	return out, nil
}
