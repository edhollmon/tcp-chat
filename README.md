![tcp-chat](assets/tcp-header.png)

# tcp-chat

A simple TCP chat server and client written in Go. Multiple clients connect to the server and messages are broadcast to all other connected clients.

## Prerequisites

- [Go 1.25+](https://go.dev/dl/)

## Getting started

```bash
git clone https://github.com/edhollmon/tcp-chat.git
cd tcp-chat
```

Build the server and client:

```bash
go build -o bin/server ./cmd/server
go build -o bin/client ./cmd/client
```

Run the server:

```bash
./bin/server
# Starting Simple TCP Server
# Server is ready
```

Run a client (in a separate terminal):

```bash
./bin/client
# Connected to :3000
```

Type a message and press Enter to send. Press `Ctrl+D` to disconnect.

## Flags

### Server

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:3000` | Address to listen on |
| `-trace` | `false` | Enable runtime tracing to `trace.out` |

```bash
./bin/server -addr :8080
./bin/server -trace && go tool trace trace.out
```

### Client

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:3000` | Server address to connect to |

```bash
./bin/client -addr :8080
```

## Connecting with other tools

You can also connect using standard CLI tools:

```bash
nc localhost 3000
```

```bash
telnet localhost 3000
```

## Project structure

```
.
├── cmd/
│   ├── server/main.go                    # Server entry point
│   └── client/main.go                    # Client entry point
└── internal/
    ├── server/simple-server.go           # TCP listener, connection dispatch, and broadcast
    └── client/simple-client.go           # TCP client — connect, read, and write loops
```
