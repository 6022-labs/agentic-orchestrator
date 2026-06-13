package images

import (
	"fmt"
	"math/rand"
)

type AbstractImageGenerator struct {
}

func NewAbstractImageGenerator() *AbstractImageGenerator {
	return &AbstractImageGenerator{}
}

// Deterministic from creator so the SVG card and rasterized icon share a border.
func (AbstractImageGenerator) borderGradientColors(creator string) []string {
	seed := int64(0)
	for _, char := range creator {
		seed += int64(char)
	}
	randSeed := rand.New(rand.NewSource(seed))
	nColors := randSeed.Intn(8) + 3
	colors := make([]string, nColors)
	for i := 0; i < nColors; i++ {
		colors[i] = fmt.Sprintf("#%06X", randSeed.Intn(0xFFFFFF))
	}
	return colors
}
