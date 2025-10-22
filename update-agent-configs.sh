#!/bin/bash
# Update all agent .mcp.json configs to use the latest github-mcp-server version

set -e

# Get the latest version tag
LATEST_TAG=$(git describe --tags --abbrev=0)
BINARY_NAME="github-mcp-server-${LATEST_TAG}.exe"
BINARY_PATH="D:/Code/bin/$BINARY_NAME"

echo "=== Updating Agent Configs to $LATEST_TAG ==="
echo

# Check if binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found: $BINARY_PATH"
    echo "Run build-versioned.sh first"
    exit 1
fi

# Update all agent configs
AGENT_WORKSPACE="D:/Code/agent-workspaces"
AGENTS=(agent1 agent2 agent3 agent4 agent5 agentx)

for agent in "${AGENTS[@]}"; do
    MCP_JSON="$AGENT_WORKSPACE/$agent/.mcp.json"

    if [ -f "$MCP_JSON" ]; then
        # Get current version from config
        CURRENT=$(grep -oP 'github-mcp-server-v[0-9]+\.[0-9]+\.[0-9]+\.exe' "$MCP_JSON" || echo "github-mcp-server.exe")

        # Replace with new version
        sed -i "s|\"D:/Code/bin/github-mcp-server.*\.exe\"|\"$BINARY_PATH\"|g" "$MCP_JSON"

        echo "✓ $agent: $CURRENT → $BINARY_NAME"
    else
        echo "⚠ $agent: .mcp.json not found"
    fi
done

echo
echo "=== Update Complete ==="
echo "All agents now use: $BINARY_PATH"
echo
echo "⚠ Reminder: Restart Claude Code for each agent to apply changes"
echo
