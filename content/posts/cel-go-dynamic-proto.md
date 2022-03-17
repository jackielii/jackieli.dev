+++
date = 2022-03-17T17:15:09Z
title = "Use dynamic proto in cel-go"
slug = "use-dynamic-proto-in-cel-go"
+++

[Cel-go](https://github.com/google/cel-go) is an amazing library for evaluating
expressions. The extensive proto support just makes it better.

Let's explore an interesting problem: if we were to build a evaluation service
without embedding the actual proto, i.e. without bundling the generated code
that defines the struct, fields etc, how can we still evaluate the expressions
that requests the contents of the proto messages?

Let's start with this proto:

```proto
syntax = "proto3";

package main;

option go_package = ".;main";

message Foo {
  string foo = 1;
}

```

```go
func messageBytes() []byte {
	foo := &Foo{Foo: "foo message"}
	b, err := proto.Marshal(foo)
	if err != nil {
		panic(err)
	}
	return b
}

func descBytes() []byte {
	set := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			protodesc.ToFileDescriptorProto(File_foo_proto),
		},
	}
	b, err := proto.Marshal(set)
	if err != nil {
		panic(err)
	}
	return b
}
```

Now if we know the bytes of `Foo` message using `messageBytes` and the proto
descriptor also in bytes (every generated code has it). Then we declare it as
var `x` and evaluate `x.foo`. And we should be able to get value `foo message`
right?

Turns out it's all supported within cel-go. We just register the all the file
descriptor related to message `foo` and use dynamic message to wrap the bytes:

```go
func exercise9() {
	// step 1: register the proto descriptors
	fileSet := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(descBytes(), fileSet)
	if err != nil { panic(err) }

	fooFullName := "main.Foo" // `main` is the proto package name
	reg, err := protodesc.NewFiles(fileSet)
	if err != nil { panic(err) }

	// step 2: wrap it with dynamicpb message
	desc, err := reg.FindDescriptorByName(protoreflect.FullName(fooFullName))
	if err != nil { panic(err) }
	msg := dynamicpb.NewMessage(desc.(protoreflect.MessageDescriptor))
	err = proto.Unmarshal(messageBytes(), msg.Interface())
	if err != nil { panic(err) }

	// step 3: cel magic
	env, _ := cel.NewEnv(
		cel.TypeDescs(fileSet),
		cel.Declarations(decls.NewVar("x", decls.NewObjectType(fooFullName))),
	)
	ast, iss := env.Compile(`x.foo`)
	if iss.Err() != nil {
		glog.Exit(iss.Err())
	}
	// Turn on optimization.
	vars := map[string]interface{}{"x": msg}
	program, _ := env.Program(ast, cel.EvalOptions(cel.OptExhaustiveEval))
	eval(program, vars)
	
	// output: foo message
}
```

Full source at [github](https://github.com/jackieli-tes/learn-cel-go)
