FROM golang:1.18
RUN mkdir /app
ADD . /app
WORKDIR /app
ENV PORT=8000
RUN go build -o main .
EXPOSE $PORT
CMD ["/app/main"]