package palm_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	httpGoogle "google.golang.org/api/transport/http"
	"io"
	"net/http"
)

var defaultParameters = &Parameters{
	Temperature:     0.2,
	MaxOutputTokens: 256,
	TopP:            0.8,
	TopK:            40,
}

type PalmClient struct {
	client  *http.Client
	palmUrl string
}

func NewClient(region string, projectId string, modelName string) *PalmClient {
	var err error
	p := &PalmClient{}
	ctx := context.Background()
	p.client, _, err = httpGoogle.NewClient(ctx)
	if err != nil {
		panic("impossible to find a credential")
	}
	p.palmUrl = p.createPalmURL(region, projectId, modelName)
	return p
}

func (p *PalmClient) createPalmURL(region string, projectId string, modelName string) string {
	return fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:predict", region, projectId, region, modelName)
}
func (p *PalmClient) CallPalmApi(prompt string, parameters *Parameters) (response *PalmResponse, err error) {
	//Create the client if not exist
	if p.client == nil {
		ctx := context.Background()
		p.client, _, err = httpGoogle.NewClient(ctx)
		if err != nil {
			panic("impossible to find a credential")
		}
	}

	// Use default parameter by default
	request := PalmRequest{Parameters: defaultParameters}
	if parameters != nil {
		request = PalmRequest{Parameters: parameters}
	}

	request.Instances = append(request.Instances, Content{prompt})

	// Managed JSON, should never fail, ignore the err.
	requestJson, _ := json.Marshal(request)

	//Call API
	rawResponse, err := p.client.Post(p.palmUrl, "application/json", bytes.NewReader(requestJson))
	if err != nil {
		errorMessage := fmt.Sprintf("api call with error %s\n", err)
		fmt.Println(errorMessage)
		err = errors.New(errorMessage)
		return
	}

	if rawResponse.StatusCode == 200 {
		respBody, err := io.ReadAll(rawResponse.Body)
		if err != nil {
			errorMessage := fmt.Sprintf("read body with error %s\n", err)
			fmt.Println(errorMessage)
			err = errors.New(errorMessage)
			return nil, err
		}
		defer rawResponse.Body.Close()

		response = &PalmResponse{}
		err = json.Unmarshal(respBody, response)
		if err != nil {
			errorMessage := fmt.Sprintf("json response parse with error %s\n", err)
			fmt.Println(errorMessage)
			err = errors.New(errorMessage)
			return nil, err
		}
	} else {
		errorMessage := fmt.Sprintf("API call with error %s\n", rawResponse.Status)
		fmt.Println(errorMessage)
		err = errors.New(errorMessage)
		return nil, err
	}
	return
}
