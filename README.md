# sfetch
CLI application to display system information for all hosts accessible over ssh.

## Installation
Fetch the package to your `$GOPATH`
```bash
go get github.com/Kafva/sfetch
```
and run `./scripts/build-release.sh` to install the binary in `$GOPATH/bin`.

## Notes
The program relies on shell scripts (`./scripts/info.sh` etc.) to determine system information. In the release build these scripts are embedded inside the source code as string literals to make the binary more portable. Using
```bash
go run main.go
```
will read the scripts directly from disk instead, which is preferable during development.
