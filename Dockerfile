FROM oven/bun AS web-build
WORKDIR /src/web
COPY web/package*.json ./
RUN bun install
COPY web ./
RUN bun run build

FROM golang:1.20 AS builder
WORKDIR /src

# copy go modules and source
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
# Copy built web assets into expected location for go:embed
COPY --from=web-build /src/web/dist ./frontend/dist

# Build the Go binary (frontend must be built before this step in CI/local)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /out/cups-web ./cmd/server

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/cups-web /cups-web
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/cups-web"]
