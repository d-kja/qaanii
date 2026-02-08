# Manga Scraper - Architecture Improvements & Design Issues

## Database Schema Not Implemented

**Location:** `manga/internals/infra/database/schemas.sql`

**Current Issue:**
The file is empty - no persistence layer exists.

**Problems:**

**Recommended Solution:**
- [ ] Design proper schema for mangas, chapters, pages, users
- [ ] Implement repository pattern for data access
- [?] Add SQLite migrations with versioning (golang-migrate)
- [ ] Use Redis for caching frequently accessed data
- [ ] Implement cache invalidation strategy

---

## Redis Configured But Unused

**Location:** `docker-compose.yml:14-22`

**Current Issue:**
Redis is provisioned but never used in the codebase.

```yaml
manga-cache:
  image: redis:7.4-alpine
  ports:
    - 6379:6379
```

**Recommended Solution:**
- [ ] Implement Redis client in shared module
- [ ] Cache search results with TTL (e.g., 1 hour)
- [ ] Cache manga metadata (longer TTL)
- [ ] Store idempotency keys to prevent duplicate scraping

---

## Hardcoded Timeouts & Magic Numbers

**Current Issue:**

```go
// search.grpc.go:125
case <-time.After(60 * time.Second):

// scraper.entity.go:60
time.Sleep(time.Second * 5)

// channels.constant.go:54
ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
```

**Recommended Solution:**
- [ ] Create configuration struct with defaults
- [?] Make timeouts configurable per-operation

---

## No Retry Logic for Scraping Operations

**Location:** `scraper/internals/domain/search/use_case/search_by_name.usecase.go`

**Recommended Solution:**
- [ ] Implement retry with exponential backoff
- [?] Add circuit breaker pattern (e.g., gobreaker)

---

## No Graceful Shutdown

**Location:** `manga/cmd/main.go`, `scraper/cmd/main.go`

**Current Issue:**
No signal handling for graceful shutdown.

**Recommended Solution:**
- [ ] Add signal handlers (SIGTERM, SIGINT)

---

## Logging Inconsistency

**Current Issue:**
Mix of `log.Printf` and no structured logging.

```go
log.Printf("[BROKER/SUBSCRIBER] - Manga queue creation failed, error: %+v\n", err)
log.Printf("Creating new page, base url [%v]\n", url)
```

**Recommended Solution:**
- [?] Add request/correlation IDs
- [ ] Add log levels based on severity
- [ ] Include context (service, operation, duration)

---

## Missing Reply Queue in Scrape Events

**Location:** `scraper/internals/infra/broker/manga/scrape-manga.event.go:49-54`

**Current Issue:**

```go
pub_message := events.ScrapedMangaMessage{
    BaseEvent: events.BaseEvent{},  // Empty metadata - no reply queue!
    Data:      response.Manga,
}
pub_message.GenerateEventId(string(events.SCRAPED_MANGA_EVENT), "n/a")
```

**Recommended Solution:**
- [ ] Forward `Metadata.Reply` from incoming message
- [ ] Ensure all event handlers preserve correlation data
- [ ] Add validation for required reply queue in request-reply patterns

---

## No Input Validation/Sanitization

**Location:** `scraper/internals/infra/http/search/search-by-name.handler.go`

**Current Issue:**
Limited validation, no sanitization.

**Recommended Solution:**
- [ ] Add comprehensive input validation
- [ ] Sanitize scraped content before storage/return
- [ ] Use parameterized queries for database
- [ ] Validate URL patterns and slugs

---

## No Health Checks / Readiness Probes

**Current Issue:**
No health endpoints exist.

**Recommended Solution:**
- [ ] Add `/health` endpoint checking:
  - [ ] RabbitMQ connection status
  - [ ] Redis connection status (when implemented)
  - [ ] Database connection status
- [ ] Add `/ready` for readiness probes
- [ ] Implement circuit breakers with health reporting

---

## Stub Implementations in Production Code

**Location:** `manga/internals/infra/grpc/manga.grpc.go`, `chapter.grpc.go`

**Recommended Solution:**
- [ ] Implement full request-reply pattern (like SearchService)

---

## Single Source Coupling

**Current Issue:**
Hardcoded to MangaBuddy.com only.

```go
BASE_URL=https://mangabuddy.com
```

**Recommended Solution:**
- [ ] Create source interface/adapter pattern
- [ ] Implement source registry
- [ ] Abstract XPath selectors per source
- [ ] Add source health monitoring
- [ ] Enable fallback to alternative sources

---

## No Rate Limiting

**Current Issue:**
No protection against abuse or overwhelming sources.

**Recommended Solution:**
- [ ] Add client rate limiting

---

## Missing Observability

**Current Issue:**
No metrics, tracing, or proper monitoring.

**Recommended Solution:**
- [ ] Implement OpenTelemetry tracing
- [ ] Add distributed tracing across services
- [ ] Create Grafana dashboards

---

## Inefficient Protobuf Conversion

**Location:** `manga/internals/infra/grpc/search.grpc.go:105-109`

**Current Issue:**

```go
mangas := []*mangav1.Manga{}
for _, manga := range message.Data {
    buf_manga := manga.ToProtobuf()
    mangas = append(mangas, &buf_manga)
}
```

**Problems:**
- [ ] Slice grows dynamically (reallocations)
- [ ] No pre-allocation despite known size
- [ ] Conversion happens on every request

**Recommended Solution:**
- [ ] Pre-allocate slice: `make([]*mangav1.Manga, 0, len(message.Data))`
- [ ] Consider caching converted results
- [ ] Use code generation for boilerplate conversions

---

## Priority Matrix

| Priority | Issue | Impact | Effort |
|----------|-------|--------|--------|
| Critical | Message ACK on failure | Data loss | Medium |
| Critical | Missing reply queue | Broken functionality | Low |
| High | Database schema | No persistence | High |
| High | Redis unused | Wasted resources | Medium |
| High | No graceful shutdown | Data loss on deploy | Medium |
| Medium | No rate limiting | Abuse risk | Medium |
| Low | Missing observability | Operability | Medium |
