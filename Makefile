all: compile
pre:
	go get github.com/Jeffail/gabs/v2
	go get github.com/aidarkhanov/nanoid/v2
	go get github.com/mattn/go-sqlite3
	go get github.com/stretchr/testify/assert
	go get github.com/jinzhu/now
	go get github.com/PaesslerAG/gval
c:
	go build  -v ./...
test:
	go test -v
rebuild: pre c test
compile: c test
