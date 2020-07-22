+++
date = "2020-08-23"
title = "BigInt to Uint8Array"
description = "How to convert BigInt to Uint8Array in javascript"
slug = "bigint to uint8array"
series = ["gRPC", "gRPC-Web", "binary"]
+++

The problem came when I wanted to pass down BigDecimal's unscaledBytes from kotlin/java to javascript through gRPC-Web.

Java's implementation of the unscaledBytes returns signed big-endian bytes

Javascript has BigInt support, but lacks of `toBytes()`, stuff around the internet all seems to handle positive numbers just fine, but fails in negative number implementation

I borrowed from a Go implementation and adapted to javascript

## Javascript implementation (using Typescript)



