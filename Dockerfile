# 
FROM golang AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o ./lb_go cmd/lb_go/main.go

#
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/lb_go ./

RUN chmod +x ./lb_go
ENTRYPOINT ["./lb_go", "start"]