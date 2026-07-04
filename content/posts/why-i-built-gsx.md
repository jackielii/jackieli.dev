+++
date = 2026-07-04T09:26:00Z
title = "Why I built gsx"
description = "Why gsx exists, and why I didn't just use templ"
slug = "why-i-built-gsx"
authors = []
tags = ["go", "templating", "gsx", "templ"]
categories = []
externalLink = ""
series = []
+++

I have been building [gsx](https://github.com/gsxhq/gsx), a Go template compiler
with JSX-style markup, `component` declarations, and plain Go output.

The obvious question is: why?

More specifically: why not just use [templ](https://templ.guide/)?

Firstly, I like templ. gsx is not built because I think templ is bad. In fact,
gsx is intentionally compatible with the templ ecosystem. A `gsx.Node` has the
same render shape as `templ.Component`:

```go
Render(ctx context.Context, w io.Writer) error
```

So this is not a "throw everything away and start again" situation.

The difference is the tradeoff.

templ tries quite hard to stay simple. The syntax is close to Go, the compiler
model is relatively simple, and many things are explicit. That is a very
reasonable design.

But I think it pays for that simplicity with ergonomics.

## The recurring pressure

After using templ, I kept running into the same category of problems.

Not correctness problems necessarily. More like friction.

Things like:

- component calls do not read like HTML
- passing rich inline markup often needs extra named components
- class composition can become verbose
- JavaScript/data interpolation is awkward
- JSON data islands need helper code
- attributes want to behave more like markup attributes
- editor/dev-loop support needs to understand the template language deeply

What convinced me this was not just my personal taste is that templ has had many
issues and proposals in exactly these areas.

There are proposals for HTML-style component authoring:
#663 (https://github.com/a-h/templ/issues/663),
#1181 (https://github.com/a-h/templ/issues/1181).

There is a proposal for anonymous/inline templ functions:
#1150 (https://github.com/a-h/templ/issues/1150).

There are discussions around passing Go data into JavaScript:
#944 (https://github.com/a-h/templ/issues/944),
#838 (https://github.com/a-h/templ/issues/838).

There was a proposal for JSON helpers:
#739 (https://github.com/a-h/templ/issues/739).

There is even a newer discussion about whether interpolating Go data inside
<script> tags is worth the parser complexity:
#1408 (https://github.com/a-h/templ/issues/1408).

Class and attribute ergonomics show up too:
#61 (https://github.com/a-h/templ/issues/61),
#902 (https://github.com/a-h/templ/issues/902),
#933 (https://github.com/a-h/templ/issues/933).

And dev-loop/tooling pressure is also there:
#318 (https://github.com/a-h/templ/issues/318).

So I think the demand is real. People like Go-checked templates, but they also
want the authoring experience to feel more like writing HTML.

## Retrofitting is hard

The hard part is that these features are not isolated.

HTML-style component calls are not just syntax. Once you have:

<Card title="Hello" class="featured" />

the compiler needs to know whether Card is a component, what props it accepts,
how attributes map to fields, what type each expression has, and what context
each value is rendered into.

If you want good editor support, the language server needs the same knowledge.

If you want safe JavaScript interpolation, the compiler needs to understand
whether the value is in a JS value position, a string position, a regex position,
an attribute value, or a JSON data island.

If you want class merging, the compiler needs to treat class as a structured
value, not just a string attribute.

You can add these things one at a time, but at some point you are no longer
keeping the compiler simple. You are building a more ambitious toolchain.

That is the point where I thought: maybe the design should start there.

## The gsx tradeoff

gsx spends more complexity in the toolchain to make templates nicer to write.

It uses go/packages and go/types to scan and analyze real Go source. This is
not free. It is more machinery than a simpler parser/code generator.

But it buys useful things.

A component can look like this:

component Card(title string, featured bool) {
  <section class={ "card", "card-featured": featured }>
    <h2>{ title }</h2>
    { if featured { <span class="badge">Featured</span> } }
    <div>{ children }</div>
  </section>
}

The markup reads like markup.

The data is still Go.

The generated output is still plain Go.

Props are named Go fields, checked by go build.

## Class ergonomics

Class composition is one of those small things that matters a lot in real UI
code.

In gsx:

<span class={ "tag", "tag--active": active }>
  { label }
</span>

This renders as:

<span class="tag tag--active">stable</span>

There is also explicit attribute forwarding:

component Button(variant string) {
  <button class="btn" data-variant={variant} { attrs... }>
    { children }
  </button>
}

component Page() {
  <Button variant="primary" class="w-full" data-test="x" hx-post="/go">
    Save
  </Button>
}

The rendered class is merged:

<button class="btn w-full" data-variant="primary" data-test="x" hx-post="/go">
  Save
</button>

This is a very boring feature, but boring features are where UI ergonomics live.

## JavaScript, but not the whole circus

I do not hate JavaScript.

I have used enough JavaScript and Node.js to know the bad parts, but I also do
not want to pretend the good parts are bad.

JSX ergonomics, fast reloads, editor tooling, browser error overlays, and tight
feedback loops are good ideas.

I just do not want the whole circus.

In gsx, JavaScript-valued attributes are explicit:

<button @click=js`openMenu()`>Open</button>

For Alpine-style attributes:

<div x-data=js`{ open: false }`>
  <button @click=js`open = !open`>Toggle</button>
  <div x-show=js`open` @click.outside=js`open = false`}>
    Contents...
  </div>
</div>

For JSON-ish attributes like hx-vals, Go values in @{ ... } holes are
serialized as JSON automatically:

component EntityFilter(entityType string, opts map[string]string) {
  <input
    type="checkbox"
    hx-post="/filter"
    hx-vals=js`{"entity_type": @{entityType}, "opts": @{opts}}`
  />
}

And for data islands:

component Widget(cfg Config) {
  <div>
    <button @click=js`toggle()`>Toggle</button>
    <script type="application/json" id="cfg">@{ cfg }</script>
  </div>
}

This is the kind of thing I wanted: Go values, HTML-shaped templates, explicit
JavaScript contexts, and generated output I can inspect.

## Context matters

A string is not just a string when rendering HTML.

Text content, attribute values, URLs, CSS, JavaScript strings, JavaScript values,
and JSON data all have different escaping rules.

gsx treats the position as part of the meaning.

For example:

component Link(u string) {
  <a href={ u |> trim }>x</a>
}

If u is:

"  javascript:alert(1)  "

the rendered output is:

<a href="about:invalid#gsx">x</a>

That is the sort of thing I want the compiler to own. Not every call site. Not a
helper remembered by convention. The template compiler can see the context, so
it should use that context.

## Tooling is part of the language

Another reason I wanted a deeper toolchain is editor support.

gsx already has:

- gsx init
- gsx dev
- gsx generate
- gsx fmt
- gsx lsp
- a VS Code extension
- a tree-sitter grammar
- a browser playground
- Vite integration

gsx dev runs the development loop in one foreground process: warm generation,
Go server rebuilds, Vite, browser reloads, and error overlays.

The language server provides diagnostics, hover, go-to-definition, references,
and formatting.

The tree-sitter grammar gives highlighting across Go, HTML, JavaScript, and CSS
inside .gsx files.

This is very much inspired by the better parts of the JavaScript development
experience. I like tight feedback. I like editor tooling. I like errors showing
up where I am working.

I just want the end result to still be Go.

## What gsx is not

gsx is not a router.

It is not an HTTP framework.

It is not an app structure.

It is just templating.

That boundary is important. Everything outside templates should stay ordinary
Go. Use net/http, chi, echo, htmx, structpages, whatever you want.

gsx should only care about turning .gsx into checked Go code that renders
HTML.

## In summary

templ made a good tradeoff: keep the model simple and close to Go.

gsx makes a different tradeoff: do more analysis in the toolchain so templates
can be more ergonomic.

That means:

- HTML-shaped component calls
- named props checked by Go
- contextual escaping
- class and attribute composition
- explicit JavaScript contexts
- JSON/data-island ergonomics
- editor-aware diagnostics and navigation
- plain Go output

It is alpha software, and the language will probably still change.

But the direction is clear: HTML as a first-class Go value, with enough tooling
behind it to make the authoring experience feel good.


I did not write it into `/Users/jackieli/personal/jackieli.dev/content/posts` because that repo is outside the current writable sandbox. The path I found
is `/Users/jackieli/personal/jackieli.dev/content/posts`; your original `../../jackieli.dev/content/posts/` did not resolve from `/Users/jackieli/
personal/gsxhq`.
