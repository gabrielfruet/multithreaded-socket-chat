FROM golang:1.22

WORKDIR /app

COPY ./server ./server

RUN cd ./server ; go build -o main .

EXPOSE 5555

CMD ["./server/main"]


