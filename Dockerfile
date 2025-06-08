FROM golang:1.24

WORKDIR /app
COPY . .

# ✅ COPY key ลง /app/config/
COPY ./config/serviceAccountKey.json /app/config/serviceAccountKey.json

RUN go mod tidy
CMD ["go", "run", "cmd/main.go"]
