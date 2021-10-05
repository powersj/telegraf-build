# telegraf-mod

Use to build Telegraf with only the specific input, output, processor, and
aggregator plugins enabled. By default, builds with all plugins disabled.

## How to Use

For Telegraf to be useful, at least one input and one output are required.
Additionally the path to the Telegraf source to build is required if not run
from that directory:

```shell
telegraf-mod --source ~/telegraf --inputs cpu --outputs file
```

The above would build Telegraf from the user's home directory with only the CPU
input and file output.

Multiple plugins can be specified either by additional flags or
comma-separated. The two commands below are identical:

```shell
telegraf-mod --source ~/telegraf --inputs cpu,mem --outputs file
telegraf-mod --source ~/telegraf --inputs cpu --inputs mem --outputs file
```

Users can pass one or more valid Telegraf TOML files to see what plugins are
required:

```shell
telegraf-mod --source ~/telegraf --config <file> [--config <file>]
```

Using both config files and explicit declarations are also possible. The final
telegraf will have the union of all flags.

## Builds

To build Telegraf the `make` command is run in the Telegraf source. If users
have issues building Telegraf, try running the `make` command directly in the
source directory yourself first.

### Important Files

The files that control whether a plugin gets built or not are as follows:

* `./plugins/inputs/all/all.go`
* `./plugins/outputs/all/all.go`
* `./plugins/processor/all/all.go`
* `./plugins/aggregators/all/all.go`

The format of these files is essentially a blank import per plugin:

```go
package all

import (
    _ "github.com/influxdata/telegraf/plugins/{aggregators|inputs|outputs|processors}/{name}"
)
```
