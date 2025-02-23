package struct

type AssetMessage struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	FileData  []byte `json:"fileData,omitempty"`
}