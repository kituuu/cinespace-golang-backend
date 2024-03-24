FROM golang:1.22.1-alpine3.18

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 10000:10000

CMD [ "go", "run", "main.go" ]