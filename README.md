# itfgolang
This is a demo created for the Inovation and Tech Forum presentation on
GopherCon and Golang. The purpose of this demo is to show the power of
running native Go code as a container deployment and showing the ease
with which a deployment might be updated.

The Dockerfile specifies what we use to build our custom Docker image
on dockerhub.com. My dockerhub account is set to automatically build 
from this Github.

The containers.yaml file specifies the Kubernetes deployment characteristics.

The main.go file is our Go code which presents the container ID and code version.
