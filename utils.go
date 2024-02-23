package main

import (
	"context"
	"fmt"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func get_aad_upn_from_id(id string, client *msgraphsdk.GraphServiceClient) (string, error) {
	result, err := client.Users().ByUserId(id).Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting user by ID: %v\n", err)
		return "", err
	}

	return *result.GetUserPrincipalName(), nil
}

func get_aad_display_name_from_id(id string, client *msgraphsdk.GraphServiceClient) (string, error) {
	result, err := client.Groups().ByGroupId(id).Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting group by ID: %v\n", err)
		return "", err
	}

	return *result.GetDisplayName(), nil
}

func get_aad_ca_named_location_from_id(id string, client *msgraphsdk.GraphServiceClient) (string, error) {
	result, err := client.Identity().ConditionalAccess().NamedLocations().ByNamedLocationId(id).Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting named location by ID: %v\n", err)
		return "", err
	}

	return *result.GetDisplayName(), nil
}

// configureCredentials configures Azure credentials.
func configureCredentials(ctx context.Context) (*azidentity.AzureCLICredential, error) {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating credentials: %v", err)
	}
	return cred, nil
}
