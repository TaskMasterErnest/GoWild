FROM golang:1.23.5-alpine3.21 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook .

FROM scratch

WORKDIR /app

COPY --from=builder /app/webhook .

EXPOSE 8443

CMD [ "./webhook" ]