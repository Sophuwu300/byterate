# byterate

Library and command line tool to calculate the time it takes to transfer a file at a given rate.

## Library Usage

```go
package main

import (
    "fmt"
    "sophuwu.site/byterate"
)

func main() {
    size, err := byterate.ParseSize("1gb")
	if err != nil {
		panic(err)
	}
    rate, err := byterate.ParseSize("1mbit")
	if err != nil {
		panic(err)
    }
    endTime, duration, err := byterate.Time(size, rate)
	if err != nil {
		panic(err)
    }
    fmt.Printf("End Time: %s\n", endTime)
    fmt.Printf("Duration: %s\n", duration)
	}
}
```

## Command Line Usage

```sh
byterate [options] <size> <rate>
```

### Options:
+ `-h` `--help` Show the help message.
+ `-d` `--duration` Print the duration of the transfer.
+ `-t` `--time` Print the time the transfer will end.

if no options are given, duration will be printed.

### Arguments:
+ `<size>` is the size of the file to transfer.
+ `<rate>` is the transfer rate as a size, always per second.

### Arguments Format:
`<size>` and `<rate>` are numbers with optional SI prefixes and units.

#### Supported units:
+ `b` for bytes
+ `bit` for bits

If no unit is given, bytes are assumed.

#### SI Prefixes:
The prefixes are case-insensitive. The prefixes are in base 10 by default. To use base 2, add an `i` to the prefix.

Base 10 prefixes: `k` `m` `g` `t` `p` `e` `z` `y`

Base 2 prefixes: `ki` `mi` `gi` `ti` `pi` `ei` `zi` `yi`

### Examples:

Duration of 10 GiB at 120 mbps:\
`byterate 10gib 120mbit`

The completion time of 16 MiB at 1.2 MiB/s:\
`byterate -t 16mib 1.2mib`

The duration and completion time of 15 GB at 1.5 MB/s:\
`byterate -dt 15g 1.5m`

## Installation

```sh
git clone sophuwu.site/byterate
cd byterate
go build -trimpath -ldflags="-s -w" -o byterate
sudo install byterate /usr/local/bin/byterate
```

## License

MIT