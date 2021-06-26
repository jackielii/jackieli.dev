+++ 
date = 2020-08-01T11:28:51+01:00
title = "gRPC-Web interceptor"
description = "a journey of adding a simple gRPC-Web interceptor"
slug = "grpc-web-interceptor" 
tags = ["javascript", "typescript"]
categories = ["gRPC-Web"]
externalLink = ""
series = []
+++

This was supposed to be 5 minute job, or at least it looked like it after I skimmed through the [documentation](https://grpc.io/blog/grpc-web-interceptor/)

This was my first attempt: 

```typescript
class AuthInterceptor<REQ extends Request, RESP = any> {
  Intercept = class {
    private stream: grpcWeb.ClientReadableStream<RESP>;

    constructor(stream: grpcWeb.ClientReadableStream<RESP>) {
      this.stream = stream;
    }

    on<F extends Function>(eventType: EventType, callback: F) {
      if (eventType === "error") {
        this.stream.on("error", (err: grpcWeb.Error) => {
          if (process.env.NODE_ENV === "development") {
            console.log("grpc web on error", err);
          }
          if (err.code === grpcWeb.StatusCode.UNAUTHENTICATED) {
            OauthHelper.redirectToSSO();
          }
          callback(err);
        });
      } else if (eventType === "data") {
        this.stream.on("data", (resp) => {
          if (process.env.NODE_ENV === "development") {
            console.log("grpc web response", (resp as any)?.toObject());
          }
          callback(resp);
        });
      } else if (eventType === "status") {
        this.stream.on("status", (status) => {
          if (process.env.NODE_ENV === "development") {
            console.log("grpc web status", status);
          }
          callback(status);
        });
      } else if (eventType === "end") {
        this.stream.on("end", callback as any);
      }
      return this;
    }

    cancel() {
      if (process.env.NODE_ENV === "development") {
        console.log("grpc web cancelled");
      }
      this.stream.cancel();
      return this;
    }
  };

  intercept(
    request: REQ,
    invoker: (
      request: REQ,
      metadata?: grpcWeb.Metadata
    ) => grpcWeb.ClientReadableStream<RESP>
  ) {
    const md = request.getMetadata();
    md["Authorization"] = `Bearer ${getAuthToken()}`;
    // cancellation
    const signal = md.signal;
    delete md.signal;
    if (process.env.NODE_ENV === "development") {
      console.log(
        "grpc-web request:",
        request.getRequestMessage()?.toObject(),
        "metadata:",
        md
      );
    }
    const stream = invoker(request);
    const newStream = new this.Intercept(stream);
    if (signal) {
      signal.addEventListener("abort", () => newStream.cancel());
    }
    return newStream;
  }
}
```

So pretty straight forward: I intercept request to add Authorization header &
intercept response to log them. 

And also I added a `signal` which is
[AbortSignal](https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal) to
handle cancellation: when making request, e.g. in createAsyncThunk from
redux-toolkit, I could just pass the signal from there. It all worked pretty
well.

_Then I did a production build, opened in incognito mode, and it just don't
work!_

Because I have service worker enabled in the prod build, so usually I have to
open a incognito window to see it - just something people do - and it just
doesn't work! The interceptor just didn't run at all: auth header is not added,
signal property is not removed before sending to request.

What's worse is somehow the production build works in normal browser window!

You can imagine at this stage I'm starting to question life, and it was 2am...

After wasting my time trying to debug and trace down the issue in the grpcWeb &
generated libraries, I found part of the problem: I was using a Promise client
and it needs a UnaryInterceptor, like so:

```ts
class UnaryAuthInterceptor<REQ extends Request, RESP extends UnaryResponse> {
  async intercept(request: REQ, invoker: (request: REQ) => Promise<RESP>) {
    const md = request.getMetadata()
    md['Authorization'] = `Bearer ${getAuthToken()}`
    if (process.env.NODE_ENV === 'development') {
      console.log('grpc-web request:', request.getRequestMessage()?.toObject(), 'metadata:', md)
    }
    // cancellation
    // const signal = md.signal // UnaryCall doesn't allow cancellation
    delete md.signal
    try {
      const resp = await invoker(request)
      if (process.env.NODE_ENV === 'development') {
        console.log('grpc-web unary response:', resp.getResponseMessage())
      }
      return resp
    } catch (e) {
      if (process.env.NODE_ENV === 'development') {
        console.log('grpc-web unary error', e)
      }
      throw e
    }
  }
}
```

Basically the same thing but needs a Promise and no cancellation.

This actually worked, but still it doesn't explain why it worked in normal browser window & not incognito mode. But at least at this stage, I'm able to have a reliable interceptor, albeit less functionality

Should I call it good enough? Probably. But it really bothered me it works in normal browser window but not the incognito mode.

I did make some progress through: I can reproduce it even in dev build in incognito mode. So the only difference between these two is ... the plugins!

I have gRPC-Web devtools plugin installed. Surely it shouldn't be the problem? Turned it off and I can reproduce the problem even in normal browser mode.

Finally: so somehow with gRPC-Web devtools, a Promise client will be turned into a RPC client that works with ClientReadableStream which is what StreamInterceptor works with. 

By using gRPC-Web devtools, it gave me a false positive that shield the problem that I was using the wrong Interceptor.

So what did I learn? **Read the documentation closely** 

And since it was 3am already. I just made my own PromiseClient from the RPC Client, this way I get the cancellation as well:

```typescript
const createPromiseClient = <PromiseClient, RpcClient = unknown>(
  client: RpcClient,
): PromiseClient => {
  const methods = Object.getPrototypeOf(client)
  return Object.keys(methods).reduce((acc, method) => {
    const rpc = methods[method].bind(client)
    acc[method] = (request: any, metadata?: grpcWeb.Metadata) =>
      new Promise((resolve, reject) => {
        rpc(request, metadata, (err: grpcWeb.Error, resp: any) => {
          if (err) {
            reject(err)
            return
          }
          resolve(resp)
        })
      })
    return acc
  }, {} as any)
}
```

[Discussions](https://github.com/jackielii/jackieli.dev/discussions)
