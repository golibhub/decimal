# Decimal

This package provides a decimal type that can represent decimal numbers with arbitrary precision. Currently,
it uses [cockroachdb/apd](https://github.com/cockroachdb/apd) package under the hood.

## Features

* `Decimal` is safe to use without initialization, it's 0 by default.
* Easy to use, it implements all the arithmetic operations.
* JSON marshaling and unmarshalling support.
* SQL scanning and value support.

## Caveats

- Most of the methods can panic. It's not ideal, but it's a trade-off to provider a simpler API.
- `Decimal.Div` method uses a default precision of 32, which can be customized by setting a new value
  to `DivPrecision` variable.
