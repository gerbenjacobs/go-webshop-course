# Golang Webshop Course

In this course we build a simple webshop with a REST API and a frontend page.

We will use some of [the learnings](https://blog.gerbenjacobs.nl/svc-an-opinionated-go-service-framework/) 
from the [svc framework](https://github.com/gerbenjacobs/svc), have a working shopping cart 
and use [JSON Web Tokens (JWT)](https://jwt.io/) for authentication.

For our frontend we'll use [Go HTML templating](https://gowebexamples.com/templates/), 
but since we'll have an API, a standalone frontend can be built as well.

So what are you gonna sell in your shop?


## Steps

This course is divided into steps, each step has a `README.md`.

The README of each step contains information about what we're doing, snippets of codes you can copy
and instructions (prefixed with __Task__) that you have to complete yourself.

It also contains all the code that's meant as a **reference**.

Since you're most likely to have your own project, you can just open this project on GitHub
and use its Markdown generator to read the READMEs.


### Step List

The steps relate to a common goal, but as such might have different lengths of time for completion.

**Step 1: Set up a simple Go web app with HTML.**

We create a Go app that returns some HTML, we use Bootstrap to create a layout, deal with routing and 
learn some basic concepts from the svc framework.

**Step 2: Let's add Products**

We have a working website with Go, but nothing is dynamic. We'll create storages and services,
learn about domain models, translation helpers and interfaces. We setup the full stack
from a Product storage to a Product handler, and display products on our page.

**Step 3: Deeper into the web**

Time to go deeper into the web. We create a 'show product' page, but this requires
changes in our routes, services and storage. We have to start dealing with errors, both domain models
and HTTP status errors. We need to deal more with Go's HTML templating and we introduce 'Flashes'.

**Step 4: A JSON REST API**

Wowee, the web is tiring. Let's start introducing a REST-ish API that replies with JSON. 
We'll learn how we can reuse our services and storage layer, but introduce a new set of handlers
and routes that help us with creating an API. Maybe someone will even make a web interface for us ;)

## Why not X?

There's a myriad of options when it comes to web frameworks, even more so in Go,
since there's really no clear winner unlike PHP, Python or Ruby.

We're doing a lot of work ourselves in this Go course, but that's also the point.

We're trying to show how things could be done, as abstraction-less as possible.