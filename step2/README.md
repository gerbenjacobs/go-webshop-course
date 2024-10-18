# Step 2: Let's add Products

In this step we're going to create a handler, service and storage for __Products__.

## Interfaces

We will use interfaces a lot in the svc framework, this allows us to make contracts with components,
but also easily test our code.


### product.go

Yep, no folders, we will put this one straight into the root of our project.

This is going to be our domain model for a __Product__.

```go
package go_webshop_course

type Product struct {
	ID          int
	Name        string
	Description string
	Image       string
	Price       float64
}
```

Floating point numbers aren't great to represent money, but hey, this is just a course!