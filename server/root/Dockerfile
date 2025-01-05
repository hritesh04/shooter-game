FROM golang AS builder

WORKDIR app

COPY go.mod go.sum ./

RUN go mod download

COPY root.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin ./root.go

FROM scratch

COPY --from=builder /bin/root /bin/root

EXPOSE 3000

CMD ["/bin/root"]