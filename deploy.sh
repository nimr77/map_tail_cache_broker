
    gcloud run deploy map-broker-jaywalk-prod \
        --project map-broker-jaywalk \
        --region us-central1 \
        --source . \
        --allow-unauthenticated \
        --min-instances 0 \
        --concurrency 50 \
        --platform managed
    
