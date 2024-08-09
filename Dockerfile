FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o fantasy-backend cmd/fantasy-backend/main.go

ENTRYPOINT ["./fantasy-backend"]
