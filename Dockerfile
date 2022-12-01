FROM gcr.io/distroless/static

# Copy the binary that goreleaser built
COPY serverservice /serverservice

# Run the web service on container startup.
ENTRYPOINT ["/serverservice"]
CMD ["serve"]
