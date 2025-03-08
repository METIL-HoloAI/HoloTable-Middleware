package config

import (
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func loadEnv() {
	envLoc, err := getConfigPath("../.env")
	if err != nil {
		logrus.Fatal("Error getting .env file path: ", err)
	}
	// Load API keys
	err = godotenv.Load(envLoc)
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	replaceEnv(reflect.ValueOf(&General).Elem())
	replaceEnv(reflect.ValueOf(&IntentDetection).Elem())

	replaceEnv(reflect.ValueOf(&ImageGen).Elem())
	replaceEnv(reflect.ValueOf(&VideoGen).Elem())
	replaceEnv(reflect.ValueOf(&GifGen).Elem())
	replaceEnv(reflect.ValueOf(&ModelGen).Elem())

	// MANUALLY TRAVERSE EACH WORKFLOW & STEP TO REPLACE ENV VARIABLES
	for workflowKey, workflow := range Workflows {
		for i := range workflow.Steps { // Iterate by reference to persist changes
			logrus.Tracef("🔍 Searching for ENV keys within Workflow: '%s', Step: %d\n", workflowKey, i)

			// Replace env variables inside Headers map
			for headerKey, headerValue := range workflow.Steps[i].Headers {
				workflow.Steps[i].Headers[headerKey] = os.ExpandEnv(headerValue) // ✅ Direct replacement
			}

			// Apply `replaceEnv()` for deeper replacements
			replaceEnv(reflect.ValueOf(&workflow.Steps[i]).Elem())

			// apply changes back into `Workflows`
			Workflows[workflowKey] = workflow
		}
	}
}

func replaceEnv(reflectedItem reflect.Value) {
	regex, err := regexp.Compile(`\&[A-Z_]+`)
	if err != nil {
		logrus.Fatal("Failed to compile regex, ", err)
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
					newVal := os.Getenv(strings.Split(match, "&")[1])
					if newVal != "" {
						reflectedItem.SetMapIndex(key, reflect.ValueOf(strings.ReplaceAll(val.String(), match, newVal)))
					}
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
			newVal := os.Getenv(strings.Split(match, "&")[1])
			if newVal != "" {
				reflectedItem.SetString(strings.ReplaceAll(reflectedItem.String(), match, newVal))
			}
		}
	}
}
