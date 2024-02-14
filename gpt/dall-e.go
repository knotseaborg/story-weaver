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

func GenerateImage(prompt string) (string, error) {
	/*
		GenerateImage generates an image using the DALL-E model based on the provided prompt.

		Parameters:
		  prompt: The prompt used for generating the image.

		Returns:
		  string: The URL of the generated image.
		  error:  An error if the image generation or download process fails.

		Example:
		  imageURL, err := GenerateImage("A cat sitting on a table")
		  if err != nil {
		      log.Fatal("Error generating image:", err)
		  }
		  fmt.Println("Generated image URL:", imageURL)
	*/

	// Generate image
	url := os.Getenv("DALL_E_URL")
	payload := []byte(fmt.Sprintf(`{
		"model":"%s",
		"prompt":"%s",
		"n":1,
		"size":"%s"
	}`, os.Getenv("DALL_E_MODEL"), prompt, os.Getenv("IMG_SIZE")))
	byteContent, err := common.RequestPOST(url, payload)
	if err != nil {
		log.Printf("Payload used: %s", string(payload))
		log.Println("Error generating image: ", err)
		return "", err
	}
	// Download image
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
	/*
		downloadImage() downloads an image from the specified URL and saves it with the given name in the /tmp folder.

		Parameters:
		  url:  The URL of the image to download.
		  name: The name to be used for saving the downloaded image file.

		Returns:
		  error: An error if there is any issue during the download or file writing process.

		Example:
		  err := downloadImage("https://example.com/image.png", "example_image")
		  if err != nil {
		      log.Fatal("Error downloading image:", err)
		  }
	*/
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
