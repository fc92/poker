docker-build:
	docker build -t poker:latest -f backend/build/package/Dockerfile .

ionic-build:
	cd frontend && ionic build

all: ionic-build docker-build

clean:
	docker image rm poker:latest poker:debug 2> /dev/null || true
