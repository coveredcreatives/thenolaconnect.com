#!/bin/bash
#!/bin/bash
go build ./cmd/cli;

docker build -t app:${TAG} -t us-central1-docker.pkg.dev/the-new-orleans-connection/com-thenolaconnect-images/app:${TAG} -f ${DOCKER_FILE} .;
