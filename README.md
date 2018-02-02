# Morton Pack/Unpack Library [![GoDoc](https://godoc.org/github.com/gojuno/go.morton?status.svg)](http://godoc.org/github.com/gojuno/go.morton) [![Build Status](https://travis-ci.org/gojuno/go.morton.svg?branch=master)](https://travis-ci.org/gojuno/go.morton)

## Basics

Check [[https://en.wikipedia.org/wiki/Z-order_curve][wikipedia]] for details.

## Example

```
import "github.com/gojuno/go.morton"

m := morton.Make64(2, 32) // 2 dimenstions 32 bits each
code := m.Pack(13, 42)    // pack two values
values := m.Unpack(code)  // should get back 13 and 42
```
