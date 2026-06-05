FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download || true
COPY . .
RUN go build -o /bin/core-backoffice .

FROM alpine:3.22
COPY --from=build /bin/core-backoffice /bin/core-backoffice
EXPOSE 3005
CMD ["/bin/core-backoffice"]
