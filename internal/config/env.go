package config

import (
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() {
	envLoc, err := getConfigPath("../.env")
	if err != nil {
		log.Fatal("Error getting .env file path: ", err)
	}
	// Load API keys
	err = godotenv.Load(envLoc)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	replaceEnv(reflect.ValueOf(&General).Elem())
	replaceEnv(reflect.ValueOf(&IntentDetection).Elem())
	replaceEnv(reflect.ValueOf(&ImageGen).Elem())
	replaceEnv(reflect.ValueOf(&VideoGen).Elem())
	replaceEnv(reflect.ValueOf(&GifGen).Elem())
	replaceEnv(reflect.ValueOf(&ModelGen).Elem())
}

func replaceEnv(reflectedItem reflect.Value) {
	regex, err := regexp.Compile(`\$[A-Z_]+`)
	if err != nil {
		log.Fatal("Failed to compile regex, ", err)
	}

	switch reflectedItem.Kind() {
	case reflect.Ptr:
		if !reflectedItem.IsNil() {
			replaceEnv(reflectedItem.Elem()) // Dereference and process
		}
	case reflect.Struct:
		for i := 0; i < reflectedItem.NumField(); i++ {
			replaceEnv(reflectedItem.Field(i))
		}
	case reflect.Map:
		for _, key := range reflectedItem.MapKeys() {
			val := reflectedItem.MapIndex(key)
			if val.Kind() == reflect.String {
				matches := regex.FindAllString(val.String(), -1)
				for _, match := range matches {
					newVal := os.Getenv(strings.Split(match, "$")[1])
					reflectedItem.SetMapIndex(key, reflect.ValueOf(strings.ReplaceAll(val.String(), match, newVal)))
				}
			} else {
				replaceEnv(val) // Process nested values
			}
		}
	case reflect.Slice:
		for i := 0; i < reflectedItem.Len(); i++ {
			replaceEnv(reflectedItem.Index(i))
		}
	case reflect.String:
		matches := regex.FindAllString(reflectedItem.String(), -1)
		for _, match := range matches {
			newVal := os.Getenv(strings.Split(match, "$")[1])
			reflectedItem.SetString(strings.ReplaceAll(reflectedItem.String(), match, newVal))
		}
	}
}
