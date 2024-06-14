package languages

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// Language is information about a programming language used during
// code review matching.
type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IsSupported returns whether the given language ID is supported by
// tasuke.
func IsSupported(id int) bool {
	_, ok := supportedLanguages[id]
	return ok
}

//go:embed languages.json
var languagesJSON []byte

var supportedLanguages = createSupportedLanguages()

func createSupportedLanguages() map[int]Language {
	var langs []Language
	if err := json.Unmarshal(languagesJSON, &langs); err != nil {
		panic(fmt.Sprintf("languages: failed to unmarshal languages: %v", err))
	}

	res := make(map[int]Language, len(langs))
	for _, lang := range langs {
		res[lang.ID] = lang
	}
	return res
}
