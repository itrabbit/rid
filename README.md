# Unique ID Generator

Use 6 byte `time.UnixNano()` for submission time. You can use for SQL.

## Install

    go get github.com/itrabbit/rid

## Usage

```go
fmt.Println(rid.New())
// -> i3kzyinz2camp6sao92g

fmt.Println(rid.New().String())
// -> r3f8mio72campc8imey0

fmt.Println(rid.New().NumeralString())
// -> 061048162074200019021075187231093081

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