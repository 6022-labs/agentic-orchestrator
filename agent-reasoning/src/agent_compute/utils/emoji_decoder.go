package utils

import (
	"sort"
	"strings"

	"github.com/kyokomi/emoji/v2"
)

var demojizeReplacer = buildDemojizeReplacer()

// buildDemojizeReplacer compiles a strings.Replacer that maps every unicode
// emoji glyph in kyokomi's table back to its canonical shortcode (the LLM
// emit form). Sort by glyph length descending so multi-codepoint sequences
// (e.g. ZWJ-joined emoji) win over their sub-glyphs.
func buildDemojizeReplacer() *strings.Replacer {
	revMap := emoji.RevCodeMap()

	type pair struct {
		unicode   string
		shortcode string
	}
	pairs := make([]pair, 0, len(revMap))
	for glyph, aliases := range revMap {
		if len(aliases) == 0 {
			continue
		}
		pairs = append(pairs, pair{unicode: glyph, shortcode: aliases[0]})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if len(pairs[i].unicode) != len(pairs[j].unicode) {
			return len(pairs[i].unicode) > len(pairs[j].unicode)
		}
		return pairs[i].unicode < pairs[j].unicode
	})

	flat := make([]string, 0, len(pairs)*2)
	for _, p := range pairs {
		flat = append(flat, p.unicode, p.shortcode)
	}
	return strings.NewReplacer(flat...)
}

func DemojizeContent(content string) string {
	return demojizeReplacer.Replace(content)
}
