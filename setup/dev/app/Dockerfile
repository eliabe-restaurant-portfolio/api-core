FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

RUN go build -o executable cmd/http/main.go

EXPOSE 8080

ENTRYPOINT ["./executable"]
