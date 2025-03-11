package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// Content represents the main structure for submitting title and optional description as part of a request body.
type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}

// MyFunctionsArguments represents the arguments required for specific function execution, including a submitter and content.
type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

// BitcoinPriceArguments defines the structure for arguments used to request Bitcoin price in a specific currency.
type BitcoinPriceArguments struct {
	Currency string `json:"currency" jsonschema:"required,description=The currency to get the Bitcoin price in (USD, EUR, GBP, etc)"`
}

// CoinGeckoResponse represents the response structure for Bitcoin price data from the CoinGecko API across multiple currencies.
type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
		EUR float64 `json:"eur"`
		GBP float64 `json:"gbp"`
		JPY float64 `json:"jpy"`
		AUD float64 `json:"aud"`
		CAD float64 `json:"cad"`
		CHF float64 `json:"chf"`
		CNY float64 `json:"cny"`
		KRW float64 `json:"krw"`
		RUB float64 `json:"rub"`
	} `json:"bitcoin"`
}

// main initializes and starts the MCP server, registers tools, prompts, and resources, and handles incoming requests.
func main() {
	log.Println("Starting MCP Server...")

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	// Register "hello" tool
	err := server.RegisterTool("hello", "Say hello to a person with a personalized greeting message", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
		log.Println("Received request for hello tool")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s! Welcome to the MCP Example.", arguments.Submitter))), nil
	})
	if err != nil {
		log.Fatalf("Error registering hello tool: %v", err)
	}

	// Register "bitcoin_price" tool
	err = server.RegisterTool("bitcoin_price", "Get the latest Bitcoin price in various currencies", func(arguments BitcoinPriceArguments) (*mcp_golang.ToolResponse, error) {
		log.Printf("Received request for bitcoin_price tool with currency: %s", arguments.Currency)

		// Default to USD if no currency is specified
		currency := arguments.Currency
		if currency == "" {
			currency = "USD"
		}

		// Call CoinGecko API to get the latest Bitcoin price
		price, err := getBitcoinPrice(currency)
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Error fetching Bitcoin price: %v", err))), nil
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("The current Bitcoin price is %.2f %s (as of %s)",
			price,
			currency,
			time.Now().Format(time.RFC1123)))), nil
	})
	if err != nil {
		log.Fatalf("Error registering bitcoin_price tool: %v", err)
	}

	// Register "prompt_test" prompt
	err = server.RegisterPrompt("prompt_test", "This is a test prompt", func(arguments Content) (*mcp_golang.PromptResponse, error) {
		log.Println("Received request for prompt_test")
		return mcp_golang.NewPromptResponse("description", mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Title)), mcp_golang.RoleUser)), nil
	})
	if err != nil {
		log.Fatalf("Error registering prompt_test: %v", err)
	}

	// Register test resource
	err = server.RegisterResource("test://resource", "resource_test", "This is a test resource", "application/json",
		func() (*mcp_golang.ResourceResponse, error) {
			log.Println("Received request for resource: test://resource")
			return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource(
				"test://resource", "This is a test resource", "application/json",
			)), nil
		})
	if err != nil {
		log.Fatalf("Error registering resource: %v", err)
	} else {
		log.Println("Successfully registered resource: test://resource") // Debug log
	}

	// Start the server
	log.Println("MCP Server is now running and waiting for requests...")
	err = server.Serve()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	select {} // Keeps the server running
}

// getBitcoinPrice retrieves the current Bitcoin price in the specified currency using the CoinGecko API.
// The function returns the price as a float64 and an error if the currency is unsupported or the API call fails.
func getBitcoinPrice(currency string) (float64, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to CoinGecko API
	resp, err := client.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd,eur,gbp,jpy,aud,cad,chf,cny,krw,rub")
	if err != nil {
		return 0, fmt.Errorf("error making request to CoinGecko API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response
	var coinGeckoResp CoinGeckoResponse
	err = json.Unmarshal(body, &coinGeckoResp)
	if err != nil {
		return 0, fmt.Errorf("error parsing JSON response: %w", err)
	}

	// Get price for requested currency
	var price float64
	switch currency {
	case "USD", "usd":
		price = coinGeckoResp.Bitcoin.USD
	case "EUR", "eur":
		price = coinGeckoResp.Bitcoin.EUR
	case "GBP", "gbp":
		price = coinGeckoResp.Bitcoin.GBP
	case "JPY", "jpy":
		price = coinGeckoResp.Bitcoin.JPY
	case "AUD", "aud":
		price = coinGeckoResp.Bitcoin.AUD
	case "CAD", "cad":
		price = coinGeckoResp.Bitcoin.CAD
	case "CHF", "chf":
		price = coinGeckoResp.Bitcoin.CHF
	case "CNY", "cny":
		price = coinGeckoResp.Bitcoin.CNY
	case "KRW", "krw":
		price = coinGeckoResp.Bitcoin.KRW
	case "RUB", "rub":
		price = coinGeckoResp.Bitcoin.RUB
	default:
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}

	return price, nil
}
