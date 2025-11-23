package services

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/settings"
)

type GenerateImageParams struct {
	Name              string
	EnsDomain         string
	Role              string
	Creator           string
	OriginEnsDomain   *string
	OriginGeneration  *uint
	BaseImageBytes    []byte
	BaseImageMimeType string
	Wallets           []string
	CreatedAt         time.Time
	Stars             uint8 // (max 8)
}

type AgentSVGImageInterface interface {
	RenderBorderGradient(creator string) string
	RenderImageData(imageBytes []byte, mimeType string) string
	RenderAgentInformation(p GenerateImageParams) string
	RenderStars(stars uint8) string
	RenderMark() (string, error)
	EstimateStrLen(value string, charWidth float32) int
}

type AgentSVGImage struct {
	agentImageGenerationSettings *settings.AgentImageGenerationSettings
}

func NewAgentSVGImage(
	agentImageGenerationSettings *settings.AgentImageGenerationSettings,
) *AgentSVGImage {
	return &AgentSVGImage{
		agentImageGenerationSettings: agentImageGenerationSettings,
	}
}

func (a *AgentSVGImage) RenderBorderGradient(creator string) string {
	var sb strings.Builder
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

	sb.WriteString(`<linearGradient id="borderGradient" x1="0" y1="1" x2="1" y2="1">`)
	for i, color := range colors {
		offset := (i * 100) / (len(colors) - 1)
		sb.WriteString(fmt.Sprintf(`<stop offset="%d%%" stop-color="%s"/>`, offset, color))
	}
	sb.WriteString(`</linearGradient>`)
	return sb.String()
}

func (a *AgentSVGImage) RenderImageData(imageBytes []byte, mimeType string) string {
	b64 := base64.StdEncoding.EncodeToString(imageBytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64)
}

func (a *AgentSVGImage) RenderAgentInformation(p GenerateImageParams) string {
	created := fmt.Sprint(p.CreatedAt.Unix() + 6102950400)

	var sb strings.Builder

	fontSize := float32(13)
	padding := 24

	if p.OriginEnsDomain != nil {
		origin := fmt.Sprintf("%s (%d)", *p.OriginEnsDomain, *p.OriginGeneration)

		originEstimatedWidth := a.EstimateStrLen(fmt.Sprintf("Origin: %s", origin), fontSize) + padding
		ensOriginNameEstimatedWidth := a.EstimateStrLen(*p.OriginEnsDomain, fontSize) + padding

		badgeWith := originEstimatedWidth
		if ensOriginNameEstimatedWidth > originEstimatedWidth {
			badgeWith = ensOriginNameEstimatedWidth
		}

		sb.WriteString(a.renderMultiLineBadge(50, 522, "Origin", origin, *p.OriginEnsDomain, badgeWith))
	}

	sb.WriteString(a.renderBadge(48, 580, "Role", p.Role, a.EstimateStrLen(fmt.Sprintf("Role: %s", p.Role), fontSize)+padding))
	sb.WriteString(a.renderBadge(48, 620, "Created", created, a.EstimateStrLen(fmt.Sprintf("Created: %s", created), fontSize)+padding))
	sb.WriteString(a.renderBadge(48, 660, "Creator", p.Creator, a.EstimateStrLen(fmt.Sprintf("Creator %s", p.Creator), fontSize)+padding))

	return sb.String()
}

func (a *AgentSVGImage) RenderStars(stars uint8) string {
	if stars <= 0 {
		return ""
	}
	if stars > 8 {
		stars = 8
	}

	starGlyph := "★"
	starString := strings.Repeat(starGlyph, int(stars))

	fontSize := 64
	lastBadgeBottom := 660 + 32
	available := 728
	y := ((lastBadgeBottom + available) / 2) + (fontSize / 2) - 4

	return fmt.Sprintf(
		`<text x="256" y="%d" text-anchor="middle" fill="white" font-size="%d" font-family="Courier, Liberation Mono, monospace">%s</text>`,
		y, fontSize, starString,
	)
}

func (a *AgentSVGImage) RenderMark() (string, error) {
	if a.agentImageGenerationSettings.ImageMarkPath == nil {
		return "", nil
	}

	markBytes, err := os.ReadFile(*a.agentImageGenerationSettings.ImageMarkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read image mark file: %w", err)
	}

	markDataURI := a.RenderImageData(markBytes, "image/png")

	cx, cy := 460, 56 // center of the circle (x = 448 + 20, y = 36 + 20)
	r := 22           // circle radius slightly larger than the image
	imgSize := 40     // image size

	return fmt.Sprintf(
		`<g>`+
			`<circle cx="%d" cy="%d" r="%d" fill="white" stroke="black" stroke-width="2" />`+
			`<image x="%d" y="%d" width="%d" height="%d" xlink:href="%s" clip-path="circle(50%%)" />`+
			`</g>`,
		cx, cy, r,
		cx-imgSize/2, cy-imgSize/2, imgSize, imgSize, markDataURI,
	), nil
}

func (a *AgentSVGImage) EstimateStrLen(value string, fontSize float32) int {
	textLen := len(value)

	charWidth := fontSize * 0.6
	return int(float32(textLen) * charWidth)
}

func (a *AgentSVGImage) renderBadge(x, y int, label, value string, width int) string {
	text := fmt.Sprintf("%s: %s", label, value)
	fmt.Printf("Rendering badge with width %d for text: %s\n", width, text)
	return fmt.Sprintf(
		`<g transform="translate(%d,%d)">`+
			`<rect width="%d" height="32" rx="8" ry="8" fill="rgba(0,0,0,0.6)" />`+
			`<text x="12" y="20" fill="white" font-family="Courier, Liberation Mono, monospace" font-size="13">%s</text>`+
			`</g>`,
		x, y, width, text,
	)
}

func (a *AgentSVGImage) renderMultiLineBadge(x, y int, label, value, subtext string, width int) string {
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
