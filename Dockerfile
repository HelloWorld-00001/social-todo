#Alpine: a mini OS - very small Linux distribution,
#often used as a base image in Docker to keep containers lightweight, secure, and efficient.
#FROM alpine
#
#WORKDIR /app/
#ADD ./app /app/
#
#ENTRYPOINT ["./app"]

# -------------------------------
# Multi-stages docker
# -------------------------------
# -------------------------------
# Stage 1: Build the Go binary
# -------------------------------
# Create working directory inside container
# Copy all project files from host to /app/ in container
# Set working directory to /app
#FROM golang:1.23-alpine AS builder
#RUN mkdir /app
FROM social-todo-lib-cached AS builder

ADD . /app/
WORKDIR /app

# Build statically linked Linux binary named demoApp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o social-app ./cmd/server

# -------------------------------
# Stage 2: Create minimal runtime image
# -------------------------------
FROM alpine
WORKDIR /app/

# Copy only the compiled Go binary from builder stage into runtime image
COPY --from=builder /app/social-app .

# Define container entrypoint → always run this binary when container starts
ENTRYPOINT ["/app/social-app"]