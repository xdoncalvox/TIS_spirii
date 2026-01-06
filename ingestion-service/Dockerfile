FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ingestion ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/ingestion /ingestion
ENV KAFKA_BROKERS=localhost:9092
ENV KAFKA_TOPIC=charger.telemetry.raw

EXPOSE 8080
ENTRYPOINT ["/ingestion"]
