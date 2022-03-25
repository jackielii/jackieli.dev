+++ 
date = 2021-06-26T12:40:10+01:00
title = "Pointers in Go - used in sql.Scanner"
description = "Golang sql.Scan without sql.NullString"
slug = "pointers-in-go-used-in-sql-scanner" 
tags = ["go","pointers","sql"]
categories = []
+++

Firstly a fun fact a lot people don't know: when passing values between
functions, it's cheaper to pass values instead of pointers in Go. Reason is
pointers to objects _could_ be allocated on heap and it takes the computer more
efforts to managed heap memory. While values are save on the stack, and stack
is cheaper. For example, it's preferable to:

```go
type Foo struct {
	Bar string
}

func (f Foo) Echo() {
	println(f.Bar)
}
```

In above example, `Echo` method will not change the value of Foo, so there is
no need to have the method receiver as a pointer.

Obviously Go compiler has gotten cleverer over the years and will do a lot more
escape analysis to determine if a pointer to struct can be allocated only on
stack. But I think it's still a good practice to use value instead of pointers
whenever you can.

Now here is a subtle bug that if used the value receive pattern incorrectly,
can be tricky to track down. It certainly cost me quite some time a while ago
and still bites me from time to time. Consider this example:

```
package main

type Bar struct {
	Baz string
}

func (b *Bar) Update() {
	b.Baz = "hello bar"
}

type Foo struct {
	Bar Bar
}

func (f Foo) UpdateBar() {
	f.Bar.Update()
}

func main() {
	bar := Bar{Baz: "doom"}
	foo := Foo{
		Bar: bar,
	}
	foo.UpdateBar()
	println(bar.Baz)
}
```

It'll print `doom` not `hello bar`. It's easy to see why in this simple
program, but when it's nested in structs, it could be a very long debugging
process.

{{< notice note >}}
If you embed a sync.Mutex in a struct and you attempt to use a value receiver,
go vet will warn you about it. Try running `go tool vet help copylocks`
{{< /notice >}}

The main trick I want to share today probably seems very normal to any C or C++
programmers. But it is very useful in practice.

### The problem

I have a proto generated struct or some struct that I don't own and I need to
use it directly in sql scan. i.e. given the following struct:

```go
type Node struct {
	ID 			int
	Name 		string
	CreatedAt 	timestamppb.Timestamp
}
```

However in database, I have `Name` and `CreatedAt` nullable, and I
don't really want to create a new struct or new variable using `sql.NullString`
just to use in `Row.Scan`.

### My solution

```go
type nullString struct {
	s *string
}

func (ts nullString) Scan(value interface{}) error {
	if value == nil {
		*ts.s = "" // nil to empty
		return nil
	}
	switch t := value.(type) {
	case string:
		*ts.s = t
	default:
		return fmt.Errorf("expect string in sql scan, got: %T", value)
	}
	return nil
}

func (n *nullString) Value() (driver.Value, error) {
	if n.s == nil {
		return "", nil
	}
	return *n.s, nil
}
```

And use it like this:

```go
var node Node

db.QueryRow("select name from node where id=?", id)
	.Scan(nullString(&node.Name))
```

Note here I don't need the pointer to the nullString, and I'm setting the value
the pointer is pointing to directly. And I didn't need to create a
`sql.NullString`.

And similarly, I can also create a `*timestamppb.Timestamp` scanner:

```go
import (
	"database/sql/driver"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type tsScanner struct {
	pt **timestamppb.Timestamp
}

// Scan implements sql.Scanner for protobuf Timestamp.
func (ts tsScanner) Scan(value interface{}) error {
	if value == nil {
		*ts.pt = nil
		return nil
	}
	switch t := value.(type) {
	case time.Time:
		tspb := timestamppb.New(t)
		*ts.pt = tspb
	default:
		return fmt.Errorf("expect time.Time in sql scan got: %T", value)
	}
	return nil
}

// Value implements driver.Valuer for protobukkf Timestamp.
func (ts tsScanner) Value() (driver.Value, error) {
	if ts.pt == nil || *ts.pt == nil {
		return nil, nil
	}
	return (*ts.pt).AsTime(), nil
}
```

And use it:

```go
var node Node

db.QueryRow("select created_at from node where id=?", id)
	.Scan(tsScanner(&node.CreatedAt))
```

In the `timestamppb.Timestamp` example, I'm using the pointer of pointer to the
struct because I need the set the pointer value, not the actual value.

It can make your head spin the whole pointer to pointer thing. Well, I'll
admit, it still does to me from time to time. Just need a bit more practice I
suppose?
