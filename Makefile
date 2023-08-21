build:
	docker build -t poker:latest -f build/package/Dockerfile .

build-debug:
	docker build -t poker:debug -f build/package/debug/Dockerfile .

all: build build-debug

clean:
	docker image rm poker:latest poker:debug 2> /dev/null || true
