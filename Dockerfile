# Start with the official Go image
FROM golang:1.23.2-bookworm

WORKDIR /app

# Install debugging tools and AWS CLI
RUN apt-get update && apt-get install -y \
    curl \
    vim \
    unzip \
    && rm -rf /var/lib/apt/lists/* \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf aws awscliv2.zip

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/api/v1/main.go

EXPOSE 8080

CMD ["./server"]
