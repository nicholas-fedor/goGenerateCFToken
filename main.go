package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
)

// Load .env file
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {
	// Get the service name from the command-line argument
	serviceName := strings.ToLower(os.Args[1])
	switch serviceName {
	case "":
		log.Fatalln("Missing service name from command-line argument")
	case "help":
		log.Fatalln("Usage: goGenerateCFToken [service name]")
	}

	// Set token and zone from loaded env variables
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		log.Fatalln("Missing API Token")
	}
	zone := os.Getenv("ZONE")

	// Create API client using the scoped API token
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Generate an API token
	newApiToken, err := generateToken(serviceName, zone, api, ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print the generated API token's value
	fmt.Println(newApiToken)
}

// Returns the Zone ID
func getZoneID(zone string, api *cloudflare.API, ctx context.Context) (string, error) {
	zones, err := api.ListZones(ctx, zone)
	if err != nil {
		return "", err
	}
	if len(zones) > 1 {
		err = fmt.Errorf("more than one zone returned")
		return "", err
	}

	return zones[0].ID, nil
}

// Generates an API token
func generateToken(serviceName string, zone string, api *cloudflare.API, ctx context.Context) (string, error) {
	// Get the Zone ID from the zone name
	zoneID, err := getZoneID(zone, api, ctx)
	if err != nil {
		return "", err
	}

	// Specify token name
	tokenName := serviceName + "." + zone

	// Output input values
	fmt.Println("Generating API token:", tokenName)

	// Specify API token to create
	resources := make(map[string]interface{})
	resources["com.cloudflare.api.account.zone."+zoneID] = "*"
	tokenToCreate := cloudflare.APIToken{
		Name: tokenName,
		Policies: []cloudflare.APITokenPolicies{{
			Effect:    "allow",
			Resources: resources,
			PermissionGroups: []cloudflare.APITokenPermissionGroups{
				{
					ID:   "c8fed203ed3043cba015a93ad1616f1f",
					Name: "Zone Read",
				},
				{
					ID:   "4755a26eedb94da69e1066d98aa820be",
					Name: "DNS Write",
				},
			},
		}},
	}

	// Send the request to the Cloudflare API
	generatedToken, err := api.CreateAPIToken(ctx, tokenToCreate)
	if err != nil {
		return "", err
	}

	return generatedToken.Value, nil
}
