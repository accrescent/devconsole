FROM golang:alpine

RUN apk --no-cache add libc-dev gcc

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /reposerver

EXPOSE 8080

CMD [ "/reposerver" ]
