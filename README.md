# Voxa Golang Server

A high-performance audio processing API built with Go and Fiber. Features voice filtering, audio-to-video conversion, and Cloudinary integration.

## Features

- **Voice Processing** - Apply various voice filters (robotic, chipmunk, deep, muffled)
- **Audio to Video** - Convert audio files to video with static background image
- **Cloudinary Integration** - Automatic media upload and storage
- **JWT Authentication** - Secure API endpoints with Bearer tokens
- **Swagger Documentation** - Interactive API documentation
- **MongoDB** - Database for user and message storage

## Tech Stack

- **Framework:** [Fiber](https://gofiber.io/) (Express-inspired Go web framework)
- **Database:** MongoDB
- **Media Processing:** FFmpeg
- **Cloud Storage:** Cloudinary
- **Documentation:** Swagger/OpenAPI
- **Authentication:** JWT

## Prerequisites

- Go 1.21+
- FFmpeg installed and in PATH
- MongoDB instance
- Cloudinary account

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/Investorharry19/voxa-golang-server.git
cd voxa-golang-server
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Install FFmpeg

**Windows:**

- Download from [ffmpeg.org](https://ffmpeg.org/download.html)
- Add to system PATH

**macOS:**

```bash
brew install ffmpeg
```

**Ubuntu/Debian:**

```bash
sudo apt update && sudo apt install ffmpeg
```

### 4. Set up environment variables

Create a `.env` file in the root directory:

```env
PORT=3000
MONGODB_URI=mongodb://localhost:27017/voxa
JWT_SECRET=your-super-secret-key

CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
```

### 5. Generate Swagger docs

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go --parseDependency -d ./,./controllers
```

### 6. Run the server

```bash
go run main.go
```

Or use npm scripts:

```bash
npm run dev
```

## API Documentation

Once the server is running, visit:

```
http://localhost:3000/swagger/index.html
```

## API Endpoints

### Authentication

| Method | Endpoint                | Description                 |
| ------ | ----------------------- | --------------------------- |
| POST   | `/account/register`     | Register new user           |
| POST   | `/account/login`        | Login and get JWT token     |
| GET    | `/account/current-user` | Get authenticated user info |

### Audio Processing

| Method | Endpoint         | Description                      |
| ------ | ---------------- | -------------------------------- |
| POST   | `/process-audio` | Apply voice filter to audio      |
| POST   | `/send-audio`    | Process and upload audio message |
| GET    | `/convert`       | Convert audio URL to video       |

## Voice Filters

| Value | Effect                     |
| ----- | -------------------------- |
| `1`   | High-pitched robotic voice |
| `2`   | Thin, chipmunk-like voice  |
| `3`   | Deep, mysterious voice     |
| `4`   | Muffled, distant voice     |

## Usage Examples

### Register User

```bash
curl -X POST http://localhost:3000/account/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "secret123"}'
```

### Login

```bash
curl -X POST http://localhost:3000/account/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "secret123"}'
```

### Process Audio

```bash
curl -X POST http://localhost:3000/process-audio \
  -F "file=@audio.mp3" \
  -F "voice=2" \
  --output processed.mp3
```

### Convert Audio to Video

```bash
curl "http://localhost:3000/convert?audioUrl=https://example.com/audio.mp3" \
  --output video.mp4
```

### Authenticated Request

```bash
curl -X GET http://localhost:3000/account/current-user \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Project Structure

```
voxa-golang-server/
├── main.go                 # Application entry point
├── controllers/            # Route handlers
│   ├── account.go          # Auth controllers
│   ├── audio.go            # Audio processing controllers
│   └── message.go          # Message controllers
├── models/                 # Data models
│   ├── user.go
│   └── message.go
├── routers/                # Route definitions
│   └── routes.go
├── middleware/             # Custom middleware
│   └── auth.go
├── utils/                  # Helper functions
│   ├── response.go
│   ├── hash.go
│   ├── jwt.go
│   └── filters.go
├── database/               # Database connection
│   └── mongo.go
├── docs/                   # Generated Swagger docs
├── .env                    # Environment variables
├── go.mod
└── README.md
```

## Docker

### Build

```bash
docker build -t voxa-server .
```

### Run

```bash
docker run -p 3000:3000 --env-file .env voxa-server
```

### Docker Compose

```yaml
version: "3.8"
services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - MONGODB_URI=mongodb://mongo:27017/voxa
    depends_on:
      - mongo

  mongo:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

volumes:
  mongo_data:
```

```bash
docker-compose up
```

## Deployment

### Using Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:latest
RUN apk add --no-cache ffmpeg
COPY --from=builder /app/server /server
COPY --from=builder /app/.env /.env
CMD ["/server"]
```

### Environment Variables for Production

| Variable                | Description                 |
| ----------------------- | --------------------------- |
| `PORT`                  | Server port (default: 3000) |
| `MONGODB_URI`           | MongoDB connection string   |
| `JWT_SECRET`            | Secret key for JWT signing  |
| `CLOUDINARY_CLOUD_NAME` | Cloudinary cloud name       |
| `CLOUDINARY_API_KEY`    | Cloudinary API key          |
| `CLOUDINARY_API_SECRET` | Cloudinary API secret       |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Harry** - [@Investorharry19](https://github.com/Investorharry19)

## Acknowledgments

- [Fiber](https://gofiber.io/) - Fast Go web framework
- [FFmpeg](https://ffmpeg.org/) - Audio/video processing
- [Cloudinary](https://cloudinary.com/) - Media storage
- [Swagger](https://swagger.io/) - API documentation
