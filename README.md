# Map Broker Service

A high-performance map tile caching and proxy service built with Go and Gin framework. This service acts as an intermediary between map clients and map tile providers, providing intelligent caching, theme support, and optimized delivery of map tiles.

## 🚀 Features

- **Map Tile Caching**: Intelligent caching system using Google Cloud Storage
- **Multiple Map Providers**: Support for various map providers (MapTiler used as example)
- **RESTful API**: Clean REST API for map tile requests
- **Automatic Caching**: Tiles are automatically cached on first request
- **Theme Support**: Light and dark map themes
- **High Performance**: Fast tile delivery with caching optimization
- **Cloud Storage Integration**: Uses Google Cloud Storage for persistent caching

## 🏗️ Architecture

The application follows a clean, modular architecture:

```
map_broker/
├── core/map/           # Core domain models and business logic
├── handlers/map/       # HTTP request handlers
├── router/            # Route definitions and middleware
├── service/map/       # Business logic services
├── service/storage_services/  # Storage abstraction layer
└── main.go           # Application entry point
```

### Core Components

- **Map Provider**: Handles different map tile providers (configurable)
- **Storage Service**: Abstracts Google Cloud Storage operations
- **Map Service**: Core business logic for tile retrieval and caching
- **Handler**: HTTP request/response handling
- **Router**: API route definitions

## 📋 Prerequisites

- Go 1.23.0 or higher
- Google Cloud Storage bucket (`map-cached`)
- Map provider API key (MapTiler used as example)
- Google Cloud credentials configured

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd map_broker
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   # Example with MapTiler (you can use any map provider)
   export MAPTILER_API_KEY="your_maptiler_api_key"
   
   # Or use your preferred map provider's API key
   # export YOUR_PROVIDER_API_KEY="your_api_key"
   ```

4. **Configure Google Cloud Storage**
   - Create a bucket named `map-cached`
   - Set up authentication (service account key or default credentials)
   - Place your service account key file in the project root

## 🚀 Running the Application

### Development Mode (with auto-reload)
```bash
./start.sh --env=debug
```

### Manual Run
```bash
go run main.go --port=8080
```

### Production Mode
```bash
go build -o map_broker
./map_broker --port=8080
```

### Deploying to Cloud Services

You can deploy the service to various platforms:

- **Google Cloud Run**
  - Build and push a Docker image
  - Deploy using Cloud Run console or CLI
- **AWS Elastic Beanstalk**
  - Create a Dockerfile
  - Deploy using Elastic Beanstalk CLI
- **Azure App Service**
  - Use Docker or Go runtime
  - Deploy via Azure CLI or portal
- **Other Platforms**
  - Any service supporting Docker or Go binaries

#### Example: Deploy with Docker
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o map_broker

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/map_broker .
CMD ["./map_broker", "--port=8080"]
```

Build and run:
```bash
docker build -t map-broker .
docker run -p 8080:8080 map-broker
```

For cloud deployment, push your image to a registry and follow your provider's deployment steps.

## 📡 API Endpoints

### Get Map Tile

**Endpoint**: `GET /map/{x}/{y}/{z}`

**Parameters**:
- `x` (path): Tile X coordinate
- `y` (path): Tile Y coordinate  
- `z` (path): Zoom level
- `provider` (query): Map provider (e.g., `maptiler`, `openstreetmap`, etc.)
- `theme` (query): Theme mode (`light` or `dark`)

**Example Request**:
```bash
curl "http://localhost:8080/map/10/5/3?provider=maptiler&theme=dark"
```

**Response**: 
- Content-Type: `image/png`
- Body: Map tile image data

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MAPTILER_API_KEY` | MapTiler API key for tile requests (example - you can use any map provider) | Yes* |

### Storage Provider Configuration

By default, the service uses Google Cloud Storage for caching tiles. To use another storage provider (e.g., AWS S3, Azure Blob Storage, local filesystem):

1. **Implement a Storage Client**: Create a new client in `service/storage_services/` (e.g., `s3_service.go`) that implements the same interface as the Google Cloud client (`GetClient`, `GetMapUploadingBucket`, etc.).
2. **Update Uploader/Downloader**: Modify `uploader.go` and `downloader.go` to use your new client based on configuration or environment variable (e.g., `STORAGE_PROVIDER=s3`).
3. **Configure Credentials**: Set up credentials for your provider (e.g., AWS keys, Azure connection string) in environment variables or config files.
4. **Test Integration**: Ensure upload/download logic works with your provider.

Example for AWS S3:
```go
// In service/storage_services/s3_service.go
// Implement GetClient and GetMapUploadingBucket for S3
```

Update your `.env` or environment variables:
```bash
export STORAGE_PROVIDER="s3"
export AWS_ACCESS_KEY_ID="your_key"
export AWS_SECRET_ACCESS_KEY="your_secret"
export AWS_BUCKET_NAME="your_bucket"
```

### Supported Map Providers

The service is designed to work with any map tile provider. Currently configured with:
- **MapTiler** (example implementation): Light and dark themes
  - Light: `streets-v2`
  - Dark: `streets-v2-dark`

**Note**: You can easily configure any map provider by modifying the provider configuration in `core/map/provider.go`. The service supports any tile-based map provider that follows standard tile URL patterns.

## 🏃‍♂️ How It Works

1. **Request Processing**: Client requests a map tile with X, Y, Z coordinates and theme
2. **Cache Check**: Service checks if tile exists in Google Cloud Storage
3. **Cache Hit**: If tile exists, returns cached tile immediately
4. **Cache Miss**: If tile doesn't exist:
   - Downloads tile from configured map provider API
   - Returns tile to client immediately
   - Asynchronously uploads tile to cache for future requests

## 📁 File Structure

```
map_broker/
├── config/                    # Configuration files
├── core/
│   └── map/
│       ├── file.go           # Map response models
│       ├── provider.go       # Map provider definitions
│       └── request.go        # Map request models and logic
├── handlers/
│   └── map/
│       └── handler.go        # HTTP handlers for map endpoints
├── router/
│   ├── router.go            # Main router configuration
│   └── map/
│       └── router.go        # Map-specific routes
├── service/
│   ├── map/
│   │   ├── checker.go       # Cache existence checking
│   │   └── save.go          # Tile saving operations
│   └── storage_services/
│       ├── downloader.go    # File download utilities
│       ├── service.go       # Storage service client
│       ├── uploader.go      # File upload operations
│       └── util.go          # Storage utility functions
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── main.go                  # Application entry point
├── start.sh                 # Development startup script
└── README.md               # This file
```

## 🔍 Key Features Explained

### Intelligent Caching
- Tiles are cached using a structured path: `{provider}/{theme}/{z}/{x}/{y}.png`
- Cache-first approach for optimal performance
- Automatic cache population on first access

### Theme Support
- Light theme: Standard map provider streets (MapTiler example)
- Dark theme: Dark variant of map provider streets (MapTiler example)
- Theme selection via query parameter
- Easily configurable for any map provider

### Error Handling
- Comprehensive error handling for invalid coordinates
- Graceful fallback for provider errors
- Detailed error messages for debugging

## 🧪 Testing

Test the service with different tile coordinates and themes:

```bash
# Light theme tile
curl "http://localhost:8080/map/10/5/3?provider=maptiler&theme=light" -o tile_light.png

# Dark theme tile  
curl "http://localhost:8080/map/10/5/3?provider=maptiler&theme=dark" -o tile_dark.png
```

## 🚀 Performance Considerations

- **Async Caching**: Tile uploads happen asynchronously to avoid blocking responses
- **Cloud Storage**: Leverages Google Cloud Storage for scalable, persistent caching
- **Efficient Routing**: Minimal overhead routing with Gin framework
- **Memory Management**: Streams large files efficiently without loading entirely into memory

## 🔒 Security

- API key management through environment variables
- Google Cloud Storage authentication
- Input validation for coordinates and parameters

## 📝 License

This project is part of the Palestine Jaywalk initiative.

---

## 🇵🇸 Made with Love for Palestine

This project was created with love and solidarity for Palestine. 🇵🇸

*Free Palestine* - We stand with the Palestinian people in their struggle for freedom, justice, and self-determination.

🇵🇸 **Palestine will be free** 🇵🇸

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📞 Support

For issues and questions, please create an issue in the repository or contact the development team.