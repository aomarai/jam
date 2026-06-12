# ==========================================
# STAGE 1: Build the Go Binary
# ==========================================
FROM golang:1.26.4-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bot ./cmd/bot

# ==========================================
# STAGE 2: Final Lightweight Runner
# ==========================================
FROM scratch

WORKDIR /app

# copy binary
COPY --from=builder /app/bot .

# create data dir for model persistence
COPY --from=builder /app/data /app/data

ENTRYPOINT ["./bot"]