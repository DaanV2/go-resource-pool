# Go Resource Pools

[![Pipeline](https://github.com/DaanV2/go-resource-pool/actions/workflows/pipeline.yaml/badge.svg)](https://github.com/DaanV2/go-resource-pool/actions/workflows/pipeline.yaml)

A simple library for resource pools of undefined types in Go.

## Installation

```bash
go get github.com/DaanV2/go-resource-pool
```

## Usage

```go
package main

import (
    "github.com/DaanV2/go-resource-pool"
)

// Amount of locks to use best to use amount of threads * 10
lockAmount := 100
items := make([]Processor, lockAmount)
//TODO add items to the slice

pool := pools.NewResourcePool[Processor](items, lockAmount)

// Accessing items:

item, returnFn := pool.Loan(index)

// Do something with the item

returnFn(item) // <= This allows others to use the item when your done, but you can replace, update or set the item to nil.


// Incase you want to grab items based on ID's or identifiy resources:

item, returnFn := pool.LoanByUint64(id)
item, returnFn := pool.LoanByString("final-filename.docx")
item, returnFn := pool.LoanByBytes([]byte{0x00, 0x01, 0x02, 0x03})
```