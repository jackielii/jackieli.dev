+++
draft = true
date = 2022-06-21T14:14:58+01:00
title = "Golang Range Loop Capture"
description = "Golang Range Loop Capture"
slug = "golang-range-loop-capture"
tags = ["golang", "Go"]
+++

It's commonly known in Go, using an range loop with closure, the closure
captures the value by reference. E.g.

```go
for i := range []int{1, 2, 3} {
    func() {
		fmt.Println(i)
	}()
}
```
