set positional-arguments

default:
    @just --list

# Build bin/grove-next2.
build:
    go build -o bin/grove-next2 ./cmd/grove-next2

test:
    go test ./...

# gofmt, vet, tests and the line budget (CONSTITUTION.md rule 12).
check:
    @test -z "$(gofmt -l .)" || { gofmt -l .; echo "gofmt: files need formatting"; exit 1; }
    go vet ./...
    go test ./...
    @code=$(find . -name '*.go' ! -name '*_test.go' -exec cat {} + | wc -l); \
     tests=$(find . -name '*_test.go' -exec cat {} + | wc -l); \
     echo "lines: $code Go, $tests test"; \
     test "$code" -le 5000 && test "$tests" -le 3000 || { echo "over the line budget"; exit 1; }

# Create the ledger fixture as a fresh Git repository at DIR, which must not exist.
fixture dir:
    mkdir "$1"
    cp -R testdata/ledger/. "$1"
    cd "$1" && git init -q -b main && git add -A && \
      git -c user.name="Fixture Owner" -c user.email=fixture@example.com commit -q -m "Ledger 0.3.0"
