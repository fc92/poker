build:
	docker build -t poker:latest -f backend/build/package/Dockerfile .


all: build

clean:
	docker image rm poker:latest poker:debug 2> /dev/null || true
