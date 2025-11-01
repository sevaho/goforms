# CLAUDE.md

## AI Guidance

- To save main context space, for code searches, inspections, troubleshooting or analysis, use code-searcher subagent where appropriate - giving the subagent full context background for the task(s) you assign it.
- After receiving tool results, carefully reflect on their quality and determine optimal next steps before proceeding. Use your thinking to plan and iterate based on this new information, and then take the best next action.
- For maximum efficiency, whenever you need to perform multiple independent operations, invoke all relevant tools simultaneously rather than sequentially.
- Before you finish, please verify your solution
- Do what has been asked; nothing more, nothing less.
- NEVER create files unless they're absolutely necessary for achieving your goal.
- ALWAYS prefer editing an existing file to creating a new one.
- NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.
- When you update or modify core context files, also update markdown documentation and memory bank
- When asked to commit changes, exclude CLAUDE.md and CLAUDE-*.md referenced memory bank system files from any commits. Never delete these files.

## Project Overview

GoForms is an open-source, self-hosted alternative to AirForms, JotForms, and FormSpree. It allows users to send form data as email without backend code via HTTP forms, with support for multiple languages (Dutch, French, English) and spam protection via Google reCAPTCHA. All mails are encrypted if stored in database.

## Architecture

The application follows a modular Go architecture:

- **Echo Framework**: HTTP server with middleware for security, recovery, and live reload
- **PostgreSQL**: Database with pgx/v5 driver and sqlc for type-safe queries
- **Repository Pattern**: Data access layer with encrypted mail storage
- **Mail Providers**: Abstracted mail sending (MailerSend, Fake provider for testing)
- **Template Engine**: Custom render engine with live reload support, html files
- **Configuration**: YAML-based config with environment variable overrides

### Key Components

- `src/app/`: Main application logic, handlers, and routing
- `src/db/`: Generated database queries and models (via sqlc)
- `src/repository/`: Data access layer with encryption
- `src/mailproviders/`: Mail provider implementations
- `src/config/`: Configuration management
- `src/pkg/`: External service integrations (reCAPTCHA, Telegram)

## Development Commands

### Essential Commands
- `make run` - Start development server with air (hot reload)
- `make test` - Run all tests with cache clearing
- `make lint` - Run golangci-lint with all rules enabled
- `make sqlgen` - Generate SQL code from queries using sqlc
- `make migrate` - Run database migrations
- `make css` - Watch and build Tailwind CSS

### Database
- `make sqlgenr` - Auto-regenerate SQL on file changes
- `go run . --migrate` - Run migrations directly
- `go run . --new-migration` - Create new migration file

### Testing
- `make testr` - Auto-run tests on file changes
- `gotest -v ./...` - Run tests with verbose output

### Deployment
- `make serve` - Run via Docker locally
- `make deploy` - Deploy to Kubernetes (requires KO_DOCKER_REPO)
- `make compose` - Run full stack with docker-compose

## Testing Strategy

The codebase uses:
- Table-driven tests with testify/suite
- HTTP mocking with httpmock
- Separate test setup per package
- In-memory test database for integration tests
- Use Ginko go testing

## Key Dependencies

- **Echo v4**: Web framework
- **pgx/v5**: PostgreSQL driver
- **sqlc**: Type-safe SQL generation
- **dbmate**: Database migrations
- **air**: Live reload for development
- **Tailwind CSS**: Styling framework
- **fernet-go**: Encryption for sensitive data

## Golang Code Guidelines

- Be explicit
- Group imports: standard library first, then third-party
- Use PascalCase for exported types/methods, camelCase for variables
- Add comments for public API and complex logic
