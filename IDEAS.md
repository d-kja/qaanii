### Async approach

Scrapers can take time to process the request, so I decided to approach 
it using a message queue to send the requests asynchronously to a separate 
service. Here's an example of a search:

```
┌─────────────┐ 
│   Client    │ 
│ Application │ 
└─────────────┘ 
   │
   │ GET /search?q=konosuba - Send request with x-idempotency-key, e.g.: [session-id]-[feature]-[request-id] -> aw32k%03ad123-search-konosuba
   │ 
   ▼
┌──────────────┐   Stores the status & key    ┌─────────────┐
│  Manga API   │ ◄──────────────────────────► │    REDIS    │
└──────────────┘                              └─────────────┘
   │
   │ Sends message through broker
   │         (RabbitMQ)
   │
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

This is a simple approach leaning more on a message architecture, but it
helps with the nature of the Scraper and with one of the features I'll be 
adding in the future, download many at once (heavy burden on the API if it was coupled).

---

### Why gRPC?

It was more of a learning choice, but it does help with the server streaming. I was 
planning on using pooling for my requests, and abuse of the REDIS cache updating the 
request status once the consumer was updated. However, I learned that gRPC has a server 
side streaming feature that does about the same, but with better performance and less 
burden on the server (multiple requests).
