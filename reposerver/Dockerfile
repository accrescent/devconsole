FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -tags=debug,nomsgpack -o /reposerver

EXPOSE 8080

CMD [ "/reposerver" ]
