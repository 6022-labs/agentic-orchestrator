FROM golang:1.27-alpine AS builder

WORKDIR /usr/local/go/src/github.com/6022-labs/agentic-orchestrator/agent-identity

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

WORKDIR /usr/local/go/src/github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

# Import the compiled executable from the first stage.
COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/local/go/src/github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/migrations/*.sql /migrations/
ENV MIGRATIONS__FOLDER_PATH=/migrations

EXPOSE 8000

# Run the compiled binary at container start.
ENTRYPOINT ["/app"]