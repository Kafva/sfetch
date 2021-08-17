# sfetch
CLI application to display system information for all hosts accessible over ssh. Hosts are grouped in a tree based on the `ProxyJump` attribute in `ssh_config`. A terminal with support for [nerdfonts](https://www.nerdfonts.com/font-downloads) is required to render OS icons. See `sfetch --help` for configuration options.

![](./.github/screenshot.png)

## Installation
Fetch the package to your `$GOPATH`
```bash
go get github.com/Kafva/sfetch
```
and run `./scripts/build-release.sh` to install the binary in `$GOPATH/bin`.

## Notes
* The program relies on shell scripts (`./scripts/info.sh` etc.) to determine system information. In the release build these scripts are embedded inside the source code as string literals to make the binary more portable. 
* Motherboard information is only available for BSD hosts if `doas` and `dmidecode` are installed.
* Using
```bash
go run main.go
```
will read the scripts directly from disk instead, which is preferable during development.

