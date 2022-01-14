FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /devportal cmd/devportal/main.go

EXPOSE 8080

CMD [ "/devportal" ]
