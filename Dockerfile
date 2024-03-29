FROM golang:alpine

WORKDIR /app

COPY ./ast /app/ast
COPY ./engine /app/engine
COPY ./lexer /app/lexer
COPY ./parser /app/parser
COPY ./token /app/token
COPY ./modes /app/modes

COPY go.mod /app
COPY main.go /app

RUN go build -o go4sql-docker

ENTRYPOINT ["./go4sql-docker"]
