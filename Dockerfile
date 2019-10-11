FROM golang:latest

WORKDIR /charly-weather

COPY . .

RUN go mod download
RUN go test $(go list ./...)
RUN go build -o charly-weather

CMD ["./charly-weather"]