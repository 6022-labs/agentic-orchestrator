package interfaces

type IPFSFileStorageInterface interface {
	UploadFile(fileBytes []byte, filename string) (*string, error)
}
