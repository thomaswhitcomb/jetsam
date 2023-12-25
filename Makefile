all: clean compile test

MODULE := "jetsam"
clean:
	@echo "==> Cleaning up previous builds."
	@rm -rf ./bin/$(MODULE)

compile:
	@echo "==> Compiling source code."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/$(MODULE)

coverage:
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out

fmt:
	@echo "==> Formatting source code."
	@gofmt -w ./

race_compile:
	@echo "==> Compiling source code."
	@go build -v -race -o ./bin/$(MODULE)

test: fmt vet
	@echo "==> Running tests."
	@go test -cover
	@echo "==> Tests complete."

vet:
	@go vet

docker: compile
	docker build .

help:
	@echo "clean\t\tremove previous builds"
	@echo "compile\t\tbuild the code"
	@echo "coverage\tgenerate and view code coverage"
	@echo "fmt\t\tformat the code"
	@echo "race_compile\tbuild the code with race detection"
	@echo "test\t\ttest the code"
	@echo "vet\t\tvet the code"
	@echo ""
	@echo "default will test, format, and compile the code"

.PNONY: all clean deps fmt help test
