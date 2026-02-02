FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/go/src/github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

WORKDIR /usr/local/go/src/github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_host

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

# Import the compiled executable from the first stage.
COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

COPY --from=builder /usr/local/go/src/github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/migrations/mysql/*.sql /migrations/database_connector_db_agent/mysql/
COPY --from=builder /usr/local/go/src/github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/migrations/mariadb/*.sql /migrations/database_connector_db_agent/mariadb/
COPY --from=builder /usr/local/go/src/github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_agent/migrations/postgres/*.sql /migrations/database_connector_db_agent/postgres/
ENV AGENT_MIGRATIONS__FOLDER_PATH=/migrations/database_connector_db_agent

COPY --from=builder /usr/local/go/src/github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector_db_orchestrator/migrations/*.sql /migrations/database_connector_db_orchestrator/
ENV ORCHESTRATOR_MIGRATIONS__FOLDER_PATH=/migrations/database_connector_db_orchestrator

EXPOSE 8000

# Run the compiled binary at container start.
ENTRYPOINT ["/app"]