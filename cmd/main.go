package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var ErrEmpty = errors.New(fmt.Sprintf("Empty project"))

func getIntentByRequest(projectID, userSentence, languageCode string) (*dialogflowpb.Intent, error) {
	//HEre we create a new SessionsClient and not an IntentsClient, because the SessionsClient struct has a DetectIntent method
	sessionClient, err := dialogflow.NewSessionsClient(context.Background(), option.WithCredentialsFile("kill-yann-olik-af9676297f7c.json"))
	if err != nil {
		return nil, err
	}
	defer sessionClient.Close()

	if projectID == "" {
		return nil, ErrEmpty
	}

	parent := fmt.Sprintf("projects/%s/locations/global/agent/environments/draft/users/-/sessions/456", projectID)

	textInput := dialogflowpb.TextInput{Text: userSentence, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}

	request := dialogflowpb.DetectIntentRequest{Session: parent, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	queryResult := response.GetQueryResult().Intent

	//fullfilmentText := queryResult.GetFulfillmentText()

	return queryResult, nil
}

func getIntents(projectID string) ([]*dialogflowpb.Intent, error) {

	intentsClient, err := dialogflow.NewIntentsClient(context.Background(), option.WithCredentialsFile("kill-yann-olik-af9676297f7c.json"))
	if err != nil {
		return nil, err
	}
	defer intentsClient.Close()

	if projectID == "" {
		return nil, ErrEmpty
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent}

	intentIterator := intentsClient.ListIntents(context.Background(), &request)

	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

func main() {
	/* intents, err := getIntents("kill-yann-olik")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(intents) */

	detectedIntent, err := getIntentByRequest("kill-yann-olik", "Je souhaiterais ingurgiter un lac entier", "fr-FR")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(detectedIntent)
}
