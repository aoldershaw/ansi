# ansi

ansi is a pluggable solution for parsing and interpreting streams of ANSI text.
The intended use-case is for processing streams of log output for 
[Concourse](github.com/concourse/concourse).

## Usage

```go
import (
    github.com/aoldershaw/ansi
)
...

output := &ansi.InMemory{}
interpreter := ansi.New(output)

interpreter.Parse([]byte("\x1b[1mbold\x1b[m not bold"))
interpreter.Parse([]byte("\nline 2"))

linesJSON, _ := json.MarshalIndent(output.Lines, "", "  ")
fmt.Println(string(linesJSON))
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

Currently, the only provided output method is `ansi.InMemory`,
which stores all the lines of text in memory. A line is represented
as an `ansi.Line`, which is a slice of `ansi.Chunk`. `ansi.Chunk`s
are intended to be concatenated in order.

### Parser

The parser can also be used independently of the interpreter.

```go
import (
    github.com/aoldershaw/ansi/action
    github.com/aoldershaw/ansi/parser
)
...

callback := action.HandlerFunc(func(a action.Action) {
    switch v := a.(type) {
        case action.Print:
            ...
        ...
    }
})
p := parser.New()

p.Parse([]byte("some bytes"))
```

If you prefer to work with channels instead of callbacks, you
can use the convenience constructor `parser.NewWithChan`

```go

callback := action.HandlerFunc(func(a action.Action) {
    switch v := action.(type) {
        case action.Print:
            ...
        ...
    }
})
p, actions, done := parser.NewWithChan()
defer done()
go func() {
    for {
        a := <-actions
        switch v := a.(type) {
            ...
        }
    }
}()
p.Parse([]byte("some bytes"))
```

## Installation

```shell script
go get -u github.com/aoldershaw/ansi
```