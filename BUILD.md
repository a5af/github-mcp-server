# GitHub MCP Server Build Process

## Auto-Versioning Build System

This repository uses an automated build system that automatically increments the patch version on every build.

## Quick Start

### Build New Version

```bash
./build-versioned.sh
```

This will:
1. Get the latest version tag (e.g., `v0.20.2`)
2. Bump patch version (e.g., `v0.20.3`)
3. Build binary as `github-mcp-server-v0.20.3.exe`
4. Create git tag `v0.20.3`
5. Copy binary to `D:/Code/bin/`

### Update All Agent Configs

```bash
./update-agent-configs.sh
```

This will:
1. Get the latest version tag
2. Update all agent `.mcp.json` files to use the new binary
3. Updates: agent1, agent2, agent3, agent4, agent5, agentx

### Push to GitHub

```bash
# Push branch
git push origin feature/automatic-token-recovery

# Push version tag
git push origin v0.20.3
```

## Complete Workflow

```bash
# 1. Make your code changes
git add .
git commit -m "feat: Add new feature"

# 2. Build with auto-version bump
./build-versioned.sh

# 3. Update all agent configs
./update-agent-configs.sh

# 4. Push to GitHub
git push origin feature/automatic-token-recovery
git push origin $(git describe --tags --abbrev=0)

# 5. Restart Claude Code for each agent
```

## Version Scheme

We use semantic versioning: `vMAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes (manual bump)
- **MINOR**: New features (manual bump)
- **PATCH**: Bug fixes, improvements (auto-bump)

### Manual Version Bumps

For major or minor version bumps, manually tag before building:

```bash
# Minor version bump (new feature)
git tag -a v0.21.0 -m "Release v0.21.0: New feature"

# Major version bump (breaking change)
git tag -a v1.0.0 -m "Release v1.0.0: Breaking changes"

# Then build
go build -o github-mcp-server-v0.21.0.exe ./cmd/github-mcp-server
```

## Binary Locations

- **Source builds**: `D:/Code/github-mcp-server/github-mcp-server-v*.exe`
- **Deployed**: `D:/Code/bin/github-mcp-server-v*.exe`
- **Agent configs**: `D:/Code/agent-workspaces/{agent}/.mcp.json`

## Agent Configuration

Each agent's `.mcp.json` points to a specific version:

```json
{
  "mcpServers": {
    "github": {
      "type": "stdio",
      "command": "D:/Code/bin/github-mcp-server-v0.20.2.exe",
      "args": ["stdio"],
      "env": {
        "GITHUB_APP_ID": "...",
        "GITHUB_APP_PRIVATE_KEY_PATH": "...",
        "GITHUB_APP_INSTALLATION_ID": "..."
      }
    }
  }
}
```

## Version History

See all versions:

```bash
git tag -l "v*" --sort=-version:refname
```

Check current version:

```bash
git describe --tags
```

## Troubleshooting

### Build fails with "invalid version tag"

Make sure you have at least one version tag:

```bash
git tag -a v0.20.0 -m "Initial version"
```

### Binary not copied to bin directory

Ensure `D:/Code/bin` exists:

```bash
mkdir -p D:/Code/bin
```

### Agent still using old version

1. Verify config updated: `cat D:/Code/agent-workspaces/agent2/.mcp.json`
2. Restart Claude Code for that agent
3. Check MCP server logs in Claude Code

## Features

### Automatic Token Recovery (v0.20.0+)

- Automatically refreshes GitHub App tokens on 401 errors
- No manual intervention required
- Zero downtime token rotation

### GitHub App Authentication (v0.20.1+)

- `get_me` tool works with GitHub App auth
- Falls back to installation API automatically
- Returns installation account info

## Release Notes

### v0.20.2 (Latest)
- Auto-versioning build system
- Merged upstream changes
- Build and deployment automation

### v0.20.1
- Fixed `get_me` tool for GitHub App authentication
- Falls back to installation API on 401/403
- Returns installation account details

### v0.20.0
- Automatic token recovery on 401 errors
- RefreshingAuthTransport for REST and GraphQL
- GitHub App authentication support

---

**Last Updated**: 2025-10-21
**Maintained By**: AgentX
