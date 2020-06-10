# detect operating system
ifeq ($(OS),Windows_NT)
    CURRENT_OS := Windows
else
    CURRENT_OS := $(shell uname -s)
endif

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -tags netgo -ldflags "-X github.com/orientwalt/htdf/version.GitCommit=${COMMIT_HASH}"
# docker part[single-node part]
DOCKER_TAG = unstable
DOCKER_IMAGE = falcon0125/htdfsingle

build: clean
ifeq ($(CURRENT_OS),Windows)
	go build  $(BUILD_FLAGS) -o ../../build/hsd.exe ../../cmd/hsd
	go build  $(BUILD_FLAGS) -o ../../build/hscli.exe ../../cmd/hscli
else
	go build  $(BUILD_FLAGS) -o ../../build/hsd ../../cmd/hsd
	go build  $(BUILD_FLAGS) -o ../../build/hscli ../../cmd/hscli
endif

clean:
	@rm -rf ../../build

docker: build
	@docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	@docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	@docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_HASH}

run:
	@docker run --rm -it ${DOCKER_IMAGE}:latest hsd start --minimum-gas-prices=0.0001htdf

shell:
	@docker run --rm -it ${DOCKER_IMAGE}:latest sh
