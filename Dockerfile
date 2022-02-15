FROM golang:1.17 as builder

WORKDIR /app

COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 go build -o lux ./cmd

FROM nginx:1.21-alpine

COPY --from=builder /app/lux /bin/lux

COPY ./example /projects

COPY ./templates /templates

CMD ["/bin/lux", "-t", "/templates", "-p", "/projects", "-o", "/config", "--exec"]
