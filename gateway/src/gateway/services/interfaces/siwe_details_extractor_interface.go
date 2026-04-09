package interfaces

type SIWEDetailsExtractorInterface interface {
	ExtractSIWEMessageDetails(message string) (*string, *uint64, *string, error)
}
