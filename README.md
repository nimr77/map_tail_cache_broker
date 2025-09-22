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

The Map Broker service uses a flexible storage abstraction layer located in `service/storage_services/` that currently implements Google Cloud Storage but is designed to support multiple storage providers.

#### Current Implementation: Google Cloud Storage

**Default Configuration**:
- **Bucket**: `map-cached` (hardcoded)
- **Authentication**: Service account key file (`map-broker-jaywalk-75c83aba05cf.json`) or default credentials
- **Access Control**: Supports both private and public file access

**Setup Requirements**:
1. Create a Google Cloud Storage bucket named `map-cached`
2. Set up authentication using one of these methods:
   - **Service Account Key**: Place `map-broker-jaywalk-75c83aba05cf.json` in project root
   - **Default Credentials**: Use `gcloud auth application-default login`
3. Grant appropriate permissions to the service account

#### Adding Another Storage Provider

The storage services component is architected to support multiple providers. Here's how to add support for additional storage providers:

**Step 1: Create Provider-Specific Service**
```go
// Create service/storage_services/aws_s3_service.go
package storage_services

import (
    "github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
    client *s3.S3
    bucket string
}

func NewS3Service() (*S3Service, error) {
    // Initialize AWS S3 client
    // Return S3Service instance
}
```

**Step 2: Implement Upload/Download Functions**
```go
func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
    // Implement S3 upload logic
    // Return S3 URL
}

func (s *S3Service) DownloadFile(path string) ([]byte, error) {
    // Implement S3 download logic
    // Return file bytes
}
```

**Step 3: Add Provider Selection Logic**
```go
func GetStorageProvider() string {
    provider := os.Getenv("STORAGE_PROVIDER")
    if provider == "" {
        return "gcs" // Default to Google Cloud Storage
    }
    return provider
}
```

**Step 4: Update Main Functions**
```go
func UploadFile(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
    provider := GetStorageProvider()
    
    switch provider {
    case "gcs":
        return UploadFileToGCS(file, fileHeader, path)
    case "s3":
        return UploadFileToS3(file, fileHeader, path)
    case "azure":
        return UploadFileToAzure(file, fileHeader, path)
    default:
        return "", fmt.Errorf("unsupported storage provider: %s", provider)
    }
}
```

**Step 5: Configure Environment Variables**

For AWS S3:
```bash
export STORAGE_PROVIDER="s3"
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
export AWS_REGION="us-east-1"
export AWS_BUCKET_NAME="your-bucket-name"
```

For Azure Blob Storage:
```bash
export STORAGE_PROVIDER="azure"
export AZURE_STORAGE_ACCOUNT="your_account"
export AZURE_STORAGE_KEY="your_key"
export AZURE_CONTAINER_NAME="your-container"
```

**Step 6: Add Dependencies**
```bash
go get github.com/aws/aws-sdk-go
go get github.com/Azure/azure-storage-blob-go
```

#### Supported Storage Providers

| Provider | Status | Configuration | Documentation |
|----------|--------|---------------|---------------|
| Google Cloud Storage | ✅ Implemented | Default | [Storage Services README](service/storage_services/README.md) |
| AWS S3 | 🔄 Extensible | `STORAGE_PROVIDER=s3` | See implementation guide above |
| Azure Blob Storage | 🔄 Extensible | `STORAGE_PROVIDER=azure` | See implementation guide above |
| Local Filesystem | 🔄 Extensible | `STORAGE_PROVIDER=local` | See implementation guide above |

#### Storage Services Architecture

The storage services component provides:

- **Unified Interface**: Consistent API for all storage operations
- **Multiple Upload Methods**: Support for multipart files and byte arrays
- **Access Control**: Both private and public file access
- **URL Parsing**: Automatic conversion between different URL formats
- **Error Handling**: Comprehensive error handling and validation
- **Extensibility**: Easy to add new storage providers

**Key Files**:
- `service.go`: Core client and bucket management
- `uploader.go`: File upload operations (public/private)
- `downloader.go`: File download and URL parsing
- `util.go`: Utility functions for file operations

For detailed implementation examples and complete documentation, see the [Storage Services README](service/storage_services/README.md).

### Supported Map Providers

The service is designed to work with any map tile provider. Currently configured with:
- **MapTiler** (example implementation): Light and dark themes
  - Light: `streets-v2`
  - Dark: `streets-v2-dark`

**Note**: You can easily configure any map provider by modifying the provider configuration in `core/map/provider.go`. The service supports any tile-based map provider that follows standard tile URL patterns.

## 🏃‍♂️ How It Works

### Map Tile Request Flow

1. **Request Processing**: Client requests a map tile with X, Y, Z coordinates and theme
2. **Cache Check**: Service checks if tile exists in Google Cloud Storage using path `{provider}/{themeMode}/{z}/{x}/{y}.png`
   - Theme mode: `0` = Dark, `1` = Light
3. **Cache Hit**: If tile exists, downloads and returns cached tile immediately
4. **Cache Miss**: If tile doesn't exist:
   - Downloads tile from configured map provider API
   - Returns tile to client immediately
   - Asynchronously uploads tile to cache for future requests

### Storage Operations

The service performs three main storage operations:

- **Check Operation**: Verifies if a tile exists in storage using `bucket.Object(path).Attrs()`
- **Download Operation**: Retrieves cached tiles using `storage_services.DownloadFile(path)`
- **Upload Operation**: Caches new tiles using `storage_services.UploadFileBytes()` (asynchronous)

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
│   └── storage_services/    # Storage abstraction layer
│       ├── downloader.go    # File download and URL parsing
│       ├── service.go       # Core storage client and bucket management
│       ├── uploader.go      # File upload operations (public/private)
│       ├── util.go          # Storage utility functions
│       └── README.md        # Comprehensive storage services documentation
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── main.go                  # Application entry point
├── start.sh                 # Development startup script
└── README.md               # This file
```

## 🔍 Key Features Explained

### Intelligent Caching
- Tiles are cached using a structured path: `{provider}/{themeMode}/{z}/{x}/{y}.png`
- Theme mode values: `0` = Dark theme, `1` = Light theme
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