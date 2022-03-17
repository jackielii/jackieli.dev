+++
draft = true
date = 2022-03-17T12:57:13Z
title = "My thoughts on OSDU"
slug = "thoughts-on-osdu"
+++

IMHO, "what level of abstraction" is the question:

Are we talking about data model standards? PPDM solves this

Are we talking about micro-service infrastructure? Service oriented architecture? Authn & Authz? AWS, Google, Azure and the likes solve this

Abstract application level API? There are certain things that can be abstracted well (metadata tags, resource labels etc), but the majority is so different that many have attempted and failed. And in some of these cases, the differences are needed: competing with each other for the better. So can we and should we?

With all that said, the open binary data formats are invaluable: OpenZGY, OpenVDS etc. Not having to guess the usage, reverse engineering the effects of the parameters are a huge time saver. Not to mention the knowledge captured in these data spec and source code.

"Open architecture" or "open source" is proven a vaiable business model and I love open source. But coming up with an complete system by a community end to end is too far fetched. By then OSDU would be yet another system to integrate, albeit simpler.
