package controller

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
	"github.com/knotseaborg/wikiSearchServer/gpt"
	"github.com/knotseaborg/wikiSearchServer/search"
)

func Run() {
	defer func() { fmt.Println("Thank you playing!") }()
	context := "Narrator: The lead actress of Avatar 2 appeared to be preparing for another movie."
	//text := "A film like nothing before. I see several large animatronic dinosaurs though."
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println(context, "\nWhat would you do?")
		text := readInput(reader)
		if text == "stop" || text == "bye" || text == "exit" {
			return
		}
		queryIntents, err := gpt.GenerateQueryIntents(context + text)
		if err != nil {
			log.Println("Error generating query intents", err)
		}
		reference := fetchReference(queryIntents, context)
		context = fmt.Sprintf("%s\nUser:%s", context, text)
		prompt := fmt.Sprintf("%s\nUse this as reference: %s\n What happens next?\n", context, reference)
		prompt = common.CleanForJSON(prompt)
		nextText, err := gpt.Completion(prompt, gpt.GPT_35)
		if err != nil {
			log.Fatal("Error could not process next dialogue: ", err)
		}
		context = fmt.Sprintf("%s\n%s", context, nextText)
		fmt.Printf("---------------WHAT HAPPENS NEXT-----------------\n%s\n", nextText)
		imgURL, err := gpt.GenerateImage(common.CleanForJSON(context))
		if err == nil {
			fmt.Printf("---------------CLICK TO SEE THE IMAGE-----------------\n%s\nf", imgURL)
		}
	}
}

func readInput(reader *bufio.Reader) string {
	// Read a line of text including spaces
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return text
}

func fetchReference(queryIntents []string, context string) string {
	/*
		Builds a reference and returns it
	*/
	reference := strings.Builder{}
	for _, queryIntent := range queryIntents {
		fmt.Println("Exploring intent:", queryIntent)
		if gpt.IsAnswerableFromContext(queryIntent, context) {
			log.Println("This intent is answerable from the context")
			continue
		} else {
			log.Println("This intent is not answerable from the context")
		}
		buf := strings.Builder{}
		queryPlan := gpt.GenerateQueryPlan(queryIntent, reference.String())
		fmt.Println("Query plan generated:", queryPlan)
		result, err := search.Query(*queryPlan)
		if err != nil {
			log.Println("Error fetching reference:", err)
			continue
		}
		for _, binding := range result.Results.Bindings {
			for v := range binding {
				buf.WriteString(binding[v].Value + ", ")
			}
		}
		if buf.Len() > 0 {
			log.Println("Reference found: ", buf.String())
			reference.WriteString(queryIntent + "\n" + buf.String()[:buf.Len()-1] + "\n") // getting rid 0f the last comma
		} else {
			log.Println("Reference not found")
		}
	}
	return reference.String()
}
