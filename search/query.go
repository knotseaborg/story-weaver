package search

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
	"github.com/knotseaborg/wikiSearchServer/gpt"
)

func Query(plan gpt.QueryPlan) (*queryResult, error) {
	query, err := buildQuery(plan)
	if err != nil {
		log.Println("Error building query:", err)
	}
	log.Println("Query built:", query)
	result, err := queryWiki(query)
	if err != nil {
		log.Println("Error querying wiki:", err)
	}
	return result, nil
}

func queryWiki(query string) (*queryResult, error) {
	baseURL := "https://query.wikidata.org/sparql"
	rawParams := map[string]string{"query": query}
	resp, err := common.RequestGET(baseURL, rawParams)
	if err != nil {
		return nil, err
	}
	result := queryResult{}
	json.Unmarshal(resp, &result)
	return &result, nil
}

func buildQuery(plan gpt.QueryPlan) (string, error) {
	/*
		Builds a query utilizing the query plan
	*/
	query := plan.Query
	for _, searchable := range plan.Searchable {
		candidates := ClusterSearch(searchable.Keyword, 5)
		itemID, err := SearchID(candidates, searchable)
		if err != nil {
			return "", nil
		}
		placeholder := fmt.Sprintf("<%s>", searchable.Name)
		query = strings.ReplaceAll(query, placeholder, itemID)
	}
	return query, nil
}
