# docker build . -t wreulicke/probes-tests
FROM golang:1.22.3

WORKDIR /app
ADD ./ /app
RUN go build -o app .

CMD ["/app/app"]
