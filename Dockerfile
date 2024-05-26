FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o api ./cmd/main.go

EXPOSE 19090

CMD ["./api"]