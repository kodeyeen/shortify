FROM golang:1.24.0 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd/api-server


FROM build AS run-test
RUN go test -v ./...


FROM alpine:latest

WORKDIR /app

COPY --from=build /app/api-server .
COPY --from=build /app/configs configs/

ENTRYPOINT ["/app/api-server"]
