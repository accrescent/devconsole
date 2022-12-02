FROM golang:1.19-alpine

RUN apk --no-cache add libc-dev gcc

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go install github.com/mattn/go-sqlite3

COPY . .

RUN go build -tags=nomsgpack -o /devconsole

EXPOSE 8080

CMD [ "/devconsole" ]
