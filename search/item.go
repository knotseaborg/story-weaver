package search

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
	"github.com/knotseaborg/wikiSearchServer/gpt"
)

func SearchID(candidates []*search, criteria gpt.Searchable) (string, error) {
	/*
		SearchID finds and returns the most suitable item ID from a list of candidate search results based on specified criteria.

		Parameters:
		candidates: A slice of search structs containing candidate search results from Wikidata.
		criteria:   The criteria (Searchable) based on which the most suitable item ID is determined.

		Returns:
		string: The most suitable item ID chosen by the user.
		error:  An error if there is any issue during the completion process or ID extraction.

		Example:
		candidates := []*search{search1, search2, ...}
		criteria := gpt.Searchable{Name: "Apple Inc.", Keyword: "Apple", Description: "Technology company"}
		itemID, err := SearchID(candidates, criteria)
		if err != nil {
			log.Fatal("Error searching ID:", err)
		}
		fmt.Println("Most suitable item ID:", itemID)
	*/
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
	resp, err := gpt.Completion(prompt, os.Getenv("GPT_MODEL_BASIC"))
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`Q[0-9]*`)
	return string(re.Find([]byte(resp))), nil
}

func ClusterSearch(keyword string, limit int) []*search {
	/*
		ClusterSearch performs clustered searching on Wikidata using a keyword.

		Parameters:
		keyword: The keyword to perform clustered searching on Wikidata.
		limit:   The maximum number of search results to retrieve for each substring.

		Returns:
		[]*search: A slice of search structs containing search results from Wikidata for each clustered search.

		Example:
		clusteredResults := ClusterSearch("Apple Inc.", 5)
		for _, result := range clusteredResults {
			fmt.Println("Clustered search results:", result)
		}
	*/
	result := []*search{}
	// Search all possible prefixes of the keyword
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
		SearchWiki searches Wikidata for wikidata entities related to the given keyword.

		Parameters:
		  keyword: The keyword to search for in Wikidata.
		  limit:   The maximum number of search results to retrieve.

		Returns:
		  *search: A pointer to the search struct containing search results from Wikidata.

		Example:
		  searchResults := SearchWiki("Apple Inc.", 5)
		  fmt.Println("Search results:", searchResults)
	*/
	rawParams := map[string]string{
		"action":   "wbsearchentities",
		"search":   keyword,
		"format":   "json",
		"language": "en",
		"uselang":  "en",
		"type":     "item",
		"limit":    fmt.Sprintf("%d", limit)}
	var search search
	body, err := common.RequestGET(os.Getenv("WIKIDATA_SEARCH_URL"), rawParams)
	if err != nil {
		log.Panic()
	}
	json.Unmarshal(body, &search)
	return &search
}
