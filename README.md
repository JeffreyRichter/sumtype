# sumtype

A Go package that provides casting capabilities between sum type projections using unsafe pointers.

## Overview

This package provides a `Caster` type that allows casting between a JSON-serializable struct and its variants. All variants must have the same memory layout, with the first field being a non-exported `xxxCaster` field of type `sumtype.Caster[Json]`.

## Features

- Type-safe casting between sum type projections
- JSON marshaling/unmarshaling support
- String representation with pretty-printed JSON
- Zero out non-relevant fields for specific variants

## Usage

```go
package main

import "github.com/JeffreyRichter/sumtype"

// Define your JSON struct and variants here
// See the test files for examples
```

## Installation

```bash
go get github.com/JeffreyRichter/sumtype
```

## Requirements

- Go 1.25 or later

## License

[Add your license here]