FROM golang:1.20-bullseye as build-stage

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build -v -o server cmd/main.go 

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /rinha-backend

COPY --from=build-stage /app/server ./server
COPY --from=build-stage /app/migrations ./migrations

EXPOSE 8000

CMD ["/rinha-backend/server"]
