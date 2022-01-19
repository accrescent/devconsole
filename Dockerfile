ARG GO_VERSION=1.17
ARG ALPINE_VERSION=3.15

# Build stage
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY cmd cmd

COPY server server

RUN CGO_ENABLED=0 go build \
    -ldflags="-extldflags=-static" \
    -o /devportal cmd/devportal/main.go


# Final stage
FROM gcr.io/distroless/static as final

USER nonroot:nonroot

WORKDIR /app

COPY --from=build --chown=nonroot:nonroot /app /app

COPY --from=build --chown=nonroot:nonroot /devportal /

COPY web web

EXPOSE 8080

ENTRYPOINT [ "/devportal" ]
