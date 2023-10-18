# Go variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
BINARY_NAME = bin/myapp

all: build

build: create-bin
	$(GOBUILD) -o $(BINARY_NAME) -v

create-bin:
	@mkdir -p bin

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

test:
	$(GOTEST) ./...
