## [Unreleased]


## [0.20.6] - 2025-10-22

### Fixed
- `get_me` now correctly detects GitHub App authentication by checking installation ID first
- Prevents GitHub App tokens from incorrectly returning OAuth user identity
- Some installation tokens can successfully call `/user` endpoint but return wrong identity (installing user instead of app account)

### Technical
- Changed detection logic: check `installationID > 0` before attempting `Users.Get`
- Ensures GitHub App installations always use `ListRepos` path regardless of token capabilities
- Updated comments to document the authentication path selection logic

### Background
The previous v0.20.5 fix used error-based detection (401/403 from Users.Get), but some GitHub App installation tokens can successfully call the `/user` endpoint. When they do, they return the installing user's identity (e.g., "a5af") instead of the app installation account. This fix ensures we always use the correct authentication path when we know we have a GitHub App (installationID > 0).


## [0.20.5] - 2025-10-22

### Fixed
- `get_me` tool now uses `ListRepos` API endpoint instead of `GetInstallation` for GitHub App authentication
- Resolves 403 "Resource not accessible by integration" errors with installation access tokens
- GitHub App installations can now retrieve account details without requiring JWT authentication

### Technical
- Replaced `client.Apps.GetInstallation(ctx, installationID)` with `client.Apps.ListRepos(ctx, nil)`
- Updated `GetMe` function in `context_tools.go` to extract owner info from first accessible repository
- Added automated tests for GitHub App authentication in `context_tools_githubapp_test.go`
- All existing tests pass with the new implementation

### Background
GitHub App authentication uses installation access tokens which have different API access than JWT tokens:
- JWT tokens: Can call `GetInstallation` endpoint
- Installation tokens: Can call `ListRepos` but not `GetInstallation`
This fix aligns with proper GitHub App authentication patterns.

## [0.20.4] - 2025-10-22

### Fixed
- `get_me` tool now correctly retrieves GitHub App installation info using actual installation ID
- Fixed 403 error when calling `get_me` with GitHub App authentication  
- Installation ID properly passed from auth provider through to tools

### Technical
- Added `GetInstallationID()` method to `GitHubAppAuthProvider`
- Updated `DefaultToolsetGroup` to accept and pass installation ID
- All tests pass with the fix

## [0.20.3] - 2025-10-22

### Added
- Automatic token recovery on 401 errors with `RefreshingAuthTransport`
- Request body size limits (10MB) for retry safety
- Infinite loop protection for token refresh retries
- Enhanced error messages with request context (method + path)
- Auto-versioning build system with `build-versioned.sh`
- Agent configuration update script `update-agent-configs.sh`
- Comprehensive build documentation in `BUILD.md`
- Comprehensive CHANGELOG.md

### Changed
- Improved `get_me` tool description for GitHub App authentication
- Better logging with request details
- Toolset management functions moved to `pkg/github/tools.go`

### Fixed
- GitHub App token expiration no longer requires server restart
- Prevents retry on oversized request bodies (>10MB)
- Stops infinite retry loops with second 401 responses
- Updated tool snapshots for compatibility

# Changelog

All notable changes to GitHub MCP Server will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.20.2] - 2025-10-21

### Added
- Auto-versioning build system
- Merged upstream changes from main branch
- Build and deployment automation scripts

### Changed
- Updated build process to automatically bump patch versions
- Improved agent configuration management

## [0.20.1] - 2025-10-21

### Fixed
- `get_me` tool now works correctly with GitHub App authentication
- Falls back to installation API on 401/403 responses
- Returns installation account details for GitHub Apps

## [0.20.0] - 2025-10-21

### Added
- Automatic GitHub App token recovery on authentication failures
- `RefreshingAuthTransport` for REST and GraphQL clients
- Background token refresh mechanism
- Zero-downtime token rotation

### Changed
- Enhanced GitHub App authentication support
- Improved error handling for expired tokens

## [0.19.0] - Previous releases

See [GitHub Releases](https://github.com/github/github-mcp-server/releases) for earlier versions.

---

## Upgrade Guide

### From 0.20.x to Current

No breaking changes. The automatic token recovery is backward compatible with existing configurations.

**Recommended**: Restart your MCP server to activate the new token recovery features.

### Configuration Changes

No configuration changes required. Existing `GITHUB_APP_*` environment variables work as-is.

---

## Breaking Changes

None in current release.

---

## Deprecations

None in current release.
