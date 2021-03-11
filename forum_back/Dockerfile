FROM golang:latest

WORKDIR /usr/src/app/
COPY . /usr/src/app/
RUN go mod download
RUN go build -o forum-back .
EXPOSE 4333

ENV TZ Asia/Almaty

CMD ["./forum_back"]

# docker build -t forum-back .
# docker volume create forum-volume
# sudo docker run --rm --name forum-back -p 4333:4333 -v forum-volume:/usr/src/app/ forum-back