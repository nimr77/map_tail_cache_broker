# Storage Services Component

Abstracts Google Cloud Storage operations for caching map tiles.

## Files
- `downloader.go`: Utilities for downloading files
- `service.go`: Storage service client and configuration
- `uploader.go`: Handles file uploads to storage
- `util.go`: Utility functions for storage operations

## Purpose
Provides a unified interface for interacting with cloud storage, including upload, download, and utility functions.

## How to Read
- Start with `service.go` for client setup
- Use `downloader.go` and `uploader.go` for file operations
- Reference `util.go` for helpers
