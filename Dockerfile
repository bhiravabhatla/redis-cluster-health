FROM golang as build
COPY . /app
WORKDIR /app
RUN go mod tidy && \
    go build -o redis-cluster-health

FROM golang
COPY --from=build  /app/redis-cluster-health /redis-cluster-health
