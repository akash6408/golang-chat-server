## WEBSOCKET-CHAT
### A WebSocket Chat Application in Go

### Compile and Run

Run with config file path:

```bash
CONFIG_PATH="$(pwd)/config/local.yml" go run ./cmd/server
```

Server listens on the configured address in `config/local.yml` under `http_server.address`.

### WebSocket connect

Use your JWT as a query param:

```text
ws://localhost:8082/ws?token=YOUR_JWT
```

Replace `localhost:8082` if you changed the configured address.