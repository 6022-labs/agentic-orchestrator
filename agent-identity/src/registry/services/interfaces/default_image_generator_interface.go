package interfaces

// DefaultImageGeneratorInterface renders the animated "default" agent frame (SVG with marquees).
type DefaultImageGeneratorInterface interface {
	Generate(params GenerateImageParams) ([]byte, error)
}
