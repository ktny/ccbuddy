# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CCBuddy is a terminal-based virtual pet application written in Go. The "buddy" is a character that grows and evolves based on ClaudeCode usage, requiring regular interaction to stay healthy.

## Development Commands

```bash
# Core development workflow
make build          # Build the ccbuddy binary
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make lint           # Run comprehensive linting (30+ linters enabled)
make fmt            # Format all Go code

# Development helpers
make dev            # Build and run the application
make clean          # Remove build artifacts
make deps           # Update and verify dependencies
make init-dev       # Install development tools (golangci-lint)

# Direct Go commands (when make unavailable)
go build -o ccbuddy ./cmd/ccbuddy
go test -v ./...
go fmt ./...
golangci-lint run
```

## Architecture

**Current Structure:**
- **cmd/ccbuddy/** - Application entry point and CLI interface
- **internal/buddy/** - Core buddy logic (state management, growth algorithms, health status)
- **internal/tui/** - Terminal UI implementation using bubbletea framework
- **internal/storage/** - Persistence layer for saving buddy state between sessions

**Design Patterns:**
- Standard Go cmd/internal architecture for clear separation
- Bubbletea Model-View-Update pattern for TUI
- JSON-based local file storage for state persistence
- Time-based state evolution requiring periodic updates

## Key Technical Decisions

1. **TUI Framework**: Bubbletea (github.com/charmbracelet/bubbletea) for terminal interface
2. **Go Version**: 1.23.0+ with modern Go features
3. **State Management**: Local JSON file persistence (~/.ccbuddy/)
4. **Code Quality**: Extremely strict linting with 30+ enabled linters
5. **Real-time Growth**: Buddy state changes based on actual ClaudeCode usage

## Core Requirements

From README.md specifications:
- **Buddy Lifecycle**: egg → hatch (via ClaudeCode usage) → growth → potential death (if neglected)
- **Parameters**: appearance, age (time-based), health status (0-100)
- **Feeding Mechanism**: ClaudeCode usage automatically feeds buddy
- **Real-time Nature**: Requires periodic care, health degrades over time
- **Command Interface**: `ccbuddy` command for interaction

## Development Guidelines

**Code Quality Standards:**
- Line length: 120 characters maximum
- Function length: 100 lines, 50 statements maximum
- Cyclomatic complexity: 15 maximum
- All linters in .golangci.yml must pass

**TDD Approach:**
- Write tests first for core buddy logic
- Test state transitions (egg → hatched → growth)
- Test time-based health degradation
- Test persistence layer thoroughly

**Bubbletea Development Pattern:**
```go
type Model struct {
    buddy *buddy.Buddy
    // other state
}

func (m Model) Init() tea.Cmd { /* ... */ }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* ... */ }
func (m Model) View() string { /* ... */ }
```

## Implementation Priority

1. **Core Buddy Model** - State struct, transitions, validation
2. **Persistence Layer** - JSON save/load with error handling
3. **Basic TUI** - Display buddy status, health bar, last fed time
4. **Time-based Logic** - Health degradation, age progression
5. **ClaudeCode Integration** - Usage detection and feeding mechanism

## File Watching Patterns

For ClaudeCode usage detection:
- Monitor file system changes in development directories
- Track git commits/activity as usage indicators
- Consider integration with ClaudeCode session logs
- Implement manual feed command as fallback (`ccbuddy feed`)

## Testing Strategy

Focus on:
- Buddy state machine transitions
- Time-based calculations (age, health degradation)
- Persistence reliability (save/load cycles)
- Edge cases (file corruption, invalid states)
- TUI rendering with various buddy states