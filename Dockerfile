# step 1: build executable binary
FROM golang:alpine AS builder

RUN mkdir /app
WORKDIR /app

# copy go mod and sum files
COPY go.mod .
COPY go.sum .

RUN go mod download

# copy the source code
COPY . .

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin

RUN chmod 0555 /app/bin

# step 2: build a small image, start from scratch
FROM scratch

WORKDIR /app

# Copy our static executable
COPY --from=builder /app/bin /app/bin

EXPOSE 3000

# migrate and seed: binary must be support command flags
# RUN ./bin --migrate
# RUN ./bin --seed

CMD ["./bin"]
