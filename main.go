package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/cloudflare/cloudflare-go/v4/zones"
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
		log.Fatalln("Missing API token")
	}
	zone := os.Getenv("ZONE")

	// Create API client using the scoped API token
	api := cloudflare.NewClient(option.WithAPIToken(token))

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
func getZoneID(zone string, api *cloudflare.Client, ctx context.Context) (string, error) {
	zones, err := api.Zones.List(ctx, zones.ZoneListParams{
		Name: cloudflare.F(zone)},
	)
	if err != nil {
		return "", err
	}
	if len(zones.Result) > 1 {
		err = fmt.Errorf("more than one zone returned")
		return "", err
	}

	return zones.Result[0].ID, nil
}

// Generates an API token
func generateToken(serviceName string, zone string, api *cloudflare.Client, ctx context.Context) (string, error) {
	// Get the Zone ID from the zone name
	zoneID, err := getZoneID(zone, api, ctx)
	if err != nil {
		return "", err
	}

	// Specify token name
	tokenName := serviceName + "." + zone

	// Output input values
	fmt.Println("Generating API token:", tokenName)

	// Specify the API token to create
	tokenToCreate := accounts.TokenNewParams{
		AccountID: cloudflare.F("eb78d65290b24279ba6f44721b3ea3c4"),
		Name:      cloudflare.F("readonly token"),
		Policies: cloudflare.F([]shared.TokenPolicyParam{{
			Effect: cloudflare.F(shared.TokenPolicyEffectAllow),
			// TODO: Pending Cloudflare to fix borked SDK to allow include the ID and Name attributes.
			PermissionGroups: cloudflare.F([]shared.TokenPolicyPermissionGroupParam{{
				// ID:   cloudflare.F("c8fed203ed3043cba015a93ad1616f1f"),
				// Name: cloudflare.F("Zone Read"),
			}, {
				// ID:   cloudflare.F("4755a26eedb94da69e1066d98aa820be"),
				// Name: cloudflare.F("DNS Write"),
			}}),
			Resources: cloudflare.F(map[string]string{
				"com.cloudflare.api.account.zone." + zoneID: "*",
			}),
		}}),
	}

	// Send the request to the Cloudflare API
	generatedToken, err := api.Accounts.Tokens.New(ctx, tokenToCreate)
	if err != nil {
		return "", err
	}

	return generatedToken.Value, nil
}
