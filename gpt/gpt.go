package gpt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/knotseaborg/wikiSearchServer/common"
)

const (
	//GPT_4 = "gpt-4"
	GPT_4  = "gpt-3.5-turbo"
	GPT_35 = "gpt-3.5-turbo"
)

func generatePrompt(keyword, description, content string) string {
	prompt := `Here is a keyword and the description of an entity. 
keyword: <KEYWORD>, description: <DESCRIPTION>

Now here is a list of entities, which consists of ID, label and description.
From this list, return the ID whose description is most satisfactory

<CONTENT>
`
	prompt = strings.Replace(prompt, "<KEYWORD>", keyword, 1)
	prompt = strings.Replace(prompt, "<DESCRIPTION>", description, 1)
	prompt = strings.Replace(prompt, "<CONTENT>", content, 1)

	return prompt
}

func GenerateQueryPlan(queryIntent, reference string) *QueryPlan {
	/*
		Generate the query plan using intent.
		Basically translates the intent into a sparQL template
	*/
	sample, err := os.ReadFile("gpt/prompts/sparQL.txt")
	if err != nil {
		log.Fatal("Error reading file", err)
	}
	if reference == "" {
		reference = "None"
	}
	prompt := fmt.Sprintf("%s\nREFERENCE:%s\nINPUT:%s\nOUTPUT:\n", sample, reference, queryIntent)
	prompt = common.CleanForJSON(prompt)
	resp, err := Completion(prompt, GPT_4)
	var plan QueryPlan
	err = json.Unmarshal([]byte(resp), &plan)
	if err != nil {
		log.Println(resp)
		log.Fatal("Error unmarshalling response: ", err)
	}
	return &plan
}

func IsAnswerableFromContext(queryIntent, context string) bool {
	/*
		Determine if the query intent can be answered from the context itself.
	*/
	prompt := fmt.Sprintf("Reference:\n%s\nCan you guess an answer this question?\n%sRespond with a \"YES\" or \"NO\" and then justify your answer.\n", context, queryIntent)
	prompt = common.CleanForJSON(prompt)
	resp, err := Completion(prompt, GPT_35)
	if err != nil {
		log.Fatal("Error Classfying source")
	}
	if strings.ToLower(resp[:2]) == "no" {
		log.Println("Justification:", resp)
		return false
	}
	return true
}

func GenerateQueryIntents(text string) ([]string, error) {
	content, err := os.ReadFile("gpt/prompts/intent.txt")
	if err != nil {
		log.Fatal(err)
	}
	prompt := fmt.Sprintf("%s\n%s", string(content), text)
	prompt = common.CleanForJSON(prompt)
	resp, err := Completion(prompt, GPT_4)
	if err != nil {
		log.Fatal("Error building query", err)
	}
	intents := []string{}
	for _, intent := range strings.Split(resp, "\n") {
		intent = strings.Trim(intent, " \n")
		if len(intent) > 0 {
			intents = append(intents, intent)
		}
	}
	return intents, nil
}

func Completion(text string, model string) (string, error) {
	// Request payload as a JSON string
	//context = "Hello!\\nThis is a test\\n"
	payload := []byte(fmt.Sprintf(`{
	"model": "%s",
	"temperature": 0.3,
	"messages": [
	  {
	    "role": "user",
	    "content": "%s"
	  }
	]
	}`, model, text))

	byteContent, err := common.RequestPOST(os.Getenv("GPT_35_URL"), payload)
	if err != nil {
		log.Fatal("Error while GPT completion", err)
	}
	var comp completion
	err = json.Unmarshal(byteContent, &comp)
	if err != nil {
		log.Println("Error completion: ", err)
		return "", err
	}
	return comp.Choices[0].Message.Content, nil
}
