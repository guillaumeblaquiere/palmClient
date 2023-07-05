package main

import (
	"fmt"
	palm_client "golang.com/gblaquiere/palm-client"
)

func main() {

	region := "us-central1"
	modelName := "text-bison@001"
	projectId := "gdglyon-cloudrun"

	prompt := "Write poem on cats"

	palmClient := palm_client.NewClient(region, projectId, modelName)

	// Use the default parameters
	response, err := palmClient.CallPalmApi(prompt, nil)

	if err != nil {
		fmt.Printf("error during the API call: %s\n", err)
		return
	}

	// You can use your own parameters if you prefer
	/*
		myParameters := &palm_client.Parameters{
			Temperature:     0.2,
			MaxOutputTokens: 256,
			TopP:            0.8,
			TopK:            40,
		}

		response, err := palmClient.CallPalmApi(prompt, myParameters)
	*/

	fmt.Printf("The initial prompt is \n%s\n\n", prompt)
	fmt.Printf("The generated answer is \n%s\n", response.Predictions[0].Content)
	return
}
