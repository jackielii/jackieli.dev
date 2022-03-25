+++ 
date = 2020-07-24T11:00:17+01:00
title = "BigInt toBytes in different languages"
description = "Difference between toBytes & fromBytes on BigInt"
slug = "bigint-tobytes" 
tags = ["bigint", "go", "java", "kotlin", "javascript", "typescript"]
categories = []
externalLink = ""
series = ["bigint"]
+++

Go & Javascript seems to be doing the same thing on BigInt: toBytes returns the
absolute value with the sign seperately. While Java returns the [Two's
complement](https://en.wikipedia.org/wiki/Two%27s_complement) representation.
It makes it just a bit trickier to get the right binary representation.

[This post]({{< relref "./js-bigint-uint8array.md" >}}) explores ways to
implement the Two's complement representation in javascript & Go

### javascript

```sh
$ node
> BigInt(-257).toString(16)
'-101'
> BigInt(257).toString(16)
'101'
>
```

It's basically the sign plus the hex of the absolute value

### Go

```go
package main

import (
	"fmt"
	"math/big"
)

func main() {
	fmt.Printf("%#v", big.NewInt(-257).Bytes())
}

// output
// []byte{0x1, 0x1}
```

From `big.Int.Bytes()` [documentation](https://pkg.go.dev/math/big?tab=doc#Int.Bytes): 

```txt
Bytes returns the absolute value of x as a big-endian byte slice.
```

### java / kotlin

```kotlin
println(Arrays.toString(BigDecimal("-257").unscaledValue().toByteArray()))

// output
// [-2, -1]
```

The doc of `unscaledValue` says:

```txt
Returns a byte array containing the two's-complement representation of this
BigInteger. The byte array will be in big-endian byte-order: the most
significant byte is in the zeroth element. The array will contain the minimum
number of bytes required to represent this BigInteger, including at least one
sign bit, which is (ceil((this.bitLength() + 1)/8)). (This representation is
compatible with the (byte[]) constructor.)
```
