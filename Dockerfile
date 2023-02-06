FROM golang:alpine

WORKDIR /usr/src/gitlab-watcher

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/gitlab-watcher .
WORKDIR /opt/gitlab-watcher

CMD ["gitlab-watcher"]