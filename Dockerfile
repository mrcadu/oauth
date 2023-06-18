FROM golang:1.20-alpine

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -o ./ ./cmd/oauth/

EXPOSE 8443

CMD [ "./oauth" ]