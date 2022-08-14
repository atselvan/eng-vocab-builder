
pre:
	go generate ./...

build:
	go build

run: pre build
	./eng-vocab-builder

gox: pre
	gox -osarch="linux/arm64 linux/amd64 darwin/amd64" -output='bin/{{.Dir}}-{{.OS}}-{{.Arch}}'

docker-build:
	gox -osarch="linux/amd64" -output='bin/{{.Dir}}-{{.OS}}'
	docker container rm -f eng-vocab-builder
	docker build --tag allanselvan/eng-vocab-builder .
	rm -rf bin

docker-run: docker-build
	docker container rm -f eng-vocab-builder
	docker container run --name eng-vocab-builder \
	-p 8000:8000 \
	-e ANKI_CONNECT_URL="http://192.168.2.18:8765" \
	-e ANKI_DECK_NAME="New Deck" \
	-e ANKI_DECK_MODEL="Basic-a39a1" \
	-d allanselvan/eng-vocab-builder

docker-push: docker-build
	docker push allanselvan/eng-vocab-builder

docker-buildx: 
	gox -osarch="linux/arm64" -output='bin/{{.Dir}}-{{.OS}}'
	docker buildx build --platform linux/arm64/v8 --tag allanselvan/eng-vocab-builder --push .
	rm -rf bin