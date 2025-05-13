FROM debian:bullseye

# Install required dependencies
RUN apt-get update && apt-get install -y wget tar git build-essential && rm -rf /var/lib/apt/lists/*

# Download and install Go 1.22.5
RUN wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm go1.22.5.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Verify Go version
RUN go version

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /backend .

EXPOSE 4000
CMD ["/backend"]
