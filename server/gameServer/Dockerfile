FROM golang AS builder

WORKDIR app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin ./main.go

FROM scratch

COPY --from=builder /bin/main /bin/main

EXPOSE 3000

CMD ["/bin/main"]