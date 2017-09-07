## go-durationfmt

go-durationfmt is a [Go](https://golang.org/) library that allows you to format [durations](https://golang.org/pkg/time/#Duration) according to a format string. You can specify years, months, days, hours, minutes, and seconds in the format string.

See the [Duration Format](#duration-format) section for an explanation of the format string rules.

See the [Caveats](#caveats) section before using this library.

### Duration Format

The format string format that go-durationfmt uses is similar to that used by [Go's fmt package](https://golang.org/pkg/fmt/#hdr-Printing). The `%` character signifies that the next character is a modifier that specifies a particular duration unit. The following is the full list of modifiers supported by go-durationfmt:

* `%y` - # of years
* `%w` - # of weeks
* `%d` - # of days
* `%h` - # of hours
* `%m` - # of minutes
* `%s` - # of seconds
* `%%` - print a percent sign

You can place a `0` before the `h`, `m`, and `s` modifiers to zeropad those values to two digits. Zeropadding is undefined for the other modifiers.

#### Format String Examples

The following examples show how a duration of 42 hours, 4 minutes, and 2 seconds will be formatted with various format strings:

| Format string          | Output             |
|------------------------|--------------------|
| `%d days, %h hours`    | `1 days, 18 hours` |
| `%m minutes`           | `2524 minutes`     |
| `%s seconds`           | `151442 seconds`   |
| `%d days, %0h:%0m:%0s` | `1 days, 18:04:02` |

### Get

Fetch and build go-durationfmt:

```
go get github.com/davidscholberg/go-durationfmt
```

### Usage

Here's a simple example for using go-durationfmt:

```go
package main

import (
    "fmt"
    "github.com/davidscholberg/go-durationfmt"
    "time"
)

func main() {
    d := (42 * time.Hour) + (4 * time.Minute) + (2 * time.Second)
    durStr, err := durationfmt.Format(d, "%d days, %h hours")
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(durStr)
    }
}
```

### Caveats

* go-durationfmt assumes that there are 24 hours in a day and 365 days in a year, which is not always the case due to such pesky things as leap years, leap seconds, and daylight savings time. **If you need your durations to account for such time jumps, then do not use this library.**
* go-durationfmt returns durations as integer values, so any fractional durations are truncated.
