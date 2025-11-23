package models

type ExternalAddResponse struct {
	Bytes      int64  `json:"Bytes"`
	Hash       string `json:"Hash"`
	Mode       string `json:"Mode"`
	Mtime      int64  `json:"Mtime"`
	MtimeNsecs int    `json:"MtimeNsecs"`
	Name       string `json:"Name"`
	Size       string `json:"Size"`
}
