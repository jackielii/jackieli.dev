+++ 
date = 2022-03-17T17:15:09Z
title = "Use dynamic proto in cel-go"
slug = "use-dynamic-proto-in-cel-go" 
+++

[Cel-go](https://github.com/google/cel-go) is an amazing library for evaluating
expressions. The extensive proto support just makes it better.

Let's explore an interesting question: if we were to build a evaluation service
without knowing the actual proto definitions, i.e. without bundling the
generated code, how can we still evaluate the expressions that requests the
fields of the proto? That is if we only have the bytes of the proto messages,
but don't have generated decoder, how do we evaluate.

Let's start with this example:

`proto`:

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
	foo := &Foo{Foo: "foo"}
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

Now if we receive the bytes of `messageBytes` and declare it as var `x` and
evaluate `x.foo`, how do we do it? Turns out it's all supported within cel-go.
We just register the all the file descriptor related to message `foo` and use
dynamic message to wrap the bytes:

```go
func exercise9() {
	// step 1: register the proto descriptors
	fileSet := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(descBytes(), fileSet)
	if err != nil { panic(err) }

	fooFullName := "main.Foo" // `main` is the proto package name
	reg, err := protodesc.NewFiles(fileSet)
	if err != nil { panic(err) }
	desc, err := reg.FindDescriptorByName(protoreflect.FullName(fooFullName))
	if err != nil { panic(err) }

	// step 2: wrap it with dynamicpb message
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
}
```

Full source at [github](https://github.com/jackieli-tes/learn-cel-go)
