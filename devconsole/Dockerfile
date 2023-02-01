FROM golang:1.20-alpine

RUN apk --no-cache add libc-dev gcc protoc protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download
RUN go install github.com/mattn/go-sqlite3

COPY . .

RUN go generate && go build -tags=debug,nomsgpack -o /devconsole

EXPOSE 8080

CMD [ "/devconsole" ]
