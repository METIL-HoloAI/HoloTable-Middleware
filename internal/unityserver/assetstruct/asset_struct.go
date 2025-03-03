package assetstruct

type AssetMessage struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	FilePath  string `json:"filePath,omitempty"`
}
