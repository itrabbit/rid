# Unique ID Generator

Use 6 byte for submission time. You can use for SQL.

#### Identifier structure:

- `4-byte` value representing the seconds since the Unix epoch

- `1-byte` hardware address CRC4 ID

- `2-byte` process id

- `3-byte` counter (random, or manual)

- `2-byte` nanoseconds (the first 2 bytes of the four byte value)

## Install

    go get github.com/itrabbit/rid

## Usage

```go
fmt.Println(rid.New())
// -> 2m9p2o6s9xp0561dedpg

fmt.Println(rid.New().String())
// -> 2m9p2pys9xxih20bgc1g

fmt.Println(rid.New().NumeralString())
// -> 021019113097218079133204240151147014
```

Custom source:

```go
src := rid.NewSource()
// src.Seed(0) <- Set start counter by custom value

id := src.NewID()
```

Get embedded info:

```go
id := rid.New()

id.Mid()    // Hardware Address CRC ID
id.Pid()    // Process Pid
id.Time()   // With an accuracy of up to 6 bytes (from 8)

id.Counter() 
```

## License and copyright

Copyright (c) 2018 IT Rabbit.

**[Original encode/decode algorithm](https://github.com/rs/xid):** Copyright (c) 2015 Olivier Poitrey <rs@dailymotion.com>. The same MIT license applies.