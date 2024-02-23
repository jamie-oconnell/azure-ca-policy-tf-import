package main

import (
	"context"
	"fmt"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// printOdataError prints OData errors if they occur.
func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("OData Error: %s\n", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("Code: %s\n", *terr.GetCode())
			fmt.Printf("Message: %s\n", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > Error: %#v\n", err, err)
	}
}

// configureCredentials configures Azure credentials.
func configureCredentials(ctx context.Context) (*azidentity.AzureCLICredential, error) {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating credentials: %v", err)
	}
	return cred, nil
}

// fetchExistingPolicies fetches existing conditional access policies.
func fetchExistingPolicies(ctx context.Context, cred *azidentity.AzureCLICredential) ([]models.ConditionalAccessPolicy, error) {
	scopes := []string{"https://graph.microsoft.com/.default"}
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	result, err := client.Identity().ConditionalAccess().Policies().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting CA policies: %v", err)
	}

	// Use PageIterator to iterate through all conditional access policies
	pageIterator, err := msgraphcore.NewPageIterator[*models.ConditionalAccessPolicy](result, client.GetAdapter(), models.CreateConditionalAccessPolicyFromDiscriminatorValue)
	if err != nil {
		return nil, fmt.Errorf("error creating page iterator: %v", err)
	}

	var policies []models.ConditionalAccessPolicy

	err = pageIterator.Iterate(ctx, func(capolicy *models.ConditionalAccessPolicy) bool {
		policies = append(policies, *capolicy)
		// Return true to continue the iteration
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating through policies: %v", err)
	}

	return policies, nil
}

func getExistingPolicies() ([]models.ConditionalAccessPolicy, error) {
	ctx := context.Background()

	// Configure Azure credentials
	cred, err := configureCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("error configuring credentials: %v", err)
	}

	// Fetch existing policies
	policies, err := fetchExistingPolicies(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("error fetching policies: %v", err)
	}

	return policies, nil
}
