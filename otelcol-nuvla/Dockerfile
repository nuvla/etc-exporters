FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o otelcol-nuvla .

EXPOSE 4317
EXPOSE 4318

ENTRYPOINT ["./otelcol-nuvla"]