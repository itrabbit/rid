# Unique ID Generator

Use 6 byte for submission time (UTC). You can use for SQL.

**! Limit values to 2310 year.**

#### Identifier structure:

- `1-byte` hardware address CRC4 ID

- `2-byte` process id

- `4-byte` value representing the seconds since the 2018-01-01 epoch

- `3-byte` counter (random, or manual)

- `2-byte` nanoseconds (the first 2 bytes of the four byte value)

## Install

    go get github.com/itrabbit/rid

## Usage

```go
fmt.Println(rid.New())
// -> t9eti00dzt0e0hkibo9g

fmt.Println(rid.New().String())
// -> t9ew400dzt1db4okkw80

fmt.Println(rid.New().NumeralString())
// -> 218093205000013254195242051055202176
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