all: minibank mariadb

minibank: bin/minibank
mysql: minibankdb
mariadb: minibankdb

bin/minibank: $(shell find src/minibank -name *.go)
	docker run --rm -v `pwd`:/go -w /go/src/minibank golang:1.9-alpine sh -c 'apk add git && go get && go build -v'
	docker build -t minibank .

minibankdb:
	docker build -t minibankdb db

run-images: minibank minibankdb db-password.txt
	docker kill minibank minibankdb || true
	docker rm minibank minibankdb || true
	docker run --rm -d --name minibankdb -e MYSQL_ROOT_PASSWORD="`cat db-password.txt`" minibankdb
	docker run --rm -d --name minibank --link minibankdb -p 80:80 minibank
	docker ps -a

db-password.txt:
	tr -dc a-zA-Z0-9@^_= < /dev/urandom | head -c 13 > $@

clean:
	sudo rm -rf bin pkg src/golang.org src/github.com src/minibank/minibank db-password.txt
	# TODO: remove the Docker images

.PHONY: all clean minibank mysql mariadb minibankdb run-images
