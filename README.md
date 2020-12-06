# ansi

ansi is a pluggable solution for parsing and interpreting streams of ANSI text.
The intended use-case is for processing streams of log output for 
[Concourse](https://github.com/concourse/concourse).

## Usage

```go
import (
    "encoding/json"

    "github.com/aoldershaw/ansi"
)

func main() {
    var lines ansi.Lines
    writer := ansi.NewWriter(&lines)

    writer.Write([]byte("\x1b[1mbold\x1b[m not bold"))
    writer.Write([]byte("\nline 2"))

    linesJSON, _ := json.MarshalIndent(lines, "", "  ")
    fmt.Println(string(linesJSON))
}

```

Output:

```
[
  [
    {
      "data": "bold",
      "style": {
        "bold": true
      }
    },
    {
      "data": " not bold",
      "style": {}
    }
  ],
  [
    {
      "data": "line 2",
      "style": {}
    }
  ]
]
```

Currently, the only provided output method is `ansi.Lines`, which stores all
the lines of text in memory. A line is a slice of `ansi.Chunk` - a stylized
chunk of text. `ansi.Chunk`s are intended to be concatenated in order.

### Parser

The parser can also be used independently of the interpreter.

```go
import (
    "github.com/aoldershaw/ansi"
)

func main() {
    parser := ansi.NewParser()

    input := []byte("some bytes")
    for _, action := range parser.ParseAll(input) {
        switch v := action.(type) {
            case ansi.Print:
                fmt.Println(string(v))
        }
    }
}
```

## Installation

```shell script
go get -u github.com/aoldershaw/ansi
```
