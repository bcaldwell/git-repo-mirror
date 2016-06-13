FROM golang:1.7-alpine

RUN apk update && apk add git openssh

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/benjamincaldwell/git-repo-mirror


# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/ghodss/yaml
RUN go install github.com/benjamincaldwell/git-repo-mirror

RUN mkdir /etc/git-repo-mirror; echo "---" > /etc/git-repo-mirror/config.yml

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/git-repo-mirror

# Document that the service listens on port 8080.
EXPOSE 8080
