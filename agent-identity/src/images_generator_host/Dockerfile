FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/go/src/github.com/6022-labs/agentic-orchestrator/agent-identity

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

WORKDIR /usr/local/go/src/github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator_host

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM chromedp/headless-shell:140.0.7259.2 AS final

# Import the compiled executable from the first stage.
COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/local/go/src/github.com/6022-labs/agentic-orchestrator/agent-identity/src/images_generator_host/images/nft_mark.png /images/nft_mark.png
ENV IMAGE_GENERATION__IMAGE_MARK_PATH=/images/nft_mark.png

EXPOSE 8000

# Run the compiled binary at container start.
ENTRYPOINT ["/app"]