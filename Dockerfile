# Build
# FIXME: multistage
FROM golang AS build

RUN ls
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN ls
RUN go build -v -o /usr/local/bin/main cmd/lb_go/main.go

CMD ["main", "start"]