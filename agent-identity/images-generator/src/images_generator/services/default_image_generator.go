package services

import (
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/images-generator/src/images_generator/settings"
)

type DefaultImageGeneratorInterface interface {
	Generate(params GenerateImageParams) ([]byte, error)
}

type DefaultImageGenerator struct {
	agentSVGImage                AgentSVGImageInterface
	agentImageGenerationSettings *settings.ImageGenerationSettings
}

func NewDefaultImageGenerator(
	agentSVGImage AgentSVGImageInterface,
	agentImageGenerationSettings *settings.ImageGenerationSettings,
) *DefaultImageGenerator {
	return &DefaultImageGenerator{
		agentSVGImage:                agentSVGImage,
		agentImageGenerationSettings: agentImageGenerationSettings,
	}
}

func (g *DefaultImageGenerator) Generate(params GenerateImageParams) ([]byte, error) {
	var sb strings.Builder

	borderGradient := g.agentSVGImage.RenderBorderGradient(params.Creator)
	imageDataURI := g.agentSVGImage.RenderImageData(params.BaseImageBytes, params.BaseImageMimeType)

	sb.WriteString(`<svg width="512" height="768" viewBox="0 0 512 768" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`)
	sb.WriteString(`<defs>`)
	sb.WriteString(`<clipPath id="clip-border"><rect width="512" height="768" rx="20" ry="20" /></clipPath>`)
	sb.WriteString(`<clipPath id="clip-image"><rect x="20" y="20" width="472" height="728" rx="16" ry="16" /></clipPath>`)
	sb.WriteString(`<clipPath id="wallet-carousel-clip"><rect x="7" y="750" width="498" height="24"/></clipPath>`)
	sb.WriteString(`<linearGradient id="fadeOverlay" x1="0" y1="0" x2="1" y2="1">`)
	sb.WriteString(`<stop offset="0%" stop-color="black" stop-opacity="1"/>`)
	sb.WriteString(`<stop offset="10%" stop-color="black" stop-opacity="0"/>`)
	sb.WriteString(`<stop offset="90%" stop-color="black" stop-opacity="0"/>`)
	sb.WriteString(`<stop offset="100%" stop-color="black" stop-opacity="1"/>`)
	sb.WriteString(`</linearGradient>`)
	sb.WriteString(borderGradient)
	sb.WriteString(`<path id="walletCarouselPath" d="M 0,762 H 512" />`)
	sb.WriteString(`</defs>`)

	sb.WriteString(`<g clip-path="url(#clip-border)">`)
	sb.WriteString(`<rect width="512" height="768" fill="url(#borderGradient)" />`)
	sb.WriteString(`<g clip-path="url(#clip-image)">`)
	sb.WriteString(fmt.Sprintf(
		`<image width="472" height="728" x="20" y="20" xlink:href="%s" preserveAspectRatio="xMidYMid slice" />`,
		imageDataURI,
	))
	sb.WriteString(`<rect x="20" y="20" width="472" height="728" fill="url(#fadeOverlay)" />`)
	sb.WriteString(`</g></g>`)

	sb.WriteString(`<g>`)
	sb.WriteString(g.renderNameAndEnsDomain(params.Name, params.EnsDomain))
	sb.WriteString(g.agentSVGImage.RenderAgentInformation(params))
	sb.WriteString(g.agentSVGImage.RenderStars(params.Stars))
	sb.WriteString(g.renderWalletCarousel(params.Wallets))

	if g.agentImageGenerationSettings.ImageMarkPath != nil {
		mark, err := g.agentSVGImage.RenderMark()
		if err != nil {
			return nil, fmt.Errorf("failed to add image mark: %w", err)
		}
		sb.WriteString(mark)
	}

	sb.WriteString(`</g></svg>`)

	imageBytes := []byte(sb.String())
	return imageBytes, nil
}

func (a *DefaultImageGenerator) renderNameAndEnsDomain(name string, ensDomain string) string {
	maxWidth := 440
	nameFontSize := float32(48)
	ensFontSize := float32(16)

	padding := 12

	backgroundWith := maxWidth
	backgroundHeight := 96

	needCarouselName := true
	nameLength := a.agentSVGImage.EstimateStrLen(name, nameFontSize)
	if nameLength < maxWidth {
		needCarouselName = false
		backgroundWith = nameLength
	}

	needCarouselEns := true
	ensLength := a.agentSVGImage.EstimateStrLen(ensDomain, ensFontSize)
	if ensLength < maxWidth {
		needCarouselEns = false
		if ensLength > backgroundWith {
			backgroundWith = ensLength
		}
	}

	backgroundWith = backgroundWith + 2*padding
	x := (512 - backgroundWith) / 2
	y := 40

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<g transform="translate(%d,%d)">`, x, y))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" rx="12" ry="12" fill="rgba(0,0,0,0.7)" />`, backgroundWith, backgroundHeight))
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d">`, backgroundWith, backgroundHeight))
	sb.WriteString(`<defs>`)
	if needCarouselName {
		nameY := 48
		nameClipY := nameY - int(nameFontSize*0.75)
		sb.WriteString(fmt.Sprintf(`<clipPath id="name-carousel-clip"><rect x="%d" y="%d" width="%d" height="%f"/></clipPath>`, padding, nameClipY, backgroundWith-2*padding, nameFontSize+8))
		sb.WriteString(fmt.Sprintf(`<path id="nameCarouselPath" d="M %d,%d H %d" />`, padding, nameY, backgroundWith-padding))
	}
	if needCarouselEns {
		ensY := 80
		ensClipY := ensY - int(ensFontSize*0.75)
		sb.WriteString(fmt.Sprintf(`<clipPath id="ens-carousel-clip"><rect x="%d" y="%d" width="%d" height="%f"/></clipPath>`, padding, ensClipY, backgroundWith-2*padding, ensFontSize+6))
		sb.WriteString(fmt.Sprintf(`<path id="ensCarouselPath" d="M %d,%d H %d" />`, padding, ensY, backgroundWith-padding))
	}
	sb.WriteString(`</defs>`)

	if needCarouselName {
		sb.WriteString(a.renderCarousel([]string{name}, 20, nameFontSize, backgroundWith-2*padding, "name-carousel-clip", "nameCarouselPath", false))
	} else {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="white" font-size="%d" font-family="Courier, Liberation Mono, monospace">%s</text>`, backgroundWith/2, 48, int(nameFontSize), name))
	}

	if needCarouselEns {
		sb.WriteString(a.renderCarousel([]string{ensDomain}, 20, ensFontSize, backgroundWith, "ens-carousel-clip", "ensCarouselPath", true))
	} else {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="white" font-size="%d" font-family="Courier, Liberation Mono, monospace">%s</text>`, backgroundWith/2, 80, int(ensFontSize), ensDomain))
	}

	sb.WriteString(`</svg></g>`)

	return sb.String()
}

func (a *DefaultImageGenerator) renderWalletCarousel(wallets []string) string {
	return a.renderCarousel(wallets, 10, 14, 512, "wallet-carousel-clip", "walletCarouselPath", false)
}

func (a *DefaultImageGenerator) renderCarousel(values []string, durationByValue uint, fontSize float32, containerLength int, clipPathId, carouselPathId string, rightToLeft bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<g clip-path="url(#%s)">`, clipPathId))
	sb.WriteString(fmt.Sprintf(`<text fill="white" font-family="Courier, Liberation Mono, monospace" font-size="%.f">`, fontSize))

	baseFirstValueLength := a.agentSVGImage.EstimateStrLen(values[0], fontSize)
	spacing := baseFirstValueLength + containerLength
	totalDistance := baseFirstValueLength + spacing*(len(values)-1) + containerLength
	duration := durationByValue * uint(len(values))

	for i, value := range values {
		var offsetStart, offsetEnd int
		if rightToLeft {
			// Start from right, move to left
			offsetStart = containerLength + spacing*i
			offsetEnd = -(baseFirstValueLength)
		} else {
			// Default: left to right
			offsetStart = -(baseFirstValueLength + spacing*i)
			offsetEnd = totalDistance + offsetStart
		}
		sb.WriteString(fmt.Sprintf(
			`<textPath xlink:href="#%s" startOffset="%d">%s<animate attributeName="startOffset" from="%d" to="%d" dur="%ds" repeatCount="indefinite" /></textPath>`,
			carouselPathId, offsetStart, value, offsetStart, offsetEnd, duration,
		))
	}

	sb.WriteString(`</text>`)
	sb.WriteString(`</g>`)
	return sb.String()
}
