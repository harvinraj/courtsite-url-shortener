# Courtsite URL Shortener Service

A high-performance prototype URL shortening and analytics proxy web service built in Go (Gin Framework).

## Technical Requirements Implemented
- **POST `/shorten_url/`**: Validates URL format and accessibility via HTTP HEAD/GET probes before returning a generated 7-character base62 key.
- **GET `/shorten_url/<key>`**: Redirects to the original URL via HTTP 302. Returns 404 if key does not exist.
- **GET `/analytics/<key>`**: Returns total visit counts and last accessed timestamp for physical location tracking (Business Intelligence).
- **Automated Tests**: Unit and HTTP integration tests included (`main_test.go`).

---

## Project Structure

```text
courtsite-url-shortener/
├── main.go            # Core service, router, handlers, & storage
├── main_test.go       # Automated unit & HTTP integration tests
├── Dockerfile         # Multi-stage container build
├── docker-compose.yml # Container orchestration
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
└── README.md          # Complete project documentation

```

---

## Technical Requirements Implemented

* **POST `/shorten_url/**`: Accepts `{ "url": "https://..." }`. Validates syntax format and executes HTTP `HEAD`/`GET` reachability checks before returning a unique 7-character base62 key.


* **GET `/shorten_url/<key>**`: Redirects visitors to the original destination URL via `HTTP 302 Found`. Returns `HTTP 404 Not Found` if the key does not exist.


* **GET `/analytics/<key>` (Stretch Goal)**: Exposes referral metrics including total visit counts and timestamp of the last access for Business Intelligence.


* **Automated Test Suite**: Integration tests verifying key generation, URL validation failure, 302 redirects, and 404 edge cases.



---

## Prerequisites & Getting Started

### 1. Run via Docker (Recommended)

Make sure Docker Desktop is running, then execute:

```bash
docker-compose up --build

```

The application will start on `http://localhost:8080`.

### 2. Run Locally (Go Installed)

```bash
go run main.go

```

### 3. Run Automated Tests

```bash
go test -v ./...

```

---

## API Usage Examples

### 1. Shorten a URL

```bash
curl -X POST http://localhost:8080/shorten_url/ \
  -H "Content-Type: application/json" \
  -d '{"url": "[https://www.google.com](https://www.google.com)"}'

```

**Response (`200 OK`):**

```json
{
  "key": "aB3x9kL"
}

```

### 2. Redirect to Original URL

```bash
curl -i http://localhost:8080/shorten_url/aB3x9kL

```

**Response (`302 Found`):**

```text
HTTP/1.1 302 Found
Location: [https://www.google.com](https://www.google.com)

```

### 3. Retrieve Analytics (Business Intelligence)

```bash
curl http://localhost:8080/analytics/aB3x9kL

```

**Response (`200 OK`):**

```json
{
  "original_url": "[https://www.google.com](https://www.google.com)",
  "short_key": "aB3x9kL",
  "created_at": "2026-08-26T16:15:00Z",
  "click_count": 1,
  "last_visited": "2026-08-26T16:16:10Z"
}

```

---

## Technical Evaluations & Architectural Decisions

### Assumptions & Trade-offs

* **HTTP 302 vs 301 Redirection:** We specifically selected `302 Found` (temporary redirect) instead of `301 Moved Permanently`. A `301` status tells browsers to cache the destination locally, bypassing our shortener service on subsequent scans and blinding our Business Intelligence analytics collection.


* **URL Reachability Checks:** Reachability validation includes a 3-second HTTP timeout with custom `User-Agent` headers and a fallback from `HEAD` to `GET`. This prevents malicious or dead URLs from being registered while ensuring slow target servers don't hang client requests.


* **In-Memory Concurrency:** Utilized Go's built-in `sync.RWMutex` around internal map access to protect against data races and panic recovery under concurrent client hits.


### Production Scaling Architecture (Millions of Requests)

To scale this service to production for high-throughput social campaigns and national physical QR deployment:

* **Caching & Persistence Layer:** Replace in-memory maps with **Redis** for sub-millisecond key-to-URL lookup caching, backed by **PostgreSQL** for persistent record storage.


* **Asynchronous Analytics Pipeline:** Offload click processing from the HTTP request path. Push scan metadata (`key`, `user-agent`, `ip_address`, `timestamp`) to an **AWS SQS** or **Kafka** queue. A dedicated consumer worker updates analytics in real-time without increasing redirect latency for users.


* **Stateless Container Deployment:** Deploy stateless Go binary containers behind an AWS Application Load Balancer (ALB) scaled dynamically via Kubernetes Horizontal Pod Autoscaling (HPA).
