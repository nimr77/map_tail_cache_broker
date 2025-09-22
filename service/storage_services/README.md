# Storage Services Component

A comprehensive storage abstraction layer that provides unified interfaces for cloud storage operations, specifically designed for caching map tiles. Currently implements Google Cloud Storage but is architected to support multiple storage providers.

## 🏗️ Architecture Overview

The storage services component follows a modular design pattern that abstracts storage operations behind a clean interface, making it easy to switch between different storage providers without changing the application logic.

```
storage_services/
├── service.go       # Core client and bucket management
├── uploader.go      # File upload operations (public/private)
├── downloader.go    # File download and URL parsing
└── util.go          # Utility functions for file operations
```

## 📁 File Descriptions

### `service.go` - Core Storage Client
**Purpose**: Manages Google Cloud Storage client initialization and bucket access.

**Key Functions**:
- `GetClient()`: Initializes and returns a Google Cloud Storage client
- `GetMapUploadingBucket()`: Returns a handle to the `map-cached` bucket

**Authentication Strategy**:
- **Primary**: Uses service account key file (`map-broker-jaywalk-75c83aba05cf.json`)
- **Fallback**: Uses default credentials (useful for development and cloud environments)

**Configuration**:
- Bucket name: `map-cached` (hardcoded)
- Credentials file: `map-broker-jaywalk-75c83aba05cf.json`

### `uploader.go` - File Upload Operations
**Purpose**: Handles all file upload operations to cloud storage with different access levels.

**Key Functions**:

1. **`UploadFile(file, fileHeader, path)`**
   - Uploads multipart files with private access
   - Sets content type from file header
   - Returns media link URL

2. **`UploadFilePublic(file, fileHeader, path)`**
   - Uploads multipart files with public read access
   - Sets ACL to allow all users to read
   - Returns media link URL

3. **`UploadFileBytes(fileBytes, contentType, path)`**
   - Uploads byte arrays with private access
   - Useful for programmatically generated content
   - Returns media link URL

4. **`UploadFilePublicBytes(fileBytes, contentType, path)`**
   - Uploads byte arrays with public read access
   - Sets ACL to allow all users to read
   - Returns media link URL

**Access Control**:
- **Private**: No ACL set (default bucket permissions apply)
- **Public**: `storage.AllUsers` with `storage.RoleReader` permission

### `downloader.go` - File Download Operations
**Purpose**: Handles file downloads from cloud storage with URL parsing capabilities.

**Key Functions**:

1. **`DownloadFile(url)`**
   - Downloads files from Firebase Storage URLs
   - Converts URL to storage path using `FromUrlToFileStoragePath`
   - Returns file content as byte array
   - Validates file is not empty

2. **`FromUrlToFileStoragePath(urlString)`**
   - Converts various URL formats to storage paths
   - Supports `gs://` URLs
   - Supports Firebase Storage download URLs
   - Handles URL decoding and path normalization

**Supported URL Formats**:
- `gs://bucket-name/path/to/file`
- `https://firebasestorage.googleapis.com/v0/b/bucket/o/path/to/file`

### `util.go` - Utility Functions
**Purpose**: Provides helper functions for file operations and HTTP downloads.

**Key Functions**:

1. **`SaveFile(data, folder, ext, filename)`**
   - Saves string data to local filesystem
   - Creates directories if they don't exist
   - Handles file creation, writing, and syncing

2. **`DownloadFileFromUrl(url)`**
   - Downloads files from HTTP URLs
   - Validates HTTP response status
   - Returns file content as byte array

## 🔧 Configuration

### Google Cloud Storage Setup

1. **Create a Storage Bucket**:
   ```bash
   gsutil mb gs://map-cached
   ```

2. **Set up Authentication**:
   
   **Option A: Service Account Key File**
   - Create a service account in Google Cloud Console
   - Download the JSON key file
   - Place it in the project root as `map-broker-jaywalk-75c83aba05cf.json`
   - Grant the service account Storage Admin permissions

   **Option B: Default Credentials**
   - Use `gcloud auth application-default login`
   - Or set `GOOGLE_APPLICATION_CREDENTIALS` environment variable

3. **Configure Bucket Permissions**:
   ```bash
   # For public read access (if needed)
   gsutil iam ch allUsers:objectViewer gs://map-cached
   ```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GOOGLE_APPLICATION_CREDENTIALS` | Path to service account key file | Optional* |
| `GOOGLE_CLOUD_PROJECT` | Google Cloud project ID | Optional* |

*Required if not using service account key file in project root

## 🚀 Usage Examples

### Map Tile Storage Path Structure

The storage services use a specific path structure for map tiles based on the theme mode system:

**Path Format**: `{provider}/{themeMode}/{z}/{x}/{y}.png`

**Theme Mode Values**:
- `0` = Dark theme (`ThemeModeDark`)
- `1` = Light theme (`ThemeModeLight`)

**Example Paths**:
- Dark theme: `maptiler/0/10/5/3.png`
- Light theme: `maptiler/1/10/5/3.png`

### Map Tile Operations Flow

The storage services support three main operations for map tiles:

#### 1. **Check Operation** - Verify Tile Existence
```go
// Check if a map tile exists in storage
func CheckIfMapTailExistInService(request map_core.MapRequest) (bool, error) {
    path := request.GetMapTilePath() // Returns: "maptiler/0/10/5/3.png"
    
    bucket, err := storage_services.GetMapUploadingBucket()
    if err != nil {
        return false, err
    }
    
    _, err = bucket.Object(path).Attrs(context.Background())
    if err == nil {
        return true, nil  // Tile exists
    }
    
    return false, nil  // Tile doesn't exist
}
```

#### 2. **Download Operation** - Retrieve Cached Tile
```go
// Download existing tile from storage
func GetTailMapFromStorageService(request map_core.MapRequest) ([]byte, error) {
    path := request.GetMapTilePath() // Returns: "maptiler/0/10/5/3.png"
    
    bytes, err := storage_services.DownloadFile(path)
    if err != nil {
        return nil, err
    }
    
    return bytes, nil
}
```

#### 3. **Upload Operation** - Cache New Tile
```go
// Upload new tile to storage (asynchronous)
func UploadMapTileToStorageService(request map_core.MapRequest, fileBytes []byte) (string, error) {
    path := request.GetMapTilePath() // Returns: "maptiler/0/10/5/3.png"
    
    url, err := storage_services.UploadFileBytes(fileBytes, "image/jpeg", path)
    if err != nil {
        log.Println("Error uploading map tile to storage service:", err.Error())
        return "", err
    }
    
    log.Println("Map tile uploaded to storage service at URL:", url)
    return url, nil
}
```

### Complete Map Tile Request Flow

```go
func GetTail(request map_core.MapRequest) ([]byte, error) {
    // Step 1: Check if tile exists in cache
    exist, err := CheckIfMapTailExistInService(request)
    if err != nil {
        return nil, err
    }
    
    if exist {
        // Step 2a: Cache hit - download from storage
        log.Println("Map tile exists in storage service, fetching from there.")
        return GetTailMapFromStorageService(request)
    }
    
    // Step 2b: Cache miss - download from map provider
    downloadingUrl, err := request.GetFullMapTailUrl()
    if err != nil {
        return nil, err
    }
    
    by, err := storage_services.DownloadFileFromUrl(downloadingUrl)
    if err != nil {
        return nil, err
    }
    
    // Step 3: Asynchronously upload to cache for future requests
    go UploadMapTileToStorageService(request, by)
    
    return by, nil
}
```

### Theme Mode Examples

```go
// Create requests for different themes
darkRequest := map_core.MapRequest{
    MapRequestMeta: map_core.MapRequestMeta{
        Provider:  "maptiler",
        ThemeMode: map_core.ThemeModeDark, // 0
    },
    X: "10",
    Y: "5", 
    Z: "3",
}

lightRequest := map_core.MapRequest{
    MapRequestMeta: map_core.MapRequestMeta{
        Provider:  "maptiler",
        ThemeMode: map_core.ThemeModeLight, // 1
    },
    X: "10",
    Y: "5",
    Z: "3", 
}

// These will generate different storage paths:
// darkRequest.GetMapTilePath()  -> "maptiler/0/3/10/5.png"
// lightRequest.GetMapTilePath() -> "maptiler/1/3/10/5.png"
```

### Direct Storage Operations

```go
// Upload a multipart file (private)
file, header, err := r.FormFile("map_tile")
if err != nil {
    return err
}
defer file.Close()

url, err := storage_services.UploadFile(file, header, "maptiler/1/10/5/3.png")
if err != nil {
    return err
}

// Upload bytes with public access
tileBytes := []byte("...") // Your tile data
url, err := storage_services.UploadFilePublicBytes(tileBytes, "image/png", "maptiler/0/10/5/3.png")

// Download from Firebase Storage URL
firebaseURL := "https://firebasestorage.googleapis.com/v0/b/map-cached/o/maptiler%2F0%2F3%2F10%2F5.png"
data, err := storage_services.DownloadFile(firebaseURL)

// Download from gs:// URL
gsURL := "gs://map-cached/maptiler/1/3/10/5.png"
data, err := storage_services.DownloadFile(gsURL)
```

### Utility Functions

```go
// Save file locally
err := storage_services.SaveFile("tile data", "downloaded_files", "png", "tile_10_5_3")

// Download from HTTP URL
data, err := storage_services.DownloadFileFromUrl("https://example.com/tile.png")
```

## 🔄 How to Add Another Storage Provider

The storage services component is designed to be extensible. Here's how to add support for additional storage providers:

### Step 1: Create Provider-Specific Service

Create a new file for your storage provider (e.g., `aws_s3_service.go`, `azure_blob_service.go`):

```go
package storage_services

import (
    "context"
    "github.com/aws/aws-sdk-go/service/s3"
    // Add other necessary imports
)

// S3Client wraps AWS S3 client
type S3Client struct {
    client *s3.S3
    bucket string
}

// GetS3Client initializes AWS S3 client
func GetS3Client() (*S3Client, error) {
    // Initialize AWS S3 client
    // Return S3Client instance
}

// GetMapUploadingBucket returns S3 bucket handle
func (s *S3Client) GetMapUploadingBucket() (*s3.Bucket, error) {
    // Return S3 bucket reference
}
```

### Step 2: Implement Provider-Specific Upload Functions

Create upload functions for your provider:

```go
// UploadFileToS3 uploads file to S3
func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
    client, err := GetS3Client()
    if err != nil {
        return "", err
    }
    
    bucket, err := client.GetMapUploadingBucket()
    if err != nil {
        return "", err
    }
    
    // Implement S3 upload logic
    // Return S3 URL
}
```

### Step 3: Implement Provider-Specific Download Functions

```go
// DownloadFileFromS3 downloads file from S3
func DownloadFileFromS3(url string) ([]byte, error) {
    client, err := GetS3Client()
    if err != nil {
        return nil, err
    }
    
    // Parse S3 URL to get bucket and key
    // Implement S3 download logic
    // Return file bytes
}
```

### Step 4: Add Configuration Support

Add environment variable support for provider selection:

```go
// In service.go, add provider selection logic
func GetStorageProvider() string {
    provider := os.Getenv("STORAGE_PROVIDER")
    if provider == "" {
        return "gcs" // Default to Google Cloud Storage
    }
    return provider
}
```

### Step 5: Update Upload/Download Functions

Modify existing functions to use the selected provider:

```go
// Update UploadFile to support multiple providers
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

### Step 6: Add Provider-Specific Configuration

Add configuration for your new provider:

```bash
# For AWS S3
export STORAGE_PROVIDER="s3"
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
export AWS_REGION="us-east-1"
export AWS_BUCKET_NAME="your-bucket-name"

# For Azure Blob Storage
export STORAGE_PROVIDER="azure"
export AZURE_STORAGE_ACCOUNT="your_account"
export AZURE_STORAGE_KEY="your_key"
export AZURE_CONTAINER_NAME="your-container"
```

### Step 7: Update Dependencies

Add necessary dependencies to `go.mod`:

```bash
go get github.com/aws/aws-sdk-go
go get github.com/Azure/azure-storage-blob-go
```

### Example: Complete AWS S3 Implementation

Here's a complete example for AWS S3:

```go
// aws_s3_service.go
package storage_services

import (
    "context"
    "fmt"
    "os"
    
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
    client *s3.S3
    bucket string
}

func NewS3Service() (*S3Service, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION")),
    })
    if err != nil {
        return nil, err
    }
    
    return &S3Service{
        client: s3.New(sess),
        bucket: os.Getenv("AWS_BUCKET_NAME"),
    }, nil
}

func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
    // Read file content
    bytes := make([]byte, fileHeader.Size)
    file.Read(bytes)
    
    // Upload to S3
    _, err := s.client.PutObject(&s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(path),
        Body:        bytes,
        ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
    })
    
    if err != nil {
        return "", err
    }
    
    // Return S3 URL
    return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, path), nil
}

func (s *S3Service) DownloadFile(path string) ([]byte, error) {
    result, err := s.client.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(path),
    })
    if err != nil {
        return nil, err
    }
    defer result.Body.Close()
    
    // Read all bytes
    bytes := make([]byte, *result.ContentLength)
    _, err = result.Body.Read(bytes)
    return bytes, err
}
```

## 🧪 Testing

Test the storage services with different scenarios:

```bash
# Test upload
curl -X POST -F "file=@tile.png" "http://localhost:8080/upload?path=test/tile.png"

# Test download
curl "http://localhost:8080/download?path=test/tile.png" -o downloaded_tile.png
```

## 🔒 Security Considerations

- **Authentication**: Always use proper authentication methods for your storage provider
- **Access Control**: Be careful with public ACLs - only use when necessary
- **Credentials**: Never commit credentials to version control
- **Environment Variables**: Use environment variables for sensitive configuration
- **URL Validation**: Validate URLs before processing to prevent path traversal attacks

## 📊 Performance Considerations

- **Connection Pooling**: Reuse storage clients when possible
- **Async Operations**: Consider async uploads for better performance
- **Chunked Uploads**: For large files, use chunked uploads
- **Caching**: Implement local caching for frequently accessed files
- **Compression**: Consider compressing files before upload

## 🐛 Troubleshooting

### Common Issues

1. **Authentication Errors**:
   - Verify service account permissions
   - Check credential file path
   - Ensure environment variables are set correctly

2. **Bucket Access Errors**:
   - Verify bucket exists and is accessible
   - Check bucket permissions
   - Ensure correct bucket name

3. **Upload Failures**:
   - Check file size limits
   - Verify content type settings
   - Check network connectivity

4. **Download Failures**:
   - Verify file exists in storage
   - Check URL format
   - Ensure proper URL encoding

### Debug Mode

Enable debug logging by setting environment variable:
```bash
export DEBUG_STORAGE=true
```

## 📝 Contributing

When contributing to storage services:

1. Follow the existing code patterns
2. Add comprehensive error handling
3. Include unit tests for new functions
4. Update this documentation
5. Test with multiple storage providers

## 🔗 Related Documentation

- [Google Cloud Storage Go Client](https://pkg.go.dev/cloud.google.com/go/storage)
- [AWS S3 Go SDK](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example.html)
- [Azure Blob Storage Go SDK](https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-go)
