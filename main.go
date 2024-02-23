package main

import (
	"context"
	"fmt"
	"log"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func main() {
	ctx := context.Background()

	// Configure Azure credentials
	cred, err := configureCredentials(ctx)
	if err != nil {
		fmt.Printf("Error configuring credentials: %v\n", err)
		return
	}

	scopes := []string{"https://graph.microsoft.com/.default"}
	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		fmt.Errorf("error creating client: %v", err)
	}
	var policies []models.ConditionalAccessPolicy
	policies, err = getExistingPolicies(graphClient)
	if err != nil {
		log.Fatalf("error getting existing policies: %v", err)
	}

	// create data.tf
	if err := createDataFile(); err != nil {
		fmt.Println("Error creating data file:", err)
		return
	}

	for _, value := range policies {
		create_azurecapolicy(value, graphClient)
	}
	// for _, value := range policies {
	// 	import_policy_to_tfstate(value)
	// }

}
