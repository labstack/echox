# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Echo Extra (echox) is a collection of examples and recipes for the Echo web framework in Go. This repository contains:

- **cookbook/**: Complete example applications demonstrating various Echo features (22 examples including JWT, CORS, WebSocket, file upload/download, graceful shutdown, etc.)
- **website/**: Docusaurus-based documentation website

## Development Commands

### Go Development
```bash
# Run tests with race detection
go test -race ./...

# Or use the Makefile
make test

# Run individual cookbook examples
cd cookbook/<example-name>
go run server.go
```

### Website Development
```bash
# Navigate to website directory first
cd website

# Install dependencies
npm install

# Start development server (http://localhost:3000)
npm start

# Build for production
npm run build

# Serve built site
npm run serve

# Alternative: run website in Docker
make serve
# or
docker run --rm -it --name echo-docs -v ${PWD}/website:/home/app -w /home/app -p 3000:3000 -u node node:lts /bin/bash -c "npm install && npm start -- --host=0.0.0.0"
```

## Architecture

### Cookbook Structure
Each cookbook example is a standalone Go application in its own directory under `cookbook/`. Examples follow a consistent pattern:
- `server.go` - Main application entry point using Echo v4
- Additional files for complex examples (handlers, models, etc.)
- Standard Echo patterns: middleware setup, route definitions, handler functions

### Key Dependencies
- **Echo v4** (`github.com/labstack/echo/v4`) - Core web framework
- **JWT libraries** - Multiple JWT implementations for authentication examples
- **WebSocket** (`github.com/gorilla/websocket`) - Real-time communication
- **Go 1.25.1** - Latest stable Go version, aligned with Echo project (supports 1.22+)

### Website Architecture
- **Docusaurus 3.1.0** - Static site generator
- **React 18.2.0** - UI framework
- **MDX** - Markdown with JSX support
- Custom GitHub codeblock theme for syntax highlighting

## Development Workflow

1. **For cookbook examples**: Navigate to specific example directory and run `go run server.go`
2. **For website changes**: Work in `website/` directory using npm commands
3. **Testing**: Use `make test` or `go test -race ./...` to run all Go tests
4. **Local preview**: Use `npm start` in website directory or `make serve` for Docker-based development

## Project Context

This is an examples/recipes repository rather than a library. Focus on:
- Complete, runnable examples in cookbook/
- Clear documentation and code comments
- Following Echo v4 best practices and patterns
- Maintaining consistency across examples