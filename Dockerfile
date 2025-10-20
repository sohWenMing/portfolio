FROM golang
# makes a working directory, from which the docker compose command is called
# this also makes the working directory into $HOME/app - WITHIN the golang image
WORKDIR /app

# first bring over go.mod and go.sum files into $HOME/app, so we can get all dependencies
COPY go.mod go.sum ./
RUN go mod download

# now copy the rest of the files, from the host working directory, into $HOME/app in the image
COPY . .

# build out the binary from ./cmd/server/ (using main.go from there)
# store in in the destination $HOME/app/server
RUN go build -v -o ./server ./cmd/server
CMD ./server

