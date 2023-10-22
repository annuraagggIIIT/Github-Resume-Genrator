package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Language    string `json:"language"`
}

func getRepositories(username string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching the data")
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading the data")
		return nil, err
	}

	var repos []Repository
	err = json.Unmarshal(body, &repos)
	if err != nil {
		fmt.Println("Error while unmarshalling the data")
		return nil, err
	}

	return repos, nil
}

func generateResume(username string, repos []Repository) string {
	resume := fmt.Sprintf("GitHub Resume\nUsername: %s\n\n", username)

	// Get most used language among all the repositories
	langCount := make(map[string]int)
	for _, repo := range repos {
		langCount[repo.Language]++
	}

	// Sort languages by count
	var languages []string
	for lang := range langCount {
		languages = append(languages, lang)
	}
	sort.Slice(languages, func(i, j int) bool {
		return langCount[languages[i]] > langCount[languages[j]]
	})

	// top theree projects 
	resume += "Skills:\n"
	for i := 0; i < 3 && i < len(languages); i++ {
		resume += fmt.Sprintf("- %s\n", languages[i])
	}

	resume += "\nTop 3 Projects:\n"
	for i := 0; i < 3 && i < len(repos); i++ {
		resume += fmt.Sprintf("- %s (%s)\n  Description: %s\n  Stars: %d, Forks: %d\n\n", repos[i].Name, repos[i].Language, repos[i].Description, repos[i].Stars, repos[i].Forks)
	}

	return resume
}

func main() {
	fmt.Println("Enter the GitHub Username:")
	var username string
	fmt.Scanln(&username)

	repos, err := getRepositories(username)
	if err != nil {
		fmt.Println("Error while fetching the data:", err)
		return
	}

	resume := generateResume(username, repos)
	fmt.Println(resume)
}
