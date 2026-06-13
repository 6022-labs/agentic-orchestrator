package images

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/settings"
)

type DefaultImageGenerator struct {
	*AbstractImageGenerator
	imageGenerationSettings *settings.ImageGenerationSettings
}

func NewDefaultImageGenerator(imageGenerationSettings *settings.ImageGenerationSettings) *DefaultImageGenerator {
	return &DefaultImageGenerator{
		AbstractImageGenerator:  NewAbstractImageGenerator(),
		imageGenerationSettings: imageGenerationSettings,
	}
}

func (g *DefaultImageGenerator) Generate(params interfaces.GenerateImageParams) ([]byte, error) {
	var sb strings.Builder

	borderGradient := g.renderBorderGradient(params.Creator)
	imageDataURI := g.renderImageData(params.BaseImageBytes, params.BaseImageMimeType)

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
	sb.WriteString(g.renderAgentInformation(params))
	sb.WriteString(g.renderStars(params.Stars))
	sb.WriteString(g.renderWalletCarousel(params.Wallets))

	if g.imageGenerationSettings.ImageMarkPath != nil {
		mark, err := g.renderMark()
		if err != nil {
			return nil, fmt.Errorf("failed to add image mark: %w", err)
		}
		sb.WriteString(mark)
	}

	sb.WriteString(`</g></svg>`)

	return []byte(sb.String()), nil
}

func (g *DefaultImageGenerator) renderNameAndEnsDomain(name string, ensDomain string) string {
	maxWidth := 440
	nameFontSize := float32(48)
	ensFontSize := float32(16)

	padding := 12

	backgroundWith := maxWidth
	backgroundHeight := 96

	needCarouselName := true
	nameLength := g.estimateStrLen(name, nameFontSize)
	if nameLength < maxWidth {
		needCarouselName = false
		backgroundWith = nameLength
	}

	needCarouselEns := true
	ensLength := g.estimateStrLen(ensDomain, ensFontSize)
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
		sb.WriteString(g.renderCarousel([]string{name}, 20, nameFontSize, backgroundWith-2*padding, "name-carousel-clip", "nameCarouselPath", false))
	} else {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="white" font-size="%d" font-family="Courier, Liberation Mono, monospace">%s</text>`, backgroundWith/2, 48, int(nameFontSize), name))
	}

	if needCarouselEns {
		sb.WriteString(g.renderCarousel([]string{ensDomain}, 20, ensFontSize, backgroundWith, "ens-carousel-clip", "ensCarouselPath", true))
	} else {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="white" font-size="%d" font-family="Courier, Liberation Mono, monospace">%s</text>`, backgroundWith/2, 80, int(ensFontSize), ensDomain))
	}

	sb.WriteString(`</svg></g>`)

	return sb.String()
}

func (g *DefaultImageGenerator) renderWalletCarousel(wallets []string) string {
	return g.renderCarousel(wallets, 10, 14, 512, "wallet-carousel-clip", "walletCarouselPath", false)
}

func (g *DefaultImageGenerator) renderCarousel(values []string, durationByValue uint, fontSize float32, containerLength int, clipPathId, carouselPathId string, rightToLeft bool) string {
	if len(values) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<g clip-path="url(#%s)">`, clipPathId))
	sb.WriteString(fmt.Sprintf(`<text fill="white" font-family="Courier, Liberation Mono, monospace" font-size="%.f">`, fontSize))

	baseFirstValueLength := g.estimateStrLen(values[0], fontSize)
	spacing := baseFirstValueLength + containerLength
	totalDistance := baseFirstValueLength + spacing*(len(values)-1) + containerLength
	duration := durationByValue * uint(len(values))

	for i, value := range values {
		var offsetStart, offsetEnd int
		if rightToLeft {
			offsetStart = containerLength + spacing*i
			offsetEnd = -baseFirstValueLength
		} else {
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

func (g *DefaultImageGenerator) renderBorderGradient(creator string) string {
	colors := g.borderGradientColors(creator)

	var sb strings.Builder
	sb.WriteString(`<linearGradient id="borderGradient" x1="0" y1="1" x2="1" y2="1">`)
	for i, color := range colors {
		offset := (i * 100) / (len(colors) - 1)
		sb.WriteString(fmt.Sprintf(`<stop offset="%d%%" stop-color="%s"/>`, offset, color))
	}
	sb.WriteString(`</linearGradient>`)
	return sb.String()
}

func (g *DefaultImageGenerator) renderImageData(imageBytes []byte, mimeType string) string {
	b64 := base64.StdEncoding.EncodeToString(imageBytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64)
}

func (g *DefaultImageGenerator) renderAgentInformation(p interfaces.GenerateImageParams) string {
	created := fmt.Sprint(p.CreatedAt.Unix() + 6102950400)

	var sb strings.Builder

	fontSize := float32(13)
	padding := 24

	if p.OriginEnsDomain != nil {
		origin := fmt.Sprintf("%s (%d)", *p.OriginEnsDomain, *p.OriginGeneration)

		originEstimatedWidth := g.estimateStrLen(fmt.Sprintf("Origin: %s", origin), fontSize) + padding
		ensOriginNameEstimatedWidth := g.estimateStrLen(*p.OriginEnsDomain, fontSize) + padding

		badgeWith := originEstimatedWidth
		if ensOriginNameEstimatedWidth > originEstimatedWidth {
			badgeWith = ensOriginNameEstimatedWidth
		}

		sb.WriteString(g.renderMultiLineBadge(50, 522, "Origin", origin, *p.OriginEnsDomain, badgeWith))
	}

	sb.WriteString(g.renderBadge(48, 580, "Role", p.Role, g.estimateStrLen(fmt.Sprintf("Role: %s", p.Role), fontSize)+padding))
	sb.WriteString(g.renderBadge(48, 620, "Created", created, g.estimateStrLen(fmt.Sprintf("Created: %s", created), fontSize)+padding))
	sb.WriteString(g.renderBadge(48, 660, "Creator", p.Creator, g.estimateStrLen(fmt.Sprintf("Creator %s", p.Creator), fontSize)+padding))

	return sb.String()
}

func (g *DefaultImageGenerator) renderStars(stars uint8) string {
	if stars == 0 {
		return ""
	}
	if stars > 8 {
		stars = 8
	}

	starString := strings.Repeat("★", int(stars))

	fontSize := 64
	lastBadgeBottom := 660 + 32
	available := 728
	y := ((lastBadgeBottom + available) / 2) + (fontSize / 2) - 4

	return fmt.Sprintf(
		`<text x="256" y="%d" text-anchor="middle" fill="white" font-size="%d" font-family="Courier, Liberation Mono, monospace">%s</text>`,
		y, fontSize, starString,
	)
}

func (g *DefaultImageGenerator) renderMark() (string, error) {
	if g.imageGenerationSettings.ImageMarkPath == nil {
		return "", nil
	}

	markBytes, err := os.ReadFile(*g.imageGenerationSettings.ImageMarkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read image mark file: %w", err)
	}

	markDataURI := g.renderImageData(markBytes, "image/png")

	cx, cy := 460, 56
	r := 22
	imgSize := 40

	return fmt.Sprintf(
		`<g>`+
			`<clipPath id="agentMarkClip"><circle cx="%d" cy="%d" r="%d" /></clipPath>`+
			`<circle cx="%d" cy="%d" r="%d" fill="white" stroke="black" stroke-width="2" />`+
			`<image x="%d" y="%d" width="%d" height="%d" xlink:href="%s" clip-path="url(#agentMarkClip)" />`+
			`</g>`,
		cx, cy, imgSize/2,
		cx, cy, r,
		cx-imgSize/2, cy-imgSize/2, imgSize, imgSize, markDataURI,
	), nil
}

func (g *DefaultImageGenerator) estimateStrLen(value string, fontSize float32) int {
	charWidth := fontSize * 0.6
	return int(float32(len(value)) * charWidth)
}

func (g *DefaultImageGenerator) renderBadge(x, y int, label, value string, width int) string {
	text := fmt.Sprintf("%s: %s", label, value)
	return fmt.Sprintf(
		`<g transform="translate(%d,%d)">`+
			`<rect width="%d" height="32" rx="8" ry="8" fill="rgba(0,0,0,0.6)" />`+
			`<text x="12" y="20" fill="white" font-family="Courier, Liberation Mono, monospace" font-size="13">%s</text>`+
			`</g>`,
		x, y, width, text,
	)
}

func (g *DefaultImageGenerator) renderMultiLineBadge(x, y int, label, value, subtext string, width int) string {
	main := fmt.Sprintf("%s: %s", label, value)
	return fmt.Sprintf(
		`<g transform="translate(%d,%d)">`+
			`<rect width="%d" height="48" rx="8" ry="8" fill="rgba(0,0,0,0.6)" />`+
			`<text x="12" y="18" fill="white" font-family="Courier, Liberation Mono, monospace" font-size="13">%s</text>`+
			`<text x="12" y="35" fill="white" font-family="Courier, Liberation Mono, monospace" font-size="12">%s</text>`+
			`</g>`,
		x, y, width, main, subtext,
	)
}
