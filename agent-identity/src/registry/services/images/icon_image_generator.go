package images

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/settings"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
	_ "golang.org/x/image/webp"
)

const (
	iconWidth       = 512
	iconHeight      = 768
	iconBorder      = 20
	iconInnerWidth  = iconWidth - 2*iconBorder
	iconInnerHeight = iconHeight - 2*iconBorder
	iconOuterRadius = 20
	iconInnerRadius = 16
)

// Bold so the icon's text matches the heavier weight the browser gives the SVG card.
var goMonoFont, _ = opentype.Parse(gomonobold.TTF)

type IconImageGenerator struct {
	*AbstractImageGenerator
	imageGenerationSettings *settings.ImageGenerationSettings
}

func NewIconImageGenerator(imageGenerationSettings *settings.ImageGenerationSettings) *IconImageGenerator {
	return &IconImageGenerator{
		AbstractImageGenerator:  NewAbstractImageGenerator(),
		imageGenerationSettings: imageGenerationSettings,
	}
}

func (g *IconImageGenerator) Generate(params interfaces.GenerateImageParams) ([]byte, error) {
	// font.Face holds mutable buffers, so build per call rather than sharing.
	canvas := &iconCanvas{
		AbstractImageGenerator: g.AbstractImageGenerator,
		img:                    image.NewRGBA(image.Rect(0, 0, iconWidth, iconHeight)),
		face12:                 g.mustFace(12),
		face13:                 g.mustFace(13),
		face14:                 g.mustFace(14),
	}

	canvas.drawBorderGradient(params.Creator)

	if len(params.BaseImageBytes) > 0 {
		if err := canvas.drawBaseImage(params.BaseImageBytes); err != nil {
			return nil, err
		}
	}

	canvas.drawFadeOverlay()
	canvas.drawAgentInformation(params)
	canvas.drawWallet(params.Wallets)

	if g.imageGenerationSettings.ImageMarkPath != nil {
		if err := canvas.drawMark(*g.imageGenerationSettings.ImageMarkPath); err != nil {
			return nil, err
		}
	}

	return canvas.encodePNG()
}

func (g *IconImageGenerator) mustFace(points float64) font.Face {
	// DPI 72 makes 1pt == 1px, matching the SVG font-size units.
	face, err := opentype.NewFace(goMonoFont, &opentype.FaceOptions{Size: points, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic(fmt.Sprintf("failed to build gomono face: %v", err))
	}
	return face
}

type iconCanvas struct {
	*AbstractImageGenerator
	img                    *image.RGBA
	face12, face13, face14 font.Face
}

func (c *iconCanvas) drawBorderGradient(creator string) {
	colors := c.borderGradientColors(creator)
	stops := make([]gradientStop, len(colors))
	for i, hex := range colors {
		stops[i] = gradientStop{pos: float64(i) / float64(len(colors)-1), col: c.parseHexColor(hex)}
	}

	grad := c.linearGradient(c.img.Bounds(), 0, 0, iconWidth, 0, stops)
	mask := c.roundedRectMask(0, 0, iconWidth, iconHeight, iconOuterRadius)
	xdraw.DrawMask(c.img, c.img.Bounds(), grad, image.Point{}, mask, image.Point{}, xdraw.Over)
}

func (c *iconCanvas) drawBaseImage(data []byte) error {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode base image: %w", err)
	}

	covered := c.coverResize(img, iconInnerWidth, iconInnerHeight)
	layer := image.NewRGBA(c.img.Bounds())
	xdraw.Draw(layer, image.Rect(iconBorder, iconBorder, iconBorder+iconInnerWidth, iconBorder+iconInnerHeight), covered, image.Point{}, xdraw.Src)

	mask := c.roundedRectMask(iconBorder, iconBorder, iconInnerWidth, iconInnerHeight, iconInnerRadius)
	xdraw.DrawMask(c.img, c.img.Bounds(), layer, image.Point{}, mask, image.Point{}, xdraw.Over)
	return nil
}

func (c *iconCanvas) drawFadeOverlay() {
	grad := c.linearGradient(c.img.Bounds(),
		iconBorder, iconBorder, iconBorder+iconInnerWidth, iconBorder+iconInnerHeight,
		[]gradientStop{
			{0, color.RGBA{A: 255}},
			{0.10, color.RGBA{}},
			{0.90, color.RGBA{}},
			{1, color.RGBA{A: 255}},
		},
	)
	mask := c.roundedRectMask(iconBorder, iconBorder, iconInnerWidth, iconInnerHeight, iconInnerRadius)
	xdraw.DrawMask(c.img, c.img.Bounds(), grad, image.Point{}, mask, image.Point{}, xdraw.Over)
}

func (c *iconCanvas) drawAgentInformation(p interfaces.GenerateImageParams) {
	if p.OriginEnsDomain != nil {
		origin := fmt.Sprintf("%s (%d)", *p.OriginEnsDomain, *p.OriginGeneration)
		c.drawMultiLineBadge(50, 522, fmt.Sprintf("Origin: %s", origin), *p.OriginEnsDomain)
	}

	created := fmt.Sprint(p.CreatedAt.Unix() + 6102950400)
	c.drawBadge(48, 580, fmt.Sprintf("Role: %s", p.Role))
	c.drawBadge(48, 620, fmt.Sprintf("Created: %s", created))
	c.drawBadge(48, 660, fmt.Sprintf("Creator: %s", p.Creator))
}

func (c *iconCanvas) drawBadge(x, y float64, text string) {
	width := c.measure(c.face13, text) + 24
	c.fillRoundedRect(x, y, width, 32, 8, color.RGBA{A: 153}) // black at 0.6 opacity
	c.drawText(c.face13, color.White, text, x+12, y+20)
}

func (c *iconCanvas) drawMultiLineBadge(x, y float64, main, subtext string) {
	width := math.Max(c.measure(c.face13, main), c.measure(c.face12, subtext)) + 24
	c.fillRoundedRect(x, y, width, 48, 8, color.RGBA{A: 153})
	c.drawText(c.face13, color.White, main, x+12, y+18)
	c.drawText(c.face12, color.White, subtext, x+12, y+35)
}

func (c *iconCanvas) drawWallet(wallets []string) {
	if len(wallets) == 0 {
		return
	}
	width := c.measure(c.face14, wallets[0])
	c.drawText(c.face14, color.White, wallets[0], 256-width/2, 762)
}

func (c *iconCanvas) drawMark(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read image mark file: %w", err)
	}
	mark, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image mark: %w", err)
	}

	const (
		cx, cy  = 460.0, 56.0
		imgSize = 40
	)

	c.fillCircle(cx, cy, 23, color.Black)
	c.fillCircle(cx, cy, 21, color.White)

	covered := c.coverResize(mark, imgSize, imgSize)
	layer := image.NewRGBA(c.img.Bounds())
	xdraw.Draw(layer, image.Rect(int(cx)-imgSize/2, int(cy)-imgSize/2, int(cx)+imgSize/2, int(cy)+imgSize/2), covered, image.Point{}, xdraw.Src)
	mask := c.circleMask(cx, cy, imgSize/2)
	xdraw.DrawMask(c.img, c.img.Bounds(), layer, image.Point{}, mask, image.Point{}, xdraw.Over)
	return nil
}

func (c *iconCanvas) fillRoundedRect(x, y, w, h, r float64, col color.Color) {
	mask := c.roundedRectMask(x, y, w, h, r)
	xdraw.DrawMask(c.img, c.img.Bounds(), image.NewUniform(col), image.Point{}, mask, image.Point{}, xdraw.Over)
}

func (c *iconCanvas) fillCircle(cx, cy, r float64, col color.Color) {
	mask := c.circleMask(cx, cy, r)
	xdraw.DrawMask(c.img, c.img.Bounds(), image.NewUniform(col), image.Point{}, mask, image.Point{}, xdraw.Over)
}

func (c *iconCanvas) encodePNG() ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, c.img); err != nil {
		return nil, fmt.Errorf("failed to encode icon png: %w", err)
	}
	return buf.Bytes(), nil
}

func (c *iconCanvas) drawText(face font.Face, col color.Color, text string, x, y float64) {
	drawer := &font.Drawer{
		Dst:  c.img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(math.Round(x * 64)), Y: fixed.Int26_6(math.Round(y * 64))},
	}
	drawer.DrawString(text)
}

func (c *iconCanvas) measure(face font.Face, text string) float64 {
	return float64(font.MeasureString(face, text)) / 64
}

type gradientStop struct {
	pos float64
	col color.RGBA
}

// stops must be ordered by pos.
func (c *iconCanvas) linearGradient(b image.Rectangle, x0, y0, x1, y1 float64, stops []gradientStop) *image.RGBA {
	img := image.NewRGBA(b)
	dx, dy := x1-x0, y1-y0
	length := dx*dx + dy*dy
	if length == 0 {
		length = 1
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			t := ((float64(x)-x0)*dx + (float64(y)-y0)*dy) / length
			img.SetRGBA(x, y, c.gradientColorAt(stops, t))
		}
	}
	return img
}

func (c *iconCanvas) gradientColorAt(stops []gradientStop, t float64) color.RGBA {
	if t <= stops[0].pos {
		return stops[0].col
	}
	last := len(stops) - 1
	if t >= stops[last].pos {
		return stops[last].col
	}
	for i := 1; i <= last; i++ {
		if t <= stops[i].pos {
			span := stops[i].pos - stops[i-1].pos
			local := 0.0
			if span > 0 {
				local = (t - stops[i-1].pos) / span
			}
			return c.lerpColor(stops[i-1].col, stops[i].col, local)
		}
	}
	return stops[last].col
}

func (c *iconCanvas) lerpColor(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: c.lerpU8(a.R, b.R, t),
		G: c.lerpU8(a.G, b.G, t),
		B: c.lerpU8(a.B, b.B, t),
		A: c.lerpU8(a.A, b.A, t),
	}
}

func (c *iconCanvas) lerpU8(a, b uint8, t float64) uint8 {
	return uint8(math.Round(float64(a) + (float64(b)-float64(a))*t))
}

func (c *iconCanvas) roundedRectMask(x, y, w, h, r float64) *image.Alpha {
	ras := vector.NewRasterizer(iconWidth, iconHeight)
	fx, fy, fw, fh, fr := float32(x), float32(y), float32(w), float32(h), float32(r)
	ras.MoveTo(fx+fr, fy)
	ras.LineTo(fx+fw-fr, fy)
	ras.QuadTo(fx+fw, fy, fx+fw, fy+fr)
	ras.LineTo(fx+fw, fy+fh-fr)
	ras.QuadTo(fx+fw, fy+fh, fx+fw-fr, fy+fh)
	ras.LineTo(fx+fr, fy+fh)
	ras.QuadTo(fx, fy+fh, fx, fy+fh-fr)
	ras.LineTo(fx, fy+fr)
	ras.QuadTo(fx, fy, fx+fr, fy)
	ras.ClosePath()
	return c.rasterizeMask(ras)
}

func (c *iconCanvas) circleMask(cx, cy, r float64) *image.Alpha {
	const kappa = 0.5522847498307936
	ras := vector.NewRasterizer(iconWidth, iconHeight)
	fcx, fcy, fr, fk := float32(cx), float32(cy), float32(r), float32(r*kappa)
	ras.MoveTo(fcx+fr, fcy)
	ras.CubeTo(fcx+fr, fcy+fk, fcx+fk, fcy+fr, fcx, fcy+fr)
	ras.CubeTo(fcx-fk, fcy+fr, fcx-fr, fcy+fk, fcx-fr, fcy)
	ras.CubeTo(fcx-fr, fcy-fk, fcx-fk, fcy-fr, fcx, fcy-fr)
	ras.CubeTo(fcx+fk, fcy-fr, fcx+fr, fcy-fk, fcx+fr, fcy)
	ras.ClosePath()
	return c.rasterizeMask(ras)
}

func (c *iconCanvas) rasterizeMask(ras *vector.Rasterizer) *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, iconWidth, iconHeight))
	ras.Draw(mask, mask.Bounds(), image.Opaque, image.Point{})
	return mask
}

func (c *iconCanvas) coverResize(img image.Image, w, h int) image.Image {
	bounds := img.Bounds()
	iw, ih := bounds.Dx(), bounds.Dy()
	if iw == 0 || ih == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	scale := math.Max(float64(w)/float64(iw), float64(h)/float64(ih))
	sw, sh := int(math.Ceil(float64(iw)*scale)), int(math.Ceil(float64(ih)*scale))

	scaled := image.NewRGBA(image.Rect(0, 0, sw, sh))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), img, bounds, xdraw.Over, nil)

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.Draw(out, out.Bounds(), scaled, image.Pt((sw-w)/2, (sh-h)/2), xdraw.Src)
	return out
}

func (c *iconCanvas) parseHexColor(s string) color.RGBA {
	v, err := strconv.ParseUint(strings.TrimPrefix(s, "#"), 16, 32)
	if err != nil {
		return color.RGBA{A: 255}
	}
	return color.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 255}
}
