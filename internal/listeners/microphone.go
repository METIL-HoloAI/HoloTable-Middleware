package listeners

import (
	"strings"

	"github.com/METIL-HoloAI/HoloTable-Middleware/internal/config"
)

func CheckForKeyword(message string) bool {
	return strings.Contains(message, config.SpeechToText.Keyword)
}
