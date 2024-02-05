package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/knotseaborg/wikiSearchServer/common"
)

const (
	DALLE_E = "dall-e-3"
)

func GenerateImage(prompt string) (string, error) {
	url := os.Getenv("DALL_E_URL")
	payload := []byte(fmt.Sprintf(`{
		"model":"%s",
		"prompt":"%s",
		"n":1,
		"size":"%s"
	}`, DALLE_E, prompt, os.Getenv("IMG_SIZE")))
	byteContent, err := common.RequestPOST(url, payload)
	if err != nil {
		log.Printf("Payload used: %s", string(payload))
		log.Println("Error generating image: ", err)
		return "", err
	}
	var img image
	json.Unmarshal(byteContent, &img)

	err = downloadImage(img.Data[0].URL, fmt.Sprintf("%d", img.Created))
	if err != nil {
		log.Fatal("Error while downloading image", err)
		return "", err
	}
	return img.Data[0].URL, nil
}

func downloadImage(url, name string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	err = os.WriteFile(fmt.Sprintf("tmp/%s.png", name), body, 0644)
	if err != nil {
		return err
	}
	return nil
}
