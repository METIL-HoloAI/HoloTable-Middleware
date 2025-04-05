package assetstruct

type AssetMessageFile struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	FilePath  string `json:"filePath,omitempty"`
}

type AssetMessageData struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	FileData  []byte `json:"fileData,omitempty"`
}
