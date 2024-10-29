# Step 3: ???

## Product detail page

Let's see if we can create a detailed product page.
We want to have a URL like `/product/:id` that will fetch data for the product with that ID,
and then show some nice HTML regarding it.

To do this, we'll have to make some changes in all the layers:
- __static__: We need to create a product page in _HTML_
- __handler__: Our product page requires its own _dedicated route_
- __service__: We don't know how to handle single products yet
- __storage__: .. this applies to our storage layer too

### static/product/product.html

Let's create a folder called `product` in our `/static/` directory and then create a file
called `product.html`.

```html
{{ define "title" }}Gopher plushie{{ end }}

{{ define "content" }}
<div class="row">
    <div class="col-6 m-auto">
        <div class="card">
            <img src="https://picsum.photos/600/300" class="card-img-top" alt="Product description">
            <div class="card-body">
                <h5 class="card-title">Gopher plushie</h5>
                <h6 class="card-subtitle mb-2 text-body-secondary">€ 12,99</h6>
                <p class="card-text">A small purple Gophier plushie, perfect for kids and adults alike.</p>
                <a href="#" class="btn btn-primary">Add to cart</a>
            </div>
        </div>

    </div>
</div>
{{ end }}
```

For now we will hardcode some data, put in a [card component](https://getbootstrap.com/docs/5.3/components/card/) 
and make a `row` with a `col-6` (i.e. half the screensize) and use `m-auto` to center this. 

### handler/handler.go

We need to add our route. Notice how we use a colon to denote a parameterized argument.

```go
r.GET("/product/:id", h.productByID)
```

### handler/product.go

```go
func (h *Handler) productByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productID := p.ByName("id")
	h.logger.DebugContext(r.Context(), "Request received",
		"method", r.Method,
		"url", r.RequestURI,
		"product_id", productID,
	)
	tmpl := template.Must(template.ParseFiles(
		"static/layout.html",
		"static/product/product.html",
	))

	// fetch our product
	// TODO: We have to implement this still..

	// set up our page data
	type pageData struct {
		User bool
	}
	data := pageData{
		User: false,
	}

	// render the templates
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}
```

For the first time we're doing something with the `httprouter.Params` argument,
as we get the product ID from the request by calling `p.ByName("id")`.

We add this product ID now to our logger as well, since it's nice information to have.

Something that's different from our `/` handler is that we now require the templates
for `layout.html` and `product/product.html`.

When we eventually have our service in order, we could 'fetch our product', but since we don't,
we're going have to leave that for now and come back later.

Then we use another `pageData` struct, this time we only really need `User` because of the
code we have in `layout.html`. Our `product.html` doesn't have any data requirements (yet).

### Let's have a quick look.

Stop your application in case it's running and restart it again (make sure you're in the `step3` folder):

```shell
go run cmd/app/main.go
```

Open your browsers and navigate to http://localhost:8000/product/1.

Now navigate to http://localhost:8000/product/gopher-plushie. 

As you can see, we still need to do some work. Validation, fetching and proper rendering!

### storage/storage.go

Back in our interfaces location, we currently have the following.

```go
type ProductRepository interface {
    GetAllProducts(context.Context) ([]app.Product, error)
}
```

Can you add another method for `GetProduct` that aside from a `context.Context` also takes a product ID.
Have a look at `app.Product` how you store your `ID`. Then make sure your method returns an `app.Product` and an `error`.

Now it's time to implement this in our in-memory repository.

### storage/product_repository_memory.go

Copy the method signature you just created in the interface, or use your IDE to auto-complete the missing
interface methods.

__Task__: Implement the method.

Alright, so I'm assuming you also have a case for when the product doesn't exist, and you most likely return an error.

Let's create a new custom error, at the top of our file.

```go
var ErrProductNotFound = errors.New("product not found")
```

In your code where you deal with the error, return this error instead.

You can also use [error wrapping](https://rollbar.com/blog/golang-wrap-and-unwrap-error/) to add more information
but keep the original error value intact.

```go
return app.Product{}, fmt.Errorf("%w: we don't know about id: %d", ErrProductNotFound, id)
```

This will return `product not found: we don't know about id: wrong-id-entered` which is its own error,
but with help from `errors.Unwrap()` and `errors.Is()` we can distinguish this as an `ErrProductNotFound` error.