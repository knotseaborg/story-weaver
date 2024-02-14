package search

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
	"github.com/knotseaborg/wikiSearchServer/gpt"
)

func ExecuteQuery(plan gpt.QueryPlan) (*queryResult, error) {
	/*
		ExecuteQuery executes a query on Wikidata based on the provided query plan.

		Parameters:
		plan: The query plan containing the query string and list of searchable entities.

		Returns:
		*queryResult: A pointer to the queryResult struct containing the result of the executed query on Wikidata.
		error:        An error if there is any issue during the query execution or building process.

		Example:
		plan := gpt.QueryPlan{
			Query: "SELECT ?company WHERE { ?company rdf:type <Company> . ?company <Location> <Location_ID> . }",
			Searchable: []gpt.Searchable{
				{Name: "Location", Keyword: "New York", Description: "A city in the United States"},
				{Name: "Company", Keyword: "Apple Inc.", Description: "A technology company"},
			},
		}
		queryResult, err := ExecuteQuery(plan)
		if err != nil {
			log.Fatal("Error querying Wikidata:", err)
		}
		fmt.Println("Query result:", queryResult)
	*/
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
	/*
		queryWiki queries Wikidata using the provided SPARQL query.

		Parameters:
		query: The SPARQL query to be executed on Wikidata.

		Returns:
		*queryResult: A pointer to the queryResult struct containing the result of the executed SPARQL query on Wikidata.
		error:        An error if there is any issue during the query execution or JSON unmarshalling.

		Example:
		sparqlQuery := "SELECT ?country ?countryLabel WHERE { ?country wdt:P31 wd:Q6256. SERVICE wikibase:label { bd:serviceParam wikibase:language \"[AUTO_LANGUAGE],en\". }}"
		queryResult, err := queryWiki(sparqlQuery)
		if err != nil {
			log.Fatal("Error querying Wikidata:", err)
		}
		fmt.Println("Query result:", queryResult)
	*/
	rawParams := map[string]string{"query": query}
	resp, err := common.RequestGET(os.Getenv("WIKIDATA_QUERY_URL"), rawParams)
	if err != nil {
		return nil, err
	}
	result := queryResult{}
	json.Unmarshal(resp, &result)
	return &result, nil
}

func buildQuery(plan gpt.QueryPlan) (string, error) {
	/*
		buildQuery builds an executable query utilizing the provided query plan.

		Parameters:
		  plan: The query plan containing the query string and list of searchable entities.

		Returns:
		  string: The built query with placeholders replaced by item IDs.
		  error:  An error if there is any issue during the building process.

		Example:
		  plan := gpt.QueryPlan{
		      Query: "SELECT ?company WHERE { ?company rdf:type <Company> . ?company <Location> <Location_ID> . }",
		      Searchable: []gpt.Searchable{
		          {Name: "Location", Keyword: "New York", Description: "A city in the United States"},
		          {Name: "Company", Keyword: "Apple Inc.", Description: "A technology company"},
		      },
		  }
		  builtQuery, err := buildQuery(plan)
		  if err != nil {
		      log.Fatal("Error building query:", err)
		  }
		  fmt.Println("Built query:", builtQuery)
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
