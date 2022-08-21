FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -tags=nomsgpack -o /reposerver

EXPOSE 8080

CMD [ "/reposerver" ]
