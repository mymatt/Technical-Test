## Test 1

- the files for Test 1 are located in the **TEST1** folder

### Pre-Requisites
- golang installation
- docker installation

### Golang Testing

- We will test the golang application locally prior to using docker and incorporating into a Travis pipeline
- Firstly, clone the repo to a local directory
```
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```
- Make sure that your GOLANG PATH, and PATH has been correctly set
- compile the golang file to create the golang-test binary
```
go build -o golang-test .
```
- run the golang application
```
./golang-test
```
- Using postman to test the API, the below tests show the results for each of the endpoints /, /go, and /opt

**endpoint /**
![Image /](https://github.com/mymatt/Technical-Test/blob/master/images/Post1.png)
**endpoint /go**
![Image /go](https://github.com/mymatt/Technical-Test/blob/master/images/Post2.png)
**endpoint /opt**
![Image /opt](https://github.com/mymatt/Technical-Test/blob/master/images/Post3.png)
### Dockerfile Testing

- Lets test that the application can be successfully run in a docker container on our local machine
- make sure you’re in the directory containing the dockerfile and run the docker build command with a tag to identify the image
```
docker build -t test:test .
```
- check image
```
docker image ls
```
- run container with the following options: name test, detached and removed when stopped
```
docker run --name test --rm -d -p 8000:8000 test:test
```
- now that the container is running, we can access the api using postman or curl
```
curl localhost:8000
```
- We get the following response
![Image Reply](https://github.com/mymatt/Technical-Test/blob/master/images/Reply1.png)
- Our server is running however we cannot connect to it on port 8000
- Lets see if we can access the api from within the container. We are using alpine image so no bash shell
```
docker exec -ti test sh
```
- install curl and check port 8000
```
apk add curl
curl localhost:8000
```
Returns **Hello, world.**

- on examining the code for the server, we can see that it is being run with address 127.0.0.1 which allows access to port 8000 only from the localhost, that is, within the container
```
s := &http.Server{
  Handler:      r,
  Addr:         "127.0.0.1:8000",
```
- changing the address to 0.0.0.0, which allows access from all addresses
```
s := &http.Server{
  Handler:      r,
  Addr:         "0.0.0.0:8000",
```
- cleanup first
```
docker container stop test
docker rmi test:test
```
- re-build, and run the container again
- test outside container with curl
```
curl localhost:8000
```
Returns **Hello, world.**

### Multi-Stage Build

- A multi-stage docker build involves multiple FROM statements that use different base images, which in turn allows the copying of artifacts only and disposing of the remaining image from the previous stage.
- The aim is reduce an image to its smallest size
- To evaluate the effectiveness of multi-stage builds, lets asses the size and layers of an image build not using multi-staging
- Lets check the size created with the current dockerfile
```
docker images test:test --format "{{.Repository}}:{{.Tag}} {{.Size}}"
```
![Image Size1](https://github.com/mymatt/Technical-Test/blob/master/images/Size1.png)  
- We can also view the layers and their sizes
```
docker history test:test
```
![Image Layers1](https://github.com/mymatt/Technical-Test/blob/master/images/Layers1.png)

### Defining Stages
- Firstly we need to label a stage, in which we can later refer to. We use AS to create a stage label
```
FROM golang:alpine AS builder
```
- we can copy from this stage using the --from option
- in this case we are copying the binary to the container
```
COPY --from=builder /app/golang-test .
```
### Scratch Image

- when copying our binary to our second stage, we define that stage with a FROM instruction that defines the image we use
- The scratch image is the most minimal image that is used as a starting point for building container
- This image is empty and contains no files or folders
- To use a scratch image that contains no dependencies for our application, we need to be able to create a binary that does not dynamically link to dependencies, which is possibly by disabling CGO
- Disabling CGO leads to a self contained binary (application and dependencies)
```
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o golang-test .
```
- Lets check the size created with the new multi-stage dockerfile
```
docker images multi:multi --format "{{.Repository}}:{{.Tag}} {{.Size}}"
```
![Image Size2](https://github.com/mymatt/Technical-Test/blob/master/images/Size2.png)
- We can also view the layers and their sizes
```
docker history multi:multi
```
![Image Layers2](https://github.com/mymatt/Technical-Test/blob/master/images/Layers2.png)
- we can see that the number and size of each layer has been reduced considerably

### Optimizing Dockerfile

- In addition to multi-stage builds some general dockerfile optimization techniques are:
1) Combine run commands using && to minimize layers
2) Keep stable instructions at top of dockerfile and place additions at the bottom
3) Using go mod download command, which takes go.mod and go.sum files and downloads the dependencies from them instead of using the source code. These files won't be changed that often so we can cache them earlier in the build process
