# go-version-lib

A Go library for parsing, comparing, and manipulating software versions across many package ecosystems.

## Features

- Parse, compare, and bump semantic versions (`semver`)
- Support for ecosystem-specific version schemes:
  - Debian (`deb`)
  - RubyGems (`gem`)
  - Maven (`mvn`)
  - npm (`npm`)
  - PEP 440 (`pep440`)
  - RPM (`rpm`)
  - Package URL (`purl`)
- Generic version handling (`version`)
- Constraints, requirements, and version collections

## Installation

```sh
go get github.com/khulnasoft/go-version-lib
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/khulnasoft/go-version-lib/pkg/semver"
)

func main() {
	v1, _ := semver.Parse("1.2.3")
	v2, _ := semver.Parse("1.2.4")

	fmt.Println(v1.LessThan(v2)) // true
	fmt.Println(v1.String())     // 1.2.3
}
```

## Packages

| Package | Description |
|---------|-------------|
| `pkg/version` | Core version parsing and comparison |
| `pkg/semver` | Semantic Versioning support |
| `pkg/deb` | Debian package versions |
| `pkg/gem` | RubyGems versions |
| `pkg/mvn` | Maven artifact versions |
| `pkg/npm` | npm package versions |
| `pkg/pep440` | PEP 440 (Python) versions |
| `pkg/rpm` | RPM package versions |
| `pkg/purl` | Package URLs |
| `pkg/part` | Version part primitives |
| `pkg/prerelease` | Pre-release comparison |
| `pkg/reference` | Version references |

## Development

```sh
make bootstrap   # install tooling
make test        # run unit tests
make lint        # run static analysis
make format      # format source
```

## License

MIT License. See [LICENSE](LICENSE).
