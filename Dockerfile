FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o ai-proxy

CMD ["./ai-proxy"]
#也是一个简单的示例

