package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/curioswitch/go-build"
	"github.com/google/go-github/v62/github"
	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/boot"
	"gopkg.in/yaml.v2"

	"github.com/curioswitch/tasuke/common/languages"
)

func main() {
	build.DefineTasks(build.ExcludeTasks("lint-go", "format-go"))

	goyek.Define(goyek.Task{
		Name: "update-languages",
		Action: func(a *goyek.A) {
			updateLanguages(a)
		},
	})

	boot.Main()
}

func updateLanguages(a *goyek.A) {
	tok, _ := auth.TokenForHost("github.com")
	if tok == "" {
		a.Fatal("Failed to GitHub token, make sure to install gh CLI.")
		return
	}

	gh := github.NewClient(http.DefaultClient).WithAuthToken(tok)

	languagesYML := fetchFile(a, gh, "github-linguist", "linguist", "lib/linguist/languages.yml")
	var allLangs map[string]map[string]any
	if err := yaml.Unmarshal(languagesYML, &allLangs); err != nil {
		a.Fatal(err)
	}

	popularYML := fetchFile(a, gh, "github", "linguist", "lib/linguist/popular.yml")
	var popular []string
	if err := yaml.Unmarshal(popularYML, &popular); err != nil {
		a.Fatal(err)
	}

	var res []languages.Language
	for lang, langData := range allLangs {
		if slices.Contains(popular, lang) {
			res = append(res, languages.Language{
				ID:   langData["language_id"].(int),
				Name: lang,
			})
		}
	}

	slices.SortFunc(res, func(a, b languages.Language) int {
		return strings.Compare(a.Name, b.Name)
	})

	resJSON, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		a.Fatal(err)
	}

	clientPath := filepath.Join("frontend", "client", "src", "pages", "settings", "languages.json")
	if err := os.WriteFile(clientPath, resJSON, 0o644); err != nil {
		a.Fatal(err)
	}

	goPath := filepath.Join("common", "go", "languages", "languages.json")
	if err := os.WriteFile(goPath, resJSON, 0o644); err != nil {
		a.Fatal(err)
	}
}

func fetchFile(a *goyek.A, gh *github.Client, owner string, repo string, path string) []byte {
	r, _, err := gh.Repositories.DownloadContents(a.Context(), owner, repo, path, nil)
	if err != nil {
		a.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		a.Fatal(err)
	}
	return b
}
