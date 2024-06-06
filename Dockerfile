#################
# Builder image
#################
FROM golang:1.22-bullseye AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build main.go

#################
# Final image
#################
FROM gcr.io/distroless/base

COPY --from=builder /app/main /

CMD ["/main"]
