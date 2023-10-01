FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build -v -o server cmd/rinha-gin/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /rinha-backend

COPY --from=builder /app/server ./server
COPY --from=builder /app/migrations ./migrations

CMD ["/rinha-backend/server"]
