## [Unreleased]

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
