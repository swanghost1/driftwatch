# driftwatch

> CLI tool that detects config drift between deployed services and their declared state

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

```bash
# Check drift for all services defined in your config
driftwatch check --config ./services.yaml

# Check a specific service
driftwatch check --service api-gateway --config ./services.yaml

# Output results as JSON
driftwatch check --config ./services.yaml --output json
```

### Example Output

```
[DRIFT]   api-gateway     replicas: declared=3, deployed=2
[DRIFT]   auth-service    image tag: declared=v1.4.2, deployed=v1.4.0
[OK]      billing-service no drift detected
```

---

## Configuration

Define your expected service state in a `services.yaml` file:

```yaml
services:
  - name: api-gateway
    replicas: 3
    image: api-gateway:v2.1.0
  - name: auth-service
    replicas: 2
    image: auth-service:v1.4.2
```

---

## License

MIT © 2024 yourusername