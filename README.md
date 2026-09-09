# Courtsite URL Shortener Service (v2 Architecture)

A high-performance, modular URL shortening and analytics service written in idiomatic Go. Refactored from a v1 monolith into a production-grade, package-by-feature architecture using Go's standard library (`net/http`) with Go 1.22+ routing, dependency injection, and thread-safe in-memory storage.

---

## Technical Features Implemented

* **`POST /shorten`**: Validates input syntax and performs $O(1)$ URL deduplication. If a URL is submitted for the first time, a 6-character Base64 key is generated and stored with an initial visit count of `1`. Subsequent submissions of the same URL reuse the existing short key and increment its visit count.
* **`GET /analytics`**: Accepts a JSON payload containing a `short_key` and returns the target URL along with its total visit metrics.
* **Zero External Framework Overhead**: Built purely on Go's standard `net/http` package, taking advantage of Go 1.22+ enhanced `http.ServeMux` method matching (`POST /shorten`, `GET /analytics`).
* **Thread-Safe In-Memory Storage**: Employs `sync.RWMutex` with dual-map indexing (`record` map by `ShortKey` and `urlIndex` map by `OriginalURL`) to safely handle high-concurrency read and write operations without data races.
* **Automated Unit & Race Testing**: Comprehensive table-driven unit tests covering domain operations, edge cases, and concurrency guarantees using `go test -race`.

---

## Project Structure

```text
courtsite-url-shortener/
├── cmd/
│   └── api/
│       └── main.go           # Application entrypoint & dependency injection wiring
├── internal/
│   └── shortener/
│       ├── handler.go        # HTTP transport layer (JSON codecs & endpoint handlers)
│       ├── memory.go         # Thread-safe in-memory store (RWMutex + URL index)
│       ├── model.go          # Domain entities and error definitions
│       ├── repository.go     # Data storage interface contract
│       ├── service.go        # Business logic and key generation
│       └── service_test.go   # Table-driven unit tests
├── Dockerfile                # Optimized multi-stage Docker build
├── docker-compose.yml        # Orchestration configuration
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── README.md                 # Complete project documentation

```

---

## Prerequisites & Getting Started

### 1. Run via Docker (Recommended)

Ensure Docker Desktop is running, then execute:

```bash
docker-compose up --build

```

The application will start on `http://localhost:8080`.

### 2. Run Locally (Go 1.22+ Installed)

```bash
go run cmd/api/main.go

```

### 3. Run Automated Tests with Race Detector

```bash
go test -v -race ./...

```

---

## API Usage Examples

### 1. Shorten a URL (or Increment Visit Count)

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "www.google.com"}'

```

**Response (`202 Accepted`):**

```json
{
  "ShortKey": "aB3x9k",
  "OriginalURL": "www.google.com",
  "Visits": 1,
  "CreatedAt": "2026-09-10T02:25:00Z"
}

```

*Note: Posting the same URL again will return the same `ShortKey` with `Visits` incremented.*

### 2. Retrieve Analytics

```bash
curl -X GET http://localhost:8080/analytics \
  -H "Content-Type: application/json" \
  -d '{"short_key": "aB3x9k"}'

```

**Response (`202 Accepted`):**

```json
{
  "OriginalUrl": "www.google.com",
  "Visits": 1
}

```

---

## Technical Evaluations & Architectural Decisions

### Assumptions & Trade-offs

* **Standard Library (`net/http`) vs. Frameworks (Gin):** Transitioned away from external web frameworks to eliminate third-party supply chain dependencies, optimize binary size, and demonstrate low-level HTTP multiplexing and context handling using Go 1.22+ `ServeMux`.
* **Atomic URL Deduplication:** Implemented a secondary `urlIndex` map inside the memory store. This allows constant time $O(1)$ lookups to detect existing URLs and increment visit counters on incoming `POST /shorten` requests under a single `sync.RWMutex` write lock.
* **In-Memory Concurrency:** Uses `sync.RWMutex` to separate read locks (`RLock` for fetching analytics) from write locks (`Lock` for saving new keys or incrementing visit counts), maximizing throughput under high read volumes.

### Production Scaling Architecture

To scale this service to production for high-throughput traffic:

* **Persistence & Distributed Cache Layer:** Replace the `Memory` store with **Redis** for sub-millisecond short-key and URL index caching, backed by **PostgreSQL** or **DynamoDB** for durable persistence. V3
* **Asynchronous Analytics Pipeline:** Offload click processing and visit analytics updates to an event stream (**AWS SQS**, **RabbitMQ**, or **Kafka**) to prevent database write bottlenecks during peak traffic spikes.