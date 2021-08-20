FROM golang:1.17 as builder

# Build the goose binary
RUN CGO_ENABLED=0 GOOS=linux go install github.com/pressly/goose/cmd/goose@v2.7.0

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY hollow /hollow

# Copy goose and database migration files
COPY --from=builder /go/bin/goose /goose
COPY db/migrations /db-migrations

# Run the web service on container startup.
ENTRYPOINT ["/hollow"]
CMD ["serve"]
