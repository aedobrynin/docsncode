# DocsnCode

DocsnCode is a tool to unite code with documentation. You can write code and explanatory text with images, diagrams and hyperlinks **at the same time**.

For example, for file like this:
```
package main

import "fmt"

// @docsncode
// ```mermaid
// graph TD;
//	A-->B;
//	A-->C;
//	B-->D;
//	C-->D;
// ```
// @docsncode

func main() {
    fmt.Println("Hello, world!")
}

```

You can get a result HTML-file with content

![example result](https://github.com/user-attachments/assets/72da1484-a526-4dc9-81ca-d5b6a8b11bfd)

## How to build

Build requires Go (it's tested with version 1.24.1). Just simply run `go build` at the repository root.

## How to run

`./docsncode project result`

There is other run parameters that are described in [docs](TODO).

## Main features

TODO: list main features

## Documentation

TODO: add link to documentation
