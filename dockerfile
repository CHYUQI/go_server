FROM golang:1.25.0 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLE=0 GOOS=linux go build -o main .



FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD [ "./main" ]


