docker push us-central1-docker.pkg.dev/the-new-orleans-connection/com-thenolaconnect-images/app:${TAG};

gcloud artifacts docker images list us-central1-docker.pkg.dev/the-new-orleans-connection/com-thenolaconnect-images/app --include-tags --format=json | jq '.[] | .';