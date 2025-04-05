package main

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Started the MCP Client CLI tool")

	c, err := client.NewStdioMCPClient(
		"npx",
		[]string{}, // ENV
		"-y",
		"@modelcontextprotocol/server-filesystem",
		"/tmp",
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	log.Infoln("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "cli-client",
		Version: "1.0.0",
	}

	initResult, err := c.Initialize(context.Background(), initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	log.Infof(
		"Initialized with server: %s %s\n\n",
		initResult.ServerInfo.Name,
		initResult.ServerInfo.Version,
	)

	log.Infoln("Listing available tools...")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.ListTools(context.Background(), toolsRequest)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		log.Infof("- %s: %s\n", tool.Name, tool.Description)
		log.Infof("Input Schema Type: %s\n", tool.InputSchema.Type)
		log.Infoln("Input Schema Properties:")
		for key, value := range tool.InputSchema.Properties {
			log.Infof("- %s:%v\n", key, value)
		}
		log.Infoln("Input Schema Required Fields:")
		for _, value := range tool.InputSchema.Required {
			log.Infof("- %s\n", value)
		}
	}
}
