package interfaces

// IconImageGeneratorInterface renders the static "icon" agent frame as a PNG (rasterized from SVG, no marquees/stars).
type IconImageGeneratorInterface interface {
	Generate(params GenerateImageParams) ([]byte, error)
}
