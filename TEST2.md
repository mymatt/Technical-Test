## Test 2

- Test 2 includes all repo files, excluding the **app** folder
- The **goapi** folder contains the golang api and test files

### Pre-Requisites
- github account setup
- dockerhub account setup with repo allocated
- travis account integrated with github account

![Image dockerhub](https://github.com/mymatt/Technical-Test/blob/master/images/Dockerhub_Account.png)

### Stages
There are 7 stages
1) Commit to GIT master branch
2) TravisCI webhook updates lastest commit to Pipeline
3) Lint Testing / Unit Testing
4) Build using dockerfile
5) Publish to Dockerhub
6) Run new image downloaded from Dockerhub
7) Test endpoints

![Image Stages](https://github.com/mymatt/Technical-Test/blob/master/images/TravisCI.png)

### Setup
1) *Github Setup*
- Github provides the 'Template repository' feature that is used to create a boilerplate for launching projects of a similar nature.
- To create a boilerplate project, go to the repo settings and select "Template repository"

![Image setup](https://github.com/mymatt/Technical-Test/blob/master/images/BoilerPlate1.png)

- To active go to:
```
https://github.com/mymatt/Technical-Test
```
- To re-use repo as a boilerplate project, select "Use this template", select the your account from drop down, and follow instructions

![Image setup](https://github.com/mymatt/Technical-Test/blob/master/images/BoilerPlate2.png)

2) *Travis setup*
- sign in to Travis-ci.com
- select sign up with github
- select account, click green activate, and select repositorie(s) to use with TravisCI
- add .travis.yml file and commit then push to git repo  
- future commits will trigger TravisCI

- set Travis environmental variables: select repo, go to settings
- Enter the following:
```
DOCKER_PASSWORD = Use Dockerhub Credentials
DOCKER_USERNAME = Use Dockerhub Credentials
REPO = Use Dockerhub Repo name (e.g technical-test)
```

![Image env](https://github.com/mymatt/Technical-Test/blob/master/images/TravisEnvVar.png)

## Golang API
- The golang API uses gorilla mux, which is used to match requests to specified endpoints. In this case we are only using the GET http method
- It will be installed during the "go install" command during the build phase within the Dockerfile
- The golang application uses modules, where a go.mod file is specified
- the structure for modules is /bin, /pkg, /appname (where the go.mod and main go file are stored)
- the application binds to port 8080
- nested struct's are used to create the API Example Response
- postman can be used to test the API. The below test shows the golang API working locally, prior to travis integration (uses dummy data for /version)

**View Response from / endpoint**
![Image /](https://github.com/mymatt/Technical-Test/blob/master/images/Rest1.png)

**View Response from /version endpoint**
![Image /version](https://github.com/mymatt/Technical-Test/blob/master/images/Rest2.png)

## Version Endpoint
- The /version endpoint returns "description" & "version" from a metadata.json file using jq (pre-installed on travis VM's)
- During the docker build stage of the Travis pipeline we pass the jq outputs to the docker build process as arguments
```
--build-arg vers=$(jq -r '.version' metadata.json)
--build-arg desc="$(jq -r '.description' metadata.json)"
```
- In the dockerfile we then create environmental variables, during the last stage of the multi-stage build (with scratch image)
```
ARG vers
ENV VERS=$vers
```
- The golang application then using the os lookupEnv function to access the environmental variable
```
version, exists := os.LookupEnv("VERS")
```

- The lastcommitsha is also returned by the /version endpoint. This is achieved using the Travis Default Environmental Variables. In this case: ${TRAVIS_COMMIT}

## To Begin
- commit to git on master branch
- To view logs go to Travis repo, select "Running" on left task bar, and choose a running job

## Testing - Golang Linting
- golint is installed and run during the first travis stage to provide linting of our golang file

## Testing - Golang Unit Testing
- The unit tests for each endpoint test for a http status of **200 OK**
- We use unmarshal to parse the JSON response according to a struct that matches the expected JSON structure
- We can then extract the version, description and lastcommitsha data
- The unit tests then attempt to match test strings for the expected values
- to simulate the container environment that runs our golang api application, the unit testing stage has the following env variables assigned
```
- GO111MODULE=on
- VERS=$(jq -r '.version' metadata.json)
- DESC="$(jq -r '.description' metadata.json)"
- SHA=${TRAVIS_COMMIT}
```
- the output of this stage is
![Image testing](https://github.com/mymatt/Technical-Test/blob/master/images/Testing.png)

## Build Stage
- The output for the Build stage will be

![Image Build](https://github.com/mymatt/Technical-Test/blob/master/images/Build.png)

- The newly created build will be visible in your dockerhub repository

![Image Dockerhub](https://github.com/mymatt/Technical-Test/blob/master/images/DockerhubRepo.png)

## Run Stage
- The run stage downloads our new image from dockerhub and runs the container. The logs will show the return values from requests to / and /version

![Image Last](https://github.com/mymatt/Technical-Test/blob/master/images/Last.png)

## Security Non-privileged User
- We need to thwart attacks to the Docker host using root access.
This can be achieved by launching the container with a non-privileged user.
- In the first “builder” stage, we are using the adduser command found in busybox images (alpine), to create an unprivileged user with no home, password or shell
```
ENV USER=usergo
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nohome" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
```
- In the second stage, which is built on an “empty” scratch image, we copy the users and groups from our first builder stage
```
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
```
- Lastly we use our newly create unprivileged user to run the binary
```
USER usergo:usergo
```
- When using unprivileged users, we need to expose ports that are greater than 1024
```
expose 8080
```

### Versioning
- Fixed tags are used for immutability. The tags are the HASH for each commit. This is to avoid the scenario of pushing new versions to the same tags
- Travis provides an Environment Variable for the hash of each commit: ${TRAVIS_COMMIT}
- In addition, we can use the following git command to retrieve the hash for the last commit
```
git log -1 --pretty=%H
```

### Security - Travis CI Best Practises
- encrypt any access keys, credentials, secrets
- avoid leaking secrets to build logs by ensuring that commands are not used that display secrets, sensitive env variables
- rotate tokens and secretes periodically
- check for any dead packages and remove access. Attackers can learn names of packages and re-register those packages and then use that package as a backdoor
- any breaches of security via logs, should then be followed by removing logs of disclosed breaches
