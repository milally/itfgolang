# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

WORKDIR /go/src/app

COPY . .

# Build the itfgolang command inside the container.
RUN go install github.com/milally/itfgolang

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

# Run the itfgolang command by default when the container starts.
#ENTRYPOINT /go/bin/itfgolang

# Document that the service listens on port 8080.
EXPOSE 8080

CMD ["go-wrapper", "run"] # ["app"]
