FROM golang:1.21 AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN go build -o /app/79a7fa75-b820 .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/79a7fa75-b820 /app/79a7fa75-b820

EXPOSE 3040

CMD [ "./79a7fa75-b820" ]