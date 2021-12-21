################################
# Build Stage
################################
FROM golang:1.16-alpine-alpine@sha256:619019b6ccbd2e49de35ce96d362556fa344e35d54150a186db93804de1e7382 AS builder

ENV GO111MODULE=on

# Copy proto module
WORKDIR /go/src/proto/
COPY ./proto .

# Copy shared module
WORKDIR /go/src/shared/
COPY ./shared .

# Copy clients module
WORKDIR /go/src/clients/
COPY ./clients .

# Copy actual code
WORKDIR /go/src/services/AlignmentService
# Resolve dependencies first for efficient caching
COPY ./services/AlignmentService/go.mod .
RUN go mod download && go mod verify && go mod tidy

# Copy & build source files
COPY ./services/AlignmentService/ .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/picapica-alignment main.go

# Adding the grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

################################
# Deployment Image
################################
FROM gcr.io/distroless/static

# Copy files from builder image
COPY --from=builder /go/bin/picapica-alignment /
COPY --from=builder /bin/grpc_health_probe ./grpc_health_probe

# Set labels accordingly
LABEL version="1.0.0-beta" release-date="2020-08-04" component="alignment"

# Define entry
CMD ["/picapica-alignment serve"]
