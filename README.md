# MCP Example Server for Cursor

**Learn more about how to use this [here](https://www.bytesizego.com/lessons/extending-cursor-mcp-golang)**


This is an example implementation of a Model Context Protocol (MCP) server in Go for use with Cursor.

## Features

- Implements a simple "hello" tool that responds with a greeting
- Provides a "bitcoin_price" tool that fetches real-time Bitcoin prices in various currencies
- Includes a test prompt
- Provides a test resource

## Prerequisites

- Go 1.24 or later
- Cursor IDE
- Internet connection (for the Bitcoin price API)

## Building the Server

```bash
go build -o mcp-example
```

## Installing in Cursor

1. Build the server using the command above
2. Open Cursor settings
3. Navigate to the "MCP Tools" section
4. Click "Add Tool"
5. Select "From Directory"
6. Browse to this directory and select it
7. Click "Install"

## Usage

Once installed, you can use the MCP tools in Cursor by:

1. Opening a chat in Cursor
2. The "hello" tool can be triggered with prompts like "Can you greet me with a personalized message?"
3. The "bitcoin_price" tool can be triggered with prompts like "What's the current Bitcoin price in USD?" or "Show me the Bitcoin price in EUR"

## Development

If you want to modify this example:

1. Edit the `main.go` file to add or modify tools, prompts, or resources
2. Rebuild the server using `go build -o mcp-example`
3. Restart Cursor to load the changes

## Structure

- `main.go` - The main server implementation
- `cursor-mcp-config.json` - Configuration file for Cursor
- `go.mod` and `go.sum` - Go module files

## API Usage

The Bitcoin price tool uses the free CoinGecko API to fetch real-time cryptocurrency prices. No API key is required for basic usage, but there are rate limits.

## License

MIT 