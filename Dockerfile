FROM golang:alpine
WORKDIR /code
# ENV TG_BOT_TOKEN=""
COPY . /code/
RUN go build -o api