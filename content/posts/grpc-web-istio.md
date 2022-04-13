+++
date = 2022-03-25T15:55:25Z
title = "GRPC & GRPC-Web multiplexed in Istio"
description = "Expose GRPC & GRPC-Web via One Port in Istio"
slug = "grpc-web-istio"
tags = ["istio", "kubernetes"]
categories = ["gRPC-Web"]
externalLink = ""
series = []
+++

Envoy always had grpc-web support very ealy on. It's used in the official
grpc-web [tutorial
docs](https://github.com/grpc/grpc-web#2-run-the-server-and-proxy). The [envoy
config](https://github.com/grpc/grpc-web/blob/8c5502186445e35002697f4bd8d1b820abdbed5d/net/grpc/gateway/examples/echo/envoy.yaml)
supports both grpc & grpcWeb via one `listener` on on port out of the box. This
is important when you have a production deployment in which some clients
(mobile apps) want to speak to the GRPC protocol directly, but web apps want to
go down to http/1.1 to work with browsers' xhr or fetch requests.

In production, we use istio as our gateway, but we have always struggled to
utilise one single port to proxy both GRPC and GRPC-Web requests, even though
envoy - the underlying proxy supports it very well. I suppose it's partially
due to the lack of documentation out there and the GRPC-Web community is still
relatively small.

As we adopted Istio & GRPC-Web quite early on, we used what some other bloggers
suggested: using a sidecar to convert HTTP1.1 request to GRPC and use a
separate gateway dedicated to GRPC. The solution roughly looks like this:

**For grpc**

`gateway:31400` -> `app`

**For grpc-web**

`gateway:443` -> `pod sidecar` -> `app`

This has a few issues:

1. We're adding an extra open port to outside world: more attacking surface
2. The new grpc-only port doesn't have TLS or needs a separate TLS config
3. We added a sidecar pod for every grpc service that we want to expose as
   GRPC-Web: wasted resources when we don't need service mesh

This is definitely not right when evnoy supports the use case so well, but why
does it feels so difficult on Istio's side?

I finally sat down today and look at this from scratch: We have an effective
envoy config, so we just need to configure istio such that the gateway envoy of
istio has the most important elements that supports both GRPC & GRPC-Web.
Namely:

1. on the listener config we need the filter [envoy.filters.http.grpc_web](https://github.com/grpc/grpc-web/blob/8c5502186445e35002697f4bd8d1b820abdbed5d/net/grpc/gateway/examples/echo/envoy.yaml#L38)
2. on the clusters config we need the [http2_protocol_options](https://github.com/grpc/grpc-web/blob/8c5502186445e35002697f4bd8d1b820abdbed5d/net/grpc/gateway/examples/echo/envoy.yaml#L45)

And we definitely don't need the sidecar for this: sidecar is for east-west
traffic, not north-south traffic. So the filter needs to be on the gateway, not
the sidecar.

Well, it turns out it's still in the same
[EnvoyFilter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)
documentation. We just need to 1) configure the filter on the gateway; 2) use
`grpc-web` as the protocol in the service. (I believe `grpc` works just as
well)

## The solution

**envoy filter**

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: gateway-grpc-web-filter
  namespace: istio-system
spec:
  workloadSelector:
    labels:
      istio: ingressgateway
  configPatches:
    - applyTo: HTTP_FILTER
      match:
	  	# importantly we're patching to the GATEWAY envoy, not sidecar
        context: GATEWAY
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
				# apply the patch before the cors filter, just like the one in
				# grpc-web example
                name: "envoy.filters.http.cors"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.grpc_web
```

**virtual service**

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hello-grpc
  namespace: default
spec:
  hosts:
    - "*"
  gateways:
    - httpbin-gateway
  http:
    - route:
        - destination:
            host: hello-grpc
            port:
              number: 12345
      match:
        - uri:
            prefix: /main.HelloService/
```

And importantly the service needs to have `grpc` in the service name or use
`appProtocol`:
[ref](https://istio.io/latest/docs/ops/configuration/traffic-management/protocol-selection/#explicit-protocol-selection)

Full example on [github](https://github.com/jackieli-tes/learn-grpc-web-istio).
There are some useful debug commands that helped.
