set positional-arguments

default:
    @just --list

# Build bin/grove. Put bin/ first on PATH to run this build as `grove`.
build:
    go build -o bin/grove ./cmd/grove

# Build bin/grove and print the line putting it first on PATH: eval "$(just dev-path)"
dev-path: build
    @echo 'export PATH="{{justfile_directory()}}/bin:$PATH"'

# Install this checkout's build as the everyday grove, in ~/.local/bin.
install:
    go build -o ~/.local/bin/grove ./cmd/grove

test:
    go test ./...

# gofmt, vet, tests, and the line and word budgets (CONSTITUTION.md rule 12).
check:
    @test -z "$(gofmt -l .)" || { gofmt -l .; echo "gofmt: files need formatting"; exit 1; }
    go vet ./...
    go test ./...
    @code=$(find . -name '*.go' ! -name '*_test.go' -exec cat {} + | wc -l); \
     tests=$(find . -name '*_test.go' -exec cat {} + | wc -l); \
     words=$(find internal/instruct \( -name '*.md' -o -name '*.yaml' \) -exec cat {} + | wc -w); \
     echo "lines: $code Go, $tests test; words: $words instruction"; \
     test "$code" -le 5000 && test "$tests" -le 3000 || { echo "over the line budget"; exit 1; }; \
     test "$words" -le 7000 || { echo "over the word budget"; exit 1; }

# Run the skill trials of dev/trials.md in Claude Code and Codex: just trials [case ...]
trials *cases:
    dev/trials/run.sh "$@"

# Create the ledger fixture as a fresh Git repository at DIR, which must not exist.
fixture dir:
    mkdir "$1"
    cp -R testdata/ledger/. "$1"
    cd "$1" && git init -q -b main && git add -A && \
      git -c user.name="Fixture Owner" -c user.email=fixture@example.com commit -q -m "Ledger 0.3.0"
