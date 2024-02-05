package weaver

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
	"github.com/knotseaborg/wikiSearchServer/gpt"
	"github.com/knotseaborg/wikiSearchServer/search"
)

func Run() {
	defer func() { fmt.Println("Thank you playing!") }()
	history := strings.Builder{}
	context := "Narrator: One of the lead actresses of Avatar 2 appeared to be preparing for another movie."
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println("\nWhat would you do?")
		text := strings.Trim(common.ReadInput(reader), "\n")
		if text == "stop" || text == "bye" || text == "exit" {
			return
		}
		queryIntents, err := gpt.GenerateQueryIntents(context + text)
		if err != nil {
			log.Println("Error generating query intents", err)
		}
		reference := generateReference(queryIntents, history.String())
		history.WriteString(reference)

		context = fmt.Sprintf("%s\nUser:%s", context, text)
		nextText := generateResponse(context, reference)
		context = fmt.Sprintf("%s\n%s", context, nextText)

		imgURL, err := generateImage(context)
		if err == nil {
			fmt.Printf("\n---------------CLICK TO SEE THE IMAGE-----------------\n%s\nf", imgURL)
		}
		fmt.Printf("\n---------------WHAT HAPPENS NEXT?-----------------\n%s\n", nextText)
		//fmt.Println(history.String())
	}
}

func generateImage(context string) (string, error) {
	limit, err := strconv.Atoi(os.Getenv("DALL_E_PROMPT_LIMIT"))
	if err != nil {
		log.Panic("Error: Unable to read DALL_E_PROMPT_LIMIT from env", err)
	}
	context = common.PruneText(context, limit)
	imgURL, err := gpt.GenerateImage(common.CleanForJSON(context))
	if err != nil {
		fmt.Println("Error generating image: ", err)
		return "", nil
	}
	return imgURL, nil
}

func generateResponse(context, reference string) string {
	prompt := fmt.Sprintf("%s\nUse this as reference: %s\n What will the Narrator say next to progress the situation?\n", context, reference)
	prompt = common.CleanForJSON(prompt)
	resp, err := gpt.Completion(prompt, gpt.GPT_35)
	if err != nil {
		log.Fatal("Error could not process next dialogue: ", err)
	}
	return resp
}

func generateReference(queryIntents []string, context string) string {
	/*
		Builds a reference and returns it
	*/
	reference := strings.Builder{}
	for _, queryIntent := range queryIntents {
		fmt.Println("Exploring intent:", queryIntent)
		if gpt.IsAnswerableFromHistory(queryIntent, context) {
			log.Println("This intent is answerable from the context")
			continue
		} else {
			log.Println("This intent is not answerable from the context")
		}
		resp := resolveIntent(queryIntent, reference.String())
		if len(resp) == 0 {
			log.Println("Reference not found")
		} else {
			log.Println("Reference found: ", resp)
			reference.WriteString(queryIntent + "\n" + resp[:len(resp)-1] + "\n") // getting rid 0f the last comma
		}
	}
	return reference.String()
}

func resolveIntent(queryIntent, reference string) string {
	response := strings.Builder{}
	queryPlan := gpt.GenerateQueryPlan(queryIntent, reference)
	fmt.Println("Query plan generated:", queryPlan)
	result, err := search.Query(*queryPlan)
	if err != nil {
		log.Println("Error fetching reference:", err)
		log.Println("For the query intent", queryIntent)
		log.Println("For the query plan: ", *queryPlan)
		return ""
	}
	for _, binding := range result.Results.Bindings {
		for v := range binding {
			response.WriteString(binding[v].Value + ", ")
		}
	}
	return response.String()
}
