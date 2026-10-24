# Install

Grove is one program, `grove`, for macOS and Linux. Install it one of these
ways.

## Pinned per project, with mise

This is the way for a team. With [mise](https://mise.jdx.dev), add Grove to
your project's `mise.toml`:

```toml
[tools]
"github:mascah/grove" = "0.1.0"
```

Then run `mise install`. Everyone who works on the project, and its CI, runs
the same version. Upgrading is a one-line pull request. In GitHub
Actions, `jdx/mise-action` installs what `mise.toml` lists, and later steps
can run `grove`:

```yaml
- uses: jdx/mise-action@v5
```

## The install script

Without mise:

```sh
curl -fsSL https://raw.githubusercontent.com/mascah/grove/main/install.sh | sh
```

The script installs the latest release in `~/.local/bin`, or in the folder
`GROVE_INSTALL_DIR` names. It tells you when that folder is not on your
PATH. It checks the download against the release's checksums, and installs
nothing if they differ. To install a given version, name its tag:

```sh
curl -fsSL https://raw.githubusercontent.com/mascah/grove/main/install.sh | sh -s -- v0.1.0
```

The script installs one version per machine, so it does not pin a project's
version. Run it again to upgrade. In CI, pass the version, so that a new
release does not change what your checks run. In GitHub Actions:

```yaml
- run: curl -fsSL https://raw.githubusercontent.com/mascah/grove/v0.1.0/install.sh | sh -s -- v0.1.0
- run: echo "$HOME/.local/bin" >> "$GITHUB_PATH"
```

The script's address names the tag too, so a later change to the script does
not change what this step runs.

## With Go

```sh
go install github.com/mascah/grove/cmd/grove@v0.1.0
```

Go puts `grove` in its bin folder, `$(go env GOPATH)/bin` unless you set
`GOBIN`.

## Download

From the [releases page](https://github.com/mascah/grove/releases), download
`checksums.txt` and the archive for your system,
`grove_<os>_<arch>.tar.gz`, where `<os>` is `darwin` or `linux` and `<arch>`
is `amd64` or `arm64`. Check the archive, in the folder you downloaded to:

```sh
sha256sum --check --ignore-missing checksums.txt
```

On macOS, use `shasum -a 256 --check --ignore-missing checksums.txt`. Unpack
the archive and put `grove` in a folder on your PATH.

## Check it

However you installed it, check it with:

```sh
grove version
```
