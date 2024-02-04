package search

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
	"github.com/knotseaborg/wikiSearchServer/gpt"
)

func SearchID(candidates []*search, criteria gpt.Searchable) (string, error) {
	/*Find and return the most suitable item ID from the candidates*/
	description := criteria.Description
	reference := strings.Builder{}
	history := map[string]struct{}{} // Used to check for duplicate ID
	for _, candidate := range candidates {
		for _, item := range candidate.Search {
			if _, ok := history[item.ID]; ok {
				continue
			}
			_, err := reference.WriteString(fmt.Sprintf("ID: %s, Label: %s, Description: %s\n", item.ID, item.Label, item.Description))
			if err != nil {
				return "", err
			}
			history[item.ID] = struct{}{}
		}
	}
	prompt := fmt.Sprintf("Given the description:\n%s\n\nWhich ID is the most suitable from the given list?\n%s\nThe most suitable is ID is: ", description, reference.String())
	prompt = common.CleanForJSON(prompt)
	resp, err := gpt.Completion(prompt, gpt.GPT_35)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`Q[0-9]*`)
	return string(re.Find([]byte(resp))), nil
}

func ClusterSearch(keyword string, limit int) []*search {
	result := []*search{}
	for i := range keyword {
		key := keyword[:i+1]
		buf := SearchWiki(key, limit)
		if len(buf.Search) == 0 {
			continue
		}
		result = append(result, buf)
	}
	return result
}

func SearchWiki(keyword string, limit int) *search {
	/*
		Searches wikidata for a given keyword
	*/
	baseURL := "https://www.wikidata.org/w/api.php"
	rawParams := map[string]string{
		"action":   "wbsearchentities",
		"search":   keyword,
		"format":   "json",
		"language": "en",
		"uselang":  "en",
		"type":     "item",
		"limit":    fmt.Sprintf("%d", limit)}
	var search search
	body, err := common.RequestGET(baseURL, rawParams)
	if err != nil {
		log.Panic()
	}
	json.Unmarshal(body, &search)
	return &search
}
