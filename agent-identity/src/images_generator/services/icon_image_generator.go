package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator/settings"
	"github.com/chromedp/chromedp"
)

type IconImageGeneratorInterface interface {
	Generate(params GenerateImageParams) ([]byte, error)
}

type IconImageGenerator struct {
	agentSVGImage                AgentSVGImageInterface
	agentImageGenerationSettings *settings.ImageGenerationSettings
}

func NewIconImageGenerator(
	agentSVGImage AgentSVGImageInterface,
	agentImageGenerationSettings *settings.ImageGenerationSettings,
) *IconImageGenerator {
	return &IconImageGenerator{
		agentSVGImage:                agentSVGImage,
		agentImageGenerationSettings: agentImageGenerationSettings,
	}
}

func (g *IconImageGenerator) Generate(params GenerateImageParams) ([]byte, error) {
	imageBytes, err := g.renderBaseSvg(params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate static image: %w", err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// limit execution time to avoid hanging chrome processes
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var buf []byte
	// Write SVG to a temp file
	tmpFile, err := os.CreateTemp("", "icon_image_*.svg")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write(imageBytes)
	if err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	fileURL := "file://" + tmpFile.Name()

	// Navigate and wait for SVG to be ready
	if err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("svg", chromedp.ByQuery),
		chromedp.EmulateViewport(512, 768),
		chromedp.Evaluate(`document.querySelector('svg').style.background = "none"`, nil),
	); err != nil {
		return nil, fmt.Errorf("failed to load SVG: %w", err)
	}

	// Take screenshot
	if err := chromedp.Run(ctx,
		chromedp.Screenshot("svg", &buf, chromedp.NodeVisible, chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return buf, nil
}

func (g *IconImageGenerator) renderBaseSvg(params GenerateImageParams) ([]byte, error) {
	var sb strings.Builder

	borderGradient := g.agentSVGImage.RenderBorderGradient(params.Creator)
	imageDataURI := g.agentSVGImage.RenderImageData(params.BaseImageBytes, params.BaseImageMimeType)

	sb.WriteString(`<svg width="512" height="768" viewBox="0 0 512 768" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`)
	sb.WriteString(`<defs>`)
	sb.WriteString(`<clipPath id="clip-border"><rect width="512" height="768" rx="20" ry="20" /></clipPath>`)
	sb.WriteString(`<clipPath id="clip-image"><rect x="20" y="20" width="472" height="728" rx="16" ry="16" /></clipPath>`)
	sb.WriteString(`<linearGradient id="fadeOverlay" x1="0" y1="0" x2="1" y2="1">`)
	sb.WriteString(`<stop offset="0%" stop-color="black" stop-opacity="1"/>`)
	sb.WriteString(`<stop offset="10%" stop-color="black" stop-opacity="0"/>`)
	sb.WriteString(`<stop offset="90%" stop-color="black" stop-opacity="0"/>`)
	sb.WriteString(`<stop offset="100%" stop-color="black" stop-opacity="1"/>`)
	sb.WriteString(`</linearGradient>`)
	sb.WriteString(borderGradient)
	sb.WriteString(`<path id="carouselPath" d="M 0,762 H 512" />`)
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
	sb.WriteString(g.agentSVGImage.RenderAgentInformation(params))
	sb.WriteString(g.agentSVGImage.RenderStars(params.Stars))
	sb.WriteString(g.renderWalletsStatic(params.Wallets))

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

func (a *IconImageGenerator) renderWalletsStatic(wallets []string) string {
	if len(wallets) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(`<g font-family="Courier, Liberation Mono, monospace" font-size="14" fill="white">`)
	y := 762
	x := 256
	walletStr := strings.Join(wallets, "  ")
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle">%s</text>`, x, y, walletStr))
	sb.WriteString(`</g>`)
	return sb.String()
}
