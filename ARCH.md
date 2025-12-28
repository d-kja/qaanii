# Architectural Design Documentation

## Primary Architectural Patterns

### 1. Clean Architecture / Hexagonal Architecture (Ports & Adapters)

The project follows a strict layered approach:

- **Domain Layer** (`internals/domain/`)
  - Business logic via use cases
  - No infrastructure dependencies
  - Pure domain entities in shared module

- **Infrastructure Layer** (`internals/infra/`)
  - HTTP handlers (Fiber framework)
  - Message broker adapters (RabbitMQ)
  - Database layer (planned/minimal)

- **Shared Kernel** (`shared/`)
  - Cross-service domain entities (Manga, Chapter, Page)
  - Event definitions and broker abstractions
  - Common utilities

### 2. Microservices Architecture

Two distinct services in a monorepo:

**Scraper Service** (`/scraper/`):
- Handles web scraping via browser automation
- Publishes events to RabbitMQ
- Fully implemented with Rod (browser automation)

**Manga Service** (`/manga/`):
- Primary API service
- Exposes REST endpoints
- Consumes scraper events (planned)
- Currently has stub implementations

### 3. Event-Driven Architecture

Asynchronous inter-service communication via RabbitMQ:

**Event Flow Pattern:**
```
Request Event → RabbitMQ → Scraper Processes →
Result Event → RabbitMQ → Manga Service (planned)
```

**Event Types:**
- `SEARCH_MANGA_EVENT` → `SEARCHED_MANGA_EVENT`
- `SCRAPE_MANGA_EVENT` → `SCRAPED_MANGA_EVENT`
- `SCRAPE_CHAPTER_EVENT` → `SCRAPED_CHAPTER_EVENT`

## Design Patterns Implemented

### Use Case Pattern (Application Services)
Every business operation is a dedicated service:

```go
type SearchByNameService struct {
    Scraper coreentities.Scraper
}

func (self *SearchByNameService) Exec(request SearchByNameRequest) (*SearchByNameResponse, error)
```

Located in: `internals/domain/*/use_case/*.usecase.go`

### Factory Pattern
Scraper creation: `shared/broker/channels/scraper.entity.go:37`

### Adapter Pattern
HTTP handlers adapt Fiber framework to domain use cases: `internals/infra/http/handlers/*.go`

### Pub/Sub Pattern
Message broker with publishers and subscribers: `shared/broker/`

### Strategy Pattern
Retry mechanisms with different strategies:
```go
RETRY, RETRY_MANY, RETRY_XPATH, RETRY_XPATH_MANY
```

### Template Method Pattern
Base event structure with metadata generation: `shared/broker/events/base.event.go:12`

## Key Architectural Decisions

### Dependency Management
- **Go Workspaces**: Multi-module monorepo (shared, manga, scraper)
- **Dependency Injection**: Via struct fields in use cases
- **No DI Framework**: Manual injection keeping it simple

### Data Flow

**HTTP Request Flow:**
```
HTTP Request → Handler → Use Case → Scraper Entity → External Site
```

**Event-Driven Flow:**
```
Manga Service → Publish Event → RabbitMQ →
Scraper Consumes → Execute Use Case → Publish Result →
Manga Service (planned)
```

### Separation of Concerns
- Handlers have zero business logic (pure delegation)
- Use cases are framework-agnostic
- Domain constants externalized (XPath queries in constants files)
- Infrastructure details completely abstracted from domain

### Error Handling Strategy
- Explicit error returns (Go idiom: `(result, error)`)
- Graceful degradation in scraping (continues on optional field failures)
- Retry mechanisms with configurable attempts
- Comprehensive logging with context prefixes
- Resource cleanup via `defer`

## Domain-Driven Design Elements

### Bounded Contexts
- **Scraper Context**: Web scraping operations, browser automation
- **Manga Context**: Catalog management, API exposition, persistence
- **Shared Kernel**: Common entities and event contracts

### Entities (Anemic Model)
Pure data structures without behavior:
- `Manga`, `Chapter`, `Page`
- JSON serialization tags
- Relationships via composition

### Aggregates
`Manga` as aggregate root containing `Chapters` → `Pages` hierarchy

### Ubiquitous Language
Consistent terminology: Manga, Chapter, Page, Scraper, Slug, Publisher, Subscriber

## Technology Stack

- **Language**: Go 1.25.5
- **Web Framework**: Fiber v2 (Express-like)
- **Message Broker**: RabbitMQ (AMQP 0.9.1)
- **Web Scraping**: Rod (Chrome DevTools Protocol) + Stealth plugin
- **Database**: SQLite (planned, minimal implementation)
- **Development**: Air (hot reload), Docker Compose

## Project Structure

```
backend/
├── shared/              # Shared kernel
│   ├── entities/        # Domain entities
│   ├── broker/          # Event infrastructure
│   └── utils/           # Common utilities
├── manga/               # API microservice
│   ├── cmd/             # Entry point
│   └── internals/
│       ├── domain/      # Business logic (use cases)
│       └── infra/       # Infrastructure (HTTP, DB, broker)
└── scraper/             # Scraping microservice
    ├── cmd/
    └── internals/
        ├── domain/      # Scraping logic
        └── infra/       # Infrastructure
```

## Configuration Management

- **Environment Variables**: `.env` files per service via godotenv
- **Infrastructure as Code**: Docker Compose for RabbitMQ
- **Constants**: Externalized domain constants (XPath queries, event names)
- **Fail-Fast**: Application exits on missing configuration

## Communication Patterns

### Intra-Service Communication
```
Handler → Use Case → Infrastructure Component
```

### Inter-Service Communication

**Publisher Registration** (on startup):
```go
broker.SetupPublishers(...)  // Creates publishers in context
```

**Subscriber Registration** (on startup):
```go
broker.SetupSubscribers(...)  // Starts consumer goroutines
```

**Message Flow**:
```
Manga Service → Search Request Event → RabbitMQ Queue →
Scraper Subscriber → Process → Publish Result Event →
Manga Service Subscriber (planned)
```

**Event Naming Convention:**
- Request events: `{ACTION}_{ENTITY}_EVENT` (e.g., `SEARCH_MANGA_EVENT`)
- Response events: `{ACTION}ED_{ENTITY}_EVENT` (e.g., `SEARCHED_MANGA_EVENT`)

**Idempotency:**
```go
func (self *BaseEvent) GenerateEventId(event_type string, user_id string) string {
    random_id := uuid.New().String()
    idempotency := fmt.Sprintf("%v-%v-%v", event_type, user_id, random_id)
    self.Metadata.Id = idempotency
    return idempotency
}
```

## Architectural Quality Assessment

### Strengths

1. **Strong Separation of Concerns**
   - Clear domain/infrastructure split
   - Use case pattern enforces single responsibility
   - Infrastructure details abstracted

2. **Testability**
   - Use cases are pure functions
   - Dependencies injected
   - No global state

3. **Scalability**
   - Microservices can scale independently
   - Async communication via message broker
   - Stateless services

4. **Maintainability**
   - Consistent project structure across services
   - Clear naming conventions
   - Well-organized directory hierarchy

5. **Resilience**
   - Retry mechanisms for scraping
   - Graceful degradation
   - Error logging without crashes
   - Resource cleanup via defer

### Areas for Improvement

1. **Database Integration**
   - Manga service has stub implementations
   - Database layer incomplete
   - No persistence layer implementation

2. **Error Handling**
   - Could benefit from custom error types
   - No error wrapping/context in some areas
   - Missing validation layer

3. **Broker Integration in Manga Service**
   - Empty broker directory
   - No consumer implementation yet
   - Events published but not consumed

4. **Testing**
   - No visible test files
   - No integration tests
   - Test infrastructure not set up

5. **Documentation**
   - Limited inline documentation
   - No package-level docs
   - API documentation missing

## Summary

This is a **well-architected microservices application** following **Clean Architecture** principles with:
- Strong separation between domain and infrastructure
- Event-driven async communication
- Use case-centric business logic
- Clear bounded contexts with a shared kernel
- Resilient scraping with retry mechanisms
- Scalable, testable, and maintainable design

The architecture is still **in active development** with the manga service's database layer and broker consumers pending implementation.
