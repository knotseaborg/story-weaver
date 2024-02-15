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

func Run(scene string) {
	/*
		Run executes the interactive storytelling game loop.

		The function initializes a string builder to track the history of interactions.
		It sets the initial context for the storytelling scenario.
		Inside a loop, it prompts the user for their action and reads the input.
		If the input is "stop", "bye", or "exit", the function returns and ends the game.
		Otherwise, it generates query intents based on the updated context and user input using the gpt.GenerateQueryIntents function.
		It then generates a reference string based on the generated query intents and the history of interactions.
		The function updates the context with the user input and generates the next response based on the updated context and reference.
		Additionally, it attempts to generate an image based on the updated context using the generateImage function.
		If successful, it displays the URL of the generated image.
		Finally, it displays the next prompt to continue the storytelling scenario.

		Parameters:
		  scene: A string which sets the scene for the story

		Example:
		  Run("")
	*/
	defer func() { fmt.Println("Thank you playing!") }()
	history := strings.Builder{}
	context := "Narrator: " + "One of the lead actresses of Avatar 2 appeared to be preparing for another movie."
	fmt.Println(context)
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
	}
}

func generateImage(context string) (string, error) {
	/*
		generateImage() generates an image based on the provided context.

		Parameters:
		  context: The context string used to generate the image.

		Returns:
		  string: The URL of the generated image.
		  error:  An error if there is any issue during the image generation process.

		Example:
		  context := "A cat sitting on a table"
		  imageURL, err := generateImage(context)
		  if err != nil {
		      log.Fatal("Error generating image:", err)
		  }
		  fmt.Println("Generated image URL:", imageURL)
	*/
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
	/*
		generateResponse() generates a response for progressing a situation based on the provided context and reference.

		Parameters:
		  context:   The context string describing the current situation.
		  reference: The reference string containing relevant information for generating the response.

		Returns:
		  string: The generated response containing the next dialogue to progress the situation.

		Example:
		  context := "In a cafe, two friends are having a conversation."
		  reference := "Friend 1: What do you want to order? Friend 2: I'll have a coffee, please."
		  response := generateResponse(context, reference)
		  fmt.Println("Generated response:", response)
	*/
	prompt := fmt.Sprintf("%s\nUse this as reference: %s\n What will the Narrator say next to progress the situation?\n", context, reference)
	prompt = common.CleanForJSON(prompt)
	resp, err := gpt.Completion(prompt, os.Getenv("GPT_MODEL_BASIC"))
	if err != nil {
		log.Fatal("Error could not process next dialogue: ", err)
	}
	return resp
}

func generateReference(queryIntents []string, context string) string {
	/*
		generateReference() builds a reference based on the provided query intents and context.

		Parameters:
		  queryIntents: A slice of query intents to generate reference for.
		  context:      The context string used to determine if a query intent can be answered from the context.

		Returns:
		  string: The constructed reference string containing resolved responses for each query intent.

		Example:
		  queryIntents := []string{"Find companies headquartered in New York", "Find restaurants in Paris"}
		  context := "A city in the United States"
		  reference := generateReference(queryIntents, context)
		  fmt.Println("Generated reference:", reference)
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
		resp := resolveIntent(queryIntent, reference.String(), 3)
		if len(resp) == 0 {
			log.Println("Reference not found")
		} else {
			log.Println("Reference found: ", resp)
			reference.WriteString(queryIntent + "\n" + resp[:len(resp)-1] + "\n") // getting rid 0f the last comma
		}
	}
	return reference.String()
}

func resolveIntent(queryIntent, reference string, maxTries int) string {
	/*
		resolveIntent() resolves the query intent into a response by executing a query plan.

		Parameters:
		queryIntent: The query intent to resolve.
		reference:   The reference string used in generating the query plan.

		Returns:
		string: The resolved response constructed from the query result bindings.

		Example:
		queryIntent := "Find companies headquartered in New York"
		reference := "Companies headquartered in New York"
		resolvedResponse := resolveIntent(queryIntent, reference)
		fmt.Println("Resolved response:", resolvedResponse)
	*/
	if maxTries == 0 {
		return ""
	}
	response := strings.Builder{}
	queryPlan := gpt.GenerateQueryPlan(queryIntent, reference)
	fmt.Println("Query plan generated:", queryPlan)
	result, err := search.ExecuteQuery(*queryPlan)
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
	if response.Len() == 0 {
		reference := reference +
			fmt.Sprintf("\n[Please Try again. Your previous query and`%s` returned no response from Wikidata. Please rebuild a query plan with different structure]",
				queryPlan.Query,
			)
		return resolveIntent(queryIntent, reference, maxTries-1)
	}
	return response.String()
}
