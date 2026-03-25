# Go Server

A simple Go server project.

## Features

- Lightweight and fast
- Easy to configure
- RESTful API support

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher

### Installation

```bash
git clone https://github.com/yourusername/go_server.git
cd go_server
go mod tidy
```

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## Project Structure

```
go_server/
├── dockerfile
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── testing.go
├── README.md
├── LICENSE
```

## Docker

This project includes a multi-stage Dockerfile (`dockerfile`) to build a minimal container:

- builder stage uses `golang:1.25.0` to download dependencies and compile the app.
- runtime stage uses `alpine:3.19` and copies the static binary from the builder.
- Exposes port `8080` and runs `./main`.

Build and run:

```bash
docker build -t go_server:latest -f dockerfile .
docker run -p 8080:8080 go_server:latest
```

## Tests

The project includes unit tests in `main_test.go`:

- `TestHelloHandler` validates hello endpoint behavior for default and error inputs.
- `TestHelloHandler_Metrics` asserts Prometheus metrics (`httpRequestsTotal`, `httpRequestDuration`) update correctly.
- A lightweight `TestEmpty` placeholder exists.

Run tests:

```bash
go test ./...
```

## Deploy to Alibaba ECS

This section shows a simple way to run this project on Alibaba Cloud Elastic Compute Service (ECS).

1. Prepare Alibaba ECS instance

- Create an ECS instance (Linux).
- Open security group inbound rules for TCP 22 (SSH) and 8080 (app).

2. Install Docker on ECS

```bash
sudo yum update -y
sudo yum install -y docker
sudo systemctl enable docker --now
```

3. Clone repo and build image on ECS

```bash
cd /home/ec2-user
git clone https://github.com/CHYUQI/go_server.git
cd go_server
docker build -t go_server:latest -f dockerfile .
```

4. Run container on ECS

```bash
docker run -d --name go_server -p 8080:8080 go_server:latest
```

5. Verify

- Access `http://<ECS_PUBLIC_IP>:8080/api/hello` from browser or curl.

### Optional: deploy from local via Docker Hub

1. Build and push to registry:

```bash
docker build -t yourdockerid/go_server:latest -f dockerfile .
docker push yourdockerid/go_server:latest
```

2. Pull and run on ECS:

```bash
docker pull yourdockerid/go_server:latest
docker run -d --name go_server -p 8080:8080 yourdockerid/go_server:latest
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

[MIT](LICENSE)