+++ 
date = 2020-07-23T01:05:39+01:00
title = "BigInt to Uint8Array"
description = "How to convert BigInt to Uint8Array in javascript"
slug = "bigint to uint8array"
series = ["gRPC", "gRPC-Web", "binary"]
tags = ["binary", "javascript", "typescript", "bigint", "uint8array"]
categories = ["javascript"]
externalLink = ""
+++

The problem came when I wanted to pass down BigDecimal's unscaledBytes from kotlin/java to javascript through gRPC-Web.

Java's implementation of the unscaledBytes returns signed big-endian bytes

Javascript has BigInt support, but lacks of `toBytes()`, stuff around the internet all seems to handle positive numbers just fine, but fails in negative number implementation

I borrowed from a Go implementation and adapted to javascript

## Javascript implementation (using Typescript)

```
const TWO_POWER_31 = 2147483648

const big0 = BigInt(0)
const big1 = BigInt(1)
const big8 = BigInt(8)

export class BigDecimal {
  private big?: bigint
  private num?: number
  private isNumber: boolean
  scale: number

  constructor(int: bigint | number, scale?: number) {
    this.isNumber = int < TWO_POWER_31 && int >= -TWO_POWER_31
    if (!this.isNumber) {
      this.big = BigInt(int)
    } else {
      const n = Number(int)
      if (Math.floor(n) !== n) {
        throw new Error('big decimal only accepts integer and scale')
      }
      this.num = n
    }
    this.scale = scale || 0
  }

  bigToUint8Array() {
    let big: bigint = this.big!!
    if (big < big0) {
      const bits: bigint = (BigInt(big.toString(2).length) / big8 + big1) * big8
      const prefix1: bigint = big1 << bits
      big += prefix1
    }
    let hex = big.toString(16)
    if (hex.length % 2) {
      hex = '0' + hex
    }
    const len = hex.length / 2
    const u8 = new Uint8Array(len)
    var i = 0
    var j = 0
    while (i < len) {
      u8[i] = parseInt(hex.slice(j, j + 2), 16)
      i += 1
      j += 2
    }
    return u8
  }

  numToUint8Array(): Uint8Array {
    let n = this.num!!
    const arr: number[] = []
    while (true) {
      if (!n) {
        break
      }
      arr.unshift(n & 0xff)
      n >>>= 8
    }
    if (arr.length === 0) {
      return new Uint8Array([0])
    }
    return new Uint8Array(arr)
  }

  unscaledBytes(): Uint8Array {
    if (this.isNumber) {
      if (this.num!! < 0) {
        this.big = BigInt(this.num)
        // TODO: use big for more concise representation, it eliminates negative padding
      } else {
        return this.numToUint8Array()
      }
    }
    return this.bigToUint8Array()
  }

  getScale(): number {
    return this.scale
  }

  toString() {
    let s: string
    if (this.isNumber) {
      s = this.num!!.toString()
    } else {
      s = this.big!!.toString()
    }
    if (this.scale == 0) {
      return s
    }
    if (this.scale > s.length) {
      throw new Error(`scale ${this.scale} shouldn't be bigger than length: ${s.length}`)
    }
    const dot = s.length - this.scale
    if (dot === 0) {
      return '0.'.concat(s)
    }
    return s.slice(0, dot).concat('.').concat(s.slice(dot))
  }

  static fromBytes(a: Uint8Array, scale?: number): BigDecimal {
    if (!a.length) {
      return new BigDecimal(0)
    }
    const hex = Buffer.from(a).toString('hex')
    let big = BigInt('0x' + hex)
    if (a[0] & 0x80) {
      const negative = BigInt('0x1' + '0'.repeat(hex.length))
      big -= negative
    }
    return new BigDecimal(big, scale)
  }

  static fromString(s: string): BigDecimal {
    const dot = s.indexOf('.')
    const minus = s.indexOf('-')
    if (dot === -1) {
      return new BigDecimal(BigInt(s), 0)
    }
    if (dot === s.length - 1) {
      s = s.slice(0, s.length - 1)
    }
    // .-1234
    if (dot == 0 && minus == 1) {
      throw new Error('invalid big decimal number'.concat(s))
    }
    s = s.slice(0, dot).concat(s.slice(dot + 1))
    // -.1234 = -0.1234
    if (dot == 1 && minus == 0) {
      // TODO
    }
    return new BigDecimal(BigInt(s), s.length - dot)
  }
}
```

## Go implementation

```
var one = big.NewInt(1)

// SetSignedBytes sets the value of n to the big-endian two's complement
// value stored in the given data. If data[0]&80 != 0, the number
// is negative. If data is empty, the result will be 0.
func SetSignedBytes(n *big.Int, data []byte) {
	n.SetBytes(data)
	if len(data) > 0 && data[0]&0x80 > 0 {
		n.Sub(n, new(big.Int).Lsh(one, uint(len(data))*8))
	}
}

// SignedBytes returns the big-endian two's complement
// form of n.
func SignedBytes(n *big.Int) []byte {
	switch n.Sign() {
	case 0:
		return []byte{0}
	case 1:
		b := n.Bytes()
		if b[0]&0x80 > 0 {
			b = append([]byte{0}, b...)
		}
		return b
	case -1:
		length := uint(n.BitLen()/8+1) * 8
		b := new(big.Int).Add(n, new(big.Int).Lsh(one, length)).Bytes()
		// When the most significant bit is on a byte
		// boundary, we can get some extra significant
		// bits, so strip them off when that happens.
		if len(b) >= 2 && b[0] == 0xff && b[1]&0x80 != 0 {
			b = b[1:]
		}
		return b
	}
	panic("unreachable")
}
```
