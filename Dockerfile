FROM golang:1.22-alpine
 
WORKDIR /app

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN go install github.com/cosmtrek/air@latest
 
COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]