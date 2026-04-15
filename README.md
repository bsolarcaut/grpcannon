# grpcannon

A lightweight load-testing CLI for gRPC services with configurable concurrency and latency histograms.

---

## Installation

```bash
go install github.com/yourusername/grpcannon@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/grpcannon.git && cd grpcannon && go build -o grpcannon .
```

---

## Usage

```bash
grpcannon [options] <host:port>
```

### Example

```bash
grpcannon \
  --proto ./api/service.proto \
  --call helloworld.Greeter/SayHello \
  --data '{"name": "world"}' \
  --concurrency 50 \
  --requests 1000 \
  localhost:50051
```

### Common Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--proto` | | Path to the `.proto` file |
| `--call` | | Fully qualified method name |
| `--data` | | JSON request payload |
| `--concurrency` | `10` | Number of concurrent workers |
| `--requests` | `200` | Total number of requests |
| `--timeout` | `20s` | Per-request timeout |
| `--insecure` | `false` | Skip TLS verification |

### Sample Output

```
Summary:
  Total requests:   1000
  Duration:         4.32s
  Requests/sec:     231.48
  Errors:           0

Latency histogram:
  p50:   18ms
  p90:   42ms
  p95:   67ms
  p99:   113ms
  max:   204ms
```

---

## License

MIT © 2024 grpcannon contributors