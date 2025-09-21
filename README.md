# Map Broker Service

A high-performance map tile caching and proxy service built with Go and Gin framework. This service acts as an intermediary between map clients and map tile providers, providing intelligent caching, theme support, and optimized delivery of map tiles.

## 🚀 Features

- **Map Tile Caching**: Intelligent caching system using Google Cloud Storage
- **Multiple Map Providers**: Support for MapTiler with light and dark themes
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

- **Map Provider**: Handles different map tile providers (MapTiler)
- **Storage Service**: Abstracts Google Cloud Storage operations
- **Map Service**: Core business logic for tile retrieval and caching
- **Handler**: HTTP request/response handling
- **Router**: API route definitions

## 📋 Prerequisites

- Go 1.23.0 or higher
- Google Cloud Storage bucket (`map-cached`)
- MapTiler API key
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
   export MAPTILER_API_KEY="your_maptiler_api_key"
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

## 📡 API Endpoints

### Get Map Tile

**Endpoint**: `GET /map/{x}/{y}/{z}`

**Parameters**:
- `x` (path): Tile X coordinate
- `y` (path): Tile Y coordinate  
- `z` (path): Zoom level
- `provider` (query): Map provider (`maptiler`)
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
| `MAPTILER_API_KEY` | MapTiler API key for tile requests | Yes |

### Google Cloud Storage

The service uses Google Cloud Storage for caching tiles. Ensure you have:

1. A bucket named `map-cached`
2. Proper authentication configured
3. Service account with Storage Admin permissions

### Supported Map Providers

Currently supports:
- **MapTiler**: Light and dark themes
  - Light: `streets-v2`
  - Dark: `streets-v2-dark`

## 🏃‍♂️ How It Works

1. **Request Processing**: Client requests a map tile with X, Y, Z coordinates and theme
2. **Cache Check**: Service checks if tile exists in Google Cloud Storage
3. **Cache Hit**: If tile exists, returns cached tile immediately
4. **Cache Miss**: If tile doesn't exist:
   - Downloads tile from MapTiler API
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
- Light theme: Standard MapTiler streets
- Dark theme: Dark variant of MapTiler streets
- Theme selection via query parameter

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

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📞 Support

For issues and questions, please create an issue in the repository or contact the development team.