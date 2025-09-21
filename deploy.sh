
#!/bin/bash

# Set script to exit on any error
set -e

# Colors for logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_info "Starting deployment of map-broker-jaywalk-prod..."

# Check if .env file exists
if [ -f ".env" ]; then
    log_info "Found .env file, will be included in deployment"
else
    log_warning ".env file not found - deployment will proceed without environment variables"
fi

# Check if gcloud is authenticated
log_info "Checking gcloud authentication..."
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    log_error "No active gcloud authentication found. Please run 'gcloud auth login'"
    exit 1
fi

log_success "gcloud authentication verified"

# Deploy to Google Cloud Run
log_info "Deploying to Google Cloud Run..."
gcloud run deploy map-broker-jaywalk-prod \
    --project map-broker-jaywalk \
    --region us-central1 \
    --source . \
    --allow-unauthenticated \
    --min-instances 0 \
    --concurrency 50 \
    --platform managed \
    --quiet

if [ $? -eq 0 ]; then
    log_success "Deployment completed successfully!"
    
    # Get the service URL
    SERVICE_URL=$(gcloud run services describe map-broker-jaywalk-prod \
        --project map-broker-jaywalk \
        --region us-central1 \
        --format="value(status.url)")
    
    log_info "Service URL: $SERVICE_URL"
else
    log_error "Deployment failed!"
    exit 1
fi

log_success "Deployment script completed"
