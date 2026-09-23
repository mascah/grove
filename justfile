set positional-arguments

default:
    @just --list

# Build ./bin/grove from this checkout.
build:
    go build -o bin/grove ./cmd/grove

# Build and run: `just run list --status active`. Unlike `go run`, keeps grove's exit code.
run *args: 
    go run ./cmd/grove "$@"

# Replace the installed ~/.local/bin/grove with a build of this checkout.
install:
    go build -o "$HOME/.local/bin/grove" ./cmd/grove
    "$HOME/.local/bin/grove" version

# Remove local branches merged into main, and their worktrees, after asking.
clean-merged:
    #!/usr/bin/env bash
    # Safe by construction: `git worktree remove` without --force refuses a worktree with
    # changes or untracked files, and `git branch -d` refuses a branch that is not merged.
    set -euo pipefail
    git worktree prune
    current=$(git branch --show-current)
    branches=$(git for-each-ref --format='%(refname:short)' --merged main refs/heads \
        | grep -vxF -e main -e "$current" || true)
    [ -n "$branches" ] || { echo "Nothing merged into main to clean."; exit 0; }
    worktree_of() {
        git worktree list --porcelain | awk -v ref="branch refs/heads/$1" \
            '/^worktree /{wt=substr($0,10)} $0==ref{print wt}'
    }
    for b in $branches; do echo "  $b  $(worktree_of "$b")"; done
    read -rp "Remove these branches and worktrees? [y/N] " ok
    [ "$ok" = y ] || exit 1
    for b in $branches; do
        wt=$(worktree_of "$b")
        if [ -n "$wt" ] && ! git worktree remove "$wt"; then
            echo "kept $b: its worktree was not removed"; continue
        fi
        git branch -d "$b" || echo "kept $b"
    done
