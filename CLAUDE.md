# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language Settings
- Default response language: Japanese
- Always respond in Japanese when I write to you in Japanese
- Technical terms and code examples can remain in English when necessary

## Overview

This is a collection of Temporal Go SDK samples demonstrating various workflow patterns and features. Each directory contains a self-contained example with its own workflows, activities, worker, and starter implementations.

## Architecture

### Sample Structure
Each sample follows a consistent pattern:
- **workflows.go** or **workflow.go**: Core workflow definitions
- **activities.go** or **activity.go**: Activity implementations  
- **starter/main.go**: Client that starts workflows
- **worker/main.go**: Worker that executes workflows and activities
- **_test.go files**: Unit tests using Temporal's test suite

### Key Samples
- `helloworld/`: Basic workflow and activity example
- `expense/`: Multi-step approval workflow with activities
- `saga/`: Saga pattern implementation for distributed transactions
- `recovery/`: Error recovery and retry patterns
- `retryactivity/`: Activity retry configurations and patterns
- `child-workflow/`: Parent-child workflow relationships
- `cancellation/`: Workflow cancellation handling
- `dsl/`: DSL-driven workflow execution from YAML

### Temporal Patterns
- Workers register workflows and activities, then listen on task queues
- Starters create clients and execute workflows with specific options
- Task queues (e.g., "hello-world") connect starters and workers
- Tests use `testsuite.WorkflowTestSuite` and `testsuite.NewTestWorkflowEnvironment()`

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./helloworld

# Run specific test
go test -run TestName ./package
```

### Running Samples
Each sample has two components that must run separately:

```bash
# Start worker (in one terminal)
go run ./helloworld/worker

# Execute workflow (in another terminal) 
go run ./helloworld/starter
```

### Build and Dependencies
```bash
# Download dependencies
go mod download

# Build all packages
go build ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Dependencies

- **go.temporal.io/sdk**: Temporal Go SDK for workflows and activities
- **github.com/stretchr/testify**: Testing framework with mocks and assertions
- **gopkg.in/yaml.v3**: YAML parsing for DSL sample
- **github.com/golang/mock**: Mock generation for testing

## Testing Notes

Tests use Temporal's built-in test environment (`testsuite.WorkflowTestSuite`) which provides:
- Workflow execution simulation without running Temporal server
- Activity mocking with `env.OnActivity()`
- DevServer integration for end-to-end testing with `testsuite.StartDevServer()`
