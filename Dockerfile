FROM alpine:3 as alpine
RUN apk add --no-cache ca-certificates

FROM scratch
# Copy ca-certs from alpine
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary that goreleaser built
COPY serverservice /serverservice

# Run the web service on container startup.
ENTRYPOINT ["/serverservice"]
CMD ["serve"]
