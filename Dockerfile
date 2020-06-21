# Sttart from base image 1.12.17:
FROM golang:1.12.17

# Configure the repo url so we can configure our work directory:
ENV REPO_URL=github.com/dzikrisyafi/kursusvirtual_oauth-api

# Setup out $GOPATH
ENV GOPATH=/golang

ENV APP_PATH=$GOPATH/src/$REPO_URL

# /app/src/github.com/dzikrisyafi/kursusvirtual_oauth-api/src

# Copy the entire source code from the current directory to $WORKPATH
ENV WORKPATH=$APP_PATH/src
COPY src $WORKPATH
WORKDIR $WORKPATH

RUN go build -o oauth-api .

# Expose port 8000 to the world:
EXPOSE 8000

CMD ["./oauth-api"]