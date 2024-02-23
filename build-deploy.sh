# Build and deploy for Debian 11

# We need to set environment variables so go knows how to build the binary
export GOOS=linux
export GOARCH=amd64

# Build the binary
go build -o bartering

# Copy file onto Grid5000 frontend
scp bartering lyon.g5k:./bartering-deployer/playbooks/bartering-protocol/bartering
