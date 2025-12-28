# Qaanii

A local manga reading platform that scrapes manga content from the web and provides a clean API for consumption by client applications.

## Overview

Qaanii is designed to give you full control over your manga reading experience by maintaining a local database of manga content. The system consists of two main services that communicate via RabbitMQ message broker:

- **Scraper Service**: Web scraping engine that extracts manga data using headless browser automation
- **Manga API Service**: RESTful API that serves manga data to client applications
- **Shared Module**: Common entities and utilities used across services

## Architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Client    │ ◄─────► │  Manga API   │ ◄─────► │   SQLite    │
│ Application │         │   Service    │         │  Database   │
└─────────────┘         └──────────────┘         └─────────────┘
                               │
                               │ RabbitMQ
                               ▼
                        ┌──────────────┐
                        │   Scraper    │
                        │   Service    │
                        └──────────────┘
                               │
                               ▼
                        ┌──────────────┐
                        │  Web Sources │
                        └──────────────┘
```

## Features

- Web scraping with headless browser automation (Rod + Stealth)
- RESTful API for manga data access
- Asynchronous scraping via message queue
- SQLite database for local storage
- Base64 image storage (no blob storage required)
- Search functionality
- Chapter and page management

## Tech Stack

### Backend
- **Language**: Go 1.25.5
- **HTTP Framework**: Fiber v2
- **Web Scraping**: go-rod (headless browser automation)
- **Message Broker**: RabbitMQ 4.2.2
- **Database**: SQLite (migrated from PostgreSQL)
- **Environment Config**: godotenv

### Infrastructure
- Docker Compose for service orchestration
- Air for hot-reload during development

## Project Structure

```
qaanii/
├── backend/
│   ├── scraper/          # Web scraping service
│   │   ├── cmd/          # Application entry point
│   │   └── internals/    # Internal implementation
│   │       ├── domain/   # Business logic
│   │       └── infra/    # Infrastructure (HTTP handlers, DB)
│   ├── manga/            # Main API service
│   │   ├── cmd/          # Application entry point
│   │   └── internals/    # Internal implementation
│   │       ├── domain/   # Business logic
│   │       └── infra/    # Infrastructure (HTTP handlers, DB)
│   ├── shared/           # Shared code between services
│   │   ├── entities/     # Common data structures
│   │   └── utils/        # Shared utilities
│   ├── docker-compose.yml
│   └── go.work           # Go workspace configuration
└── client/               # (Planned) Client application
```

## Data Models

### Manga
```go
{
  "url": "string",           // Original source URL
  "slug": "string",          // URL-friendly identifier
  "name": "string",
  "description": "string",
  "tags": ["string"],
  "image": "string",         // Base64 encoded
  "image_type": "string",
  "status": "string",
  "last_update": "string",
  "chapters": [Chapter]
}
```

### Chapter
```go
{
  "title": "string",
  "link": "string",
  "time": "string",
  "pages": [Page]
}
```

### Page
```go
{
  "order": int,
  "image": "string",         // Base64 encoded
  "image_type": "string"
}
```

## Getting Started

### Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose
- Air (optional, for development)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd qaanii
```

2. Set up environment variables
```bash
cd backend
cp .env.example .env
# Edit .env with your configuration
```

3. Start the RabbitMQ broker
```bash
docker-compose up -d
```

4. Install dependencies
```bash
go work sync
```

5. Run the services

**Scraper Service:**
```bash
cd backend/scraper
go run cmd/main.go
# Or with hot-reload:
air
```

**Manga API Service:**
```bash
cd backend/manga
go run cmd/main.go
# Or with hot-reload:
air
```

## Environment Variables

Copy the `.env.example` file to `.env` and configure:

```env
# RabbitMQ Configuration
RABBITMQ_DEFAULT_USER=guest
RABBITMQ_DEFAULT_PASS=guest
RABBITMQ_DEFAULT_VHOST=/

# Database Configuration (Legacy - Now using SQLite)
MANGA_DB_NAME=mangas
MANGA_DB_USER=docker
MANGA_DB_PASSWORD=docker
```

## API Endpoints

### Search
- `GET /search?q=<query>` - Search manga by name

### Manga
- `GET /manga/:slug` - Get manga details
- `GET /manga/:slug/chapter/:chapter` - Get chapter with pages

## Development

### Go Workspace

This project uses Go workspaces to manage multiple modules. The workspace includes:
- `./shared` - Common entities and utilities
- `./manga` - Main API service
- `./scraper` - Web scraping service

### Hot Reload

Both services support hot-reload using Air. Configuration files are located at:
- `backend/scraper/.air.toml`
- `backend/manga/.air.toml`

### Adding New Scrapers

New scrapers can be implemented by following the domain-driven design structure in the scraper service:
1. Add scraper logic to `scraper/internals/domain/core/`
2. Register use cases in `scraper/internals/domain/*/use_case/`
3. Create HTTP handlers in `scraper/internals/infra/http/*/handler.go`

## Design Decisions

### Base64 Image Storage
Images are stored as Base64 strings directly in the database to avoid the complexity of blob storage solutions while keeping the system self-contained and portable.

### SQLite over PostgreSQL
Switched from PostgreSQL to SQLite for simplicity and portability, making it easier to run locally without external database dependencies.

### Message Queue Communication
RabbitMQ enables asynchronous communication between the API and scraper services, allowing the API to remain responsive while scraping operations run in the background.

## Roadmap

- [ ] Client application for manga reading
- [ ] Caching layer for frequently accessed content
- [ ] Multi-source scraper support
- [ ] User library management
- [ ] Reading progress tracking
- [ ] Offline reading support

## Contributing

This is a personal project, but suggestions and feedback are welcome.

## License

This project is for personal use and educational purposes.
