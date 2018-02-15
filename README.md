# Unique ID Generator

Use 6 byte for submission time (UTC). You can use for SQL.

**! Limit values to 2310 year.**

#### Identifier structure:

- `4-byte` value representing the seconds since the 2018-01-01 epoch

- `1-byte` hardware address CRC4 ID

- `2-byte` process id

- `3-byte` counter (random, or manual)

- `2-byte` nanoseconds (the first 2 bytes of the four byte value)

## Install

    go get github.com/itrabbit/rid

## Usage

```go
fmt.Println(rid.New())
// -> 006yxbosapiby5gfy6gg

fmt.Println(rid.New().String())
// -> 006yxfysaq0bynaq2070

fmt.Println(rid.New().NumeralString())
// -> 000013238191218086000191085088016014
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