# Decorder

A declaration order linter for golang. In case of this tool declarations are `type`, `const`, `var` and `func`.

This linter can check

* the general order of declarations (e.g. all global constants are always defined before variables)
* that `init()` is the first function inside a file (if defined)
* that multiple (global) `const`, `var` and `type` statements are not allowed (go supports e.g. a single `const`
  statement with parenthesis for all constant declarations)

## Installation

```
go get gitlab.com/bosi/decorder/cmd
```

## Usage

```shell
# with default options
./decorder ./...

# custom declaration order
./decorder -dec-order var,const,func,type ./...

# disable declaration order check
./decorder -disable-dec-order-check ./...

# disable check for multiple declarations statements
./decorder -disable-dec-num-check ./...

# disable check that init func is always first function
./decorder -disable-init-func-first-check ./...
```