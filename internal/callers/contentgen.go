package callers

import (
	"encoding/json"
	"fmt"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader"
	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/configloader/structs"
)

func LoadIntentDetectionResponse(JSONData []byte) {
	//read in JSON data from intent detection
	var intentDetectionResponse structs.IntentDetectionResponse
	if err := json.Unmarshal(JSONData, &intentDetectionResponse); err != nil {
		fmt.Println("Error unmarshalling intent detection response")
		fmt.Println(err)
		return
	}

	//load content gen yaml based off JSON data
	var apiConfig structs.APIConfig
	var err error
	switch intentDetectionResponse.ContentType {
	case "image":
		apiConfig, err = configloader.GetImage()
	case "video":
		apiConfig, err = configloader.GetVideo()
	case "gif":
		apiConfig, err = configloader.GetGif()
	case "3d":
		apiConfig, err = configloader.Get3d()
	default:
		fmt.Println("Intent detection provided invalid content type")
		return
	}

	if err != nil {
		fmt.Println("Error loading content gen settings")
		fmt.Println(err)
		return
	}

	BuildAPICall(intentDetectionResponse, apiConfig)
}

func BuildAPICall(intentDetectionResponse structs.IntentDetectionResponse, apiConfig structs.APIConfig) {
	// Build API call using intentDetectionResponse and yamlConfig
	//fmt.Println(intentDetectionResponse)
	//fmt.Println(apiConfig)

}
