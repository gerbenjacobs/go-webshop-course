# Step 3: Deeper into the web

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

For now we will hardcode some data, put it in a [card component](https://getbootstrap.com/docs/5.3/components/card/) 
and make a `row` with a `col-6` (i.e. half the screensize) and use `m-auto` to center this. 
We haven't touched images yet either, so let's use a placeholder website to autogenerate us something.

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

Finally we use another `pageData` struct, this time we only really need `User` because of the
code we have in `layout.html`. Our `product.html` doesn't have any data requirements (yet).

### Let's have a quick look.

Stop your application in case it's running and restart it again (make sure you're in the `step3` folder):

```shell
go run cmd/app/main.go
```

Open your browsers and navigate to http://localhost:8000/product/1.

Now navigate to http://localhost:8000/product/gopher-plushie. 


```
Oct 29 19:45:40.989 INF Server started address=localhost:8000
Oct 29 19:45:45.350 DBG Request received method=GET url=/
Oct 29 19:45:48.161 DBG Request received method=GET url=/product/1 product_id=1
Oct 29 19:48:01.897 DBG Request received method=GET url=/product/gopher-plushie product_id=gopher-plushie
```

As you can see, we still need to do some work. Validation, fetching and proper rendering!

### storage/storage.go

Back in our interfaces location, we currently have the following.

```go
type ProductRepository interface {
    GetAllProducts(context.Context) ([]app.Product, error)
}
```

__Task__: Can you add another method for `GetProduct()` that aside from a `context.Context` also takes a product ID.
Have a look at `app.Product` how to store your `ID`. Then make sure your method returns an `app.Product` and an `error`.

Now it's time to implement this in our in-memory repository.

### storage/product_repository_memory.go

Copy the method signature you just created in the interface, or use your IDE to auto-complete the missing
interface methods.

__Task__: Implement the method.

Alright, so I'm assuming you also have a case for when the product doesn't exist, and you most likely return an error.

Let's create a new custom error, __but__ we will put this at the top of `/product.go`, our domain model.

```go
var ErrProductNotFound = errors.New("product not found")
```

Back in your repository code where you deal with the error, return this error instead. As per the svc guidelines, 
certain errors are also considered domain models. Something application specific such as 'product not found' 
is a good example, all layers can understand this.

You can also use [error wrapping](https://rollbar.com/blog/golang-wrap-and-unwrap-error/) to add more information
while keeping the original error intact.

```go
return app.Product{}, fmt.Errorf("%w: we don't know about id: %d", app.ErrProductNotFound, id)
```

This will return `product not found: we don't know about id: wrong-id-entered` which is its own error,
but with help from `errors.Unwrap()` and `errors.Is()` we can distinguish this as an `app.ErrProductNotFound` error.

The `%w` verb in the `fmt` package is a special character for 'wrapped errors'.

### services/services.go

Now we go to _services_.

First we have to update our service interface. We need something that can fulfill our `/product/:id` 
page with data. How about..

```go
type ProductService interface {
	ListProducts(context.Context) ([]app.Product, error)
	ShowProduct(context.Context, int) (app.Product, error)
}
```

It's very difficult to name things. We could have gone for `FetchProduct()` or even `GetProduct()`, just like 
our _storage_ interface.

In the case of services, we should be as specific as possible. We are creating this method specifically for 
'showing a detailed product'.

For example, if we ever have to deal with authorization and a user is not allowed to see a specific product, 
we could implement this in our `ShowProduct()` method. If you want a method to just always get a product, 
without any checks, maybe then we could introduce a `GetProduct()` method.

Both methods could still use `storage.GetProduct()` under the hood. And that's the beauty of having 
separate interfaces for services and storage.

### services/product.go

Now we need to implement `ShowProduct()`. 

Since our _service_ interface technically matches our _storage_ interface, we could try:

```go
func (p *ProductSvc) ShowProduct(ctx context.Context, productID int) (app.Product, error) {
	return p.repo.GetProduct(ctx, productID)
}
```

We just return the call from `repo` straight back, it's a `(app.Product, error)` return signature.

### handler/product.go

We can now complete the circle, and implement the `// TODO` in our product handler.

```go
func (h *Handler) productByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productIDParam := p.ByName("id")
	h.logger.DebugContext(r.Context(), "Request received",
		"method", r.Method,
		"url", r.RequestURI,
		"product_id", productIDParam,
	)

	// validate our product ID
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "couldn't convert product ID to int", "error", err)
		http.Error(w, "couldn't convert product ID to int", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"static/layout.html",
		"static/product/product.html",
	))

	// fetch our product
	product, err := h.Product.ShowProduct(r.Context(), productID)
	switch {
	case errors.Is(err, app.ErrProductNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		// an unknown error occured
		h.logger.Error("something went wrong", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// set up our page data
	type pageData struct {
		User    bool
		Product app.Product
	}
	data := pageData{
		User:    false,
		Product: product,
	}

	// render the templates
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}
```

Alright, there's quite a few things we have to deal with now.

- First of, we probably notice that our `productID` is of type `string` and we require `int`
	- So we use `strconv.Atoi` (which is short for 'string convert' package, with 'ascii to int' method)
	- If this errors, we can return a `400 Bad Request`
	- At this point I also renamed it to `productIDParam`, to indicate this is raw data from the request
- Now we call our new method: `h.Product.ShowProduct()`
	- We use a `switch` statement to deal with errors (maybe overkill for 2 conditions, but we might end up with a lot)
	- The first error we try to deal with is our newly created `app.ProductNotFound`, this results of course in `404 Not Found`
	- Any other error (regular old `err != nil` check) we consider an unknown error (clearly not from our domain model): `500 Internal Server Error`
- Our `pageData` is now also old, let's add our `Product` argument (and make sure to fill it)

When testing our page for product 1, 2, 3 and "gopher-plushie", we get:

```
Oct 29 20:59:31.340 DBG Request received method=GET url=/product/1 product_id=1
Oct 29 20:59:34.703 DBG Request received method=GET url=/product/2 product_id=2
Oct 29 20:59:37.168 DBG Request received method=GET url=/product/3 product_id=3
Oct 29 20:59:45.692 DBG Request received method=GET url=/product/gopher-plushie product_id=gopher-plushie
Oct 29 20:59:45.692 ERR couldn't convert product ID to int error="strconv.Atoi: parsing \"gopher-plushie\": invalid syntax"
```

Which is a 200, 200, 404, 400 respectively. 

However, these `http.Error()` calls are quite rough, it literally shows you a page with nothing but an error message and
the proper HTTP status code. This is fine for a JSON API, but not for our HTML site.

We'll have a look at 'Flashes' later, to give us proper user-friendly error handling.

### static/product/product.html

We however have more important business to do. We need to make our `product.html` page dynamic.

In `homepage.html` we saw us use `{{ . }}`, this dot means the current scope. 
We got that because we were using a `{{ range }}`.

Our scope in `product.html` will basically be the entire `pageData` struct we created. 
So as a hint, to get the product title: `{{ .Product.Name }}`

__Task__: Replace all static data in `product.html` with Go templating. (Maybe leave the image alone)

Restart your application and check multiple products. Is the data changing?

_Make sure you don't forget to update the text in our `{{ define "title" }}` template, 
and also the `<img>` has an `alt` field._

_FYI: If you **only** make changes in your HTML, you can reload your page without reloading the Go app_

### Formatting of the price

Most likely the price looks good for you, but what happens if you were to give discounts and your `12.99` was now 20% off?

We are currently using `%.2f` in our `fmt` methods to indicate 'give us a float with 2 point decimal accuracy'.
So even discounted values will be set correctly, but we're losing the raw data.

Sometimes you want to have both, raw data that we can put in HTML `data-` attributes for use by Javascript
and a nicely formatted price for the actual text.

```html
<h6 class="card-subtitle mb-2 text-body-secondary" data-price="12.99">€12.99</h6>
```

We could solve this by adding functions to our HTML template, but I've had better luck with creating helper 
functions on our model.

Let's go to `product.go` and add the following:

```go
func (p Product) FormattedPrice() string {
	return fmt.Sprintf("€%.2f", p.Price)
}
```

In our `product.html` go back to `{{ .Product.Price }}` and replace it with `{{ .Product.FormattedPrice }}`. 
Make sure to also delete the extra Euro symbol.

Refresh your application and make sure it works. Go's HTML templating allows you to call functions on models.

__Note__ that you have to remove the `()`, it will handle functions as if it was an attribute.

_We can even use `{{ .Product.String }}` and it will use our custom `String()` method that we build in step 2._

### static/homepage.html

You can go back to `static/homepage.html` and maybe add a little button in the `<li>` item, so you can easily
visit a product-detail page.

```html
<a href="/product/{{ .ID }}" class="btn btn-sm btn-primary">View</a>
```

## Error handling and flashes

Alright, our happy path works. But everytime we get an error, including a simple '404' we have a jarring user experience.

I've tried looking for history or an authorative source on flashes, but I can't find anything.

Long story short: a flash message is an error that's stored in a session and displayed on your site once.
Once it has been displayed, it will be deleted. Hence it's 'gone in a flash'.

It's a common concept in CRM frameworks such as Django (Python) or Laravel (PHP).

### 404

But in general we can also deal with 404's in a different way.

Our `httprouter` package has special 'catch-all' method.

Let's open our `handler/handler.go` file and add a method:

```go
func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.logger.WarnContext(r.Context(), "Page not found",
		"method", r.Method,
		"url", r.RequestURI,
	)
	tmpl := template.Must(template.ParseFiles(
		"static/layout.html",
		"static/404.html",
	))
	w.WriteHeader(http.StatusNotFound)
	if err := tmpl.Execute(w, nil); err != nil {
		h.logger.Error("failed to execute layout", "error", err)
		http.Error(w, "failed to create layout", http.StatusInternalServerError)
		return
	}
}
```

Within our `New()` method we can then add this line to our router:

```go
r.NotFound = http.HandlerFunc(h.notFound)
```

This `r.NotFound` is a special func that we can overwrite. It requires a `http.Handler` 
(so not the special `httprouter` one), however we can easily encapsulate this with `http.HandlerFunc()`
which creates a `http.Handler` out of any func we give it.

### static/404.html

Oh yeah, so that means we need to add a custom 404 page!

```html
{{ define "title" }}Page not found{{ end }}

{{ define "content" }}
<div class="row padding">
    <div class="col">
        <h2>Oops, couldn't find what you're looking for</h2>

        <div class="text-center">
            <a href="/" class="btn btn-primary">Return home</a>
        </div>
    </div>
</div>
{{ end }}
```

See if you can trigger your application to throw a 404 error..

### It's not working how we expect

So you might have tried to visit http://localhost:8000/product/3 (or any ID that results in a 404)
and found out that we're still not getting our custom 404 page.

This is because we're manually returning a 404. Let's go back to `handler/product.go` and change
our error case for 404.

From:

```go
	case errors.Is(err, app.ErrProductNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
```

To:

```go
	case errors.Is(err, app.ErrProductNotFound):
		h.notFound(w, r)
		return
```

So instead of giving just raw bytes and a status code, we're pointing back to our custom 404 handler,
which gives us back a proper layout but also a status code.

Now both `/product/3` and `/some-route-we-dont-know` return the same 404 page.

Maybe you should spend some time on making your 404 fancy. Find a nice gif on [giphy.com](https://giphy.com)

### Flashes

Let's add flashes though..

#### static/layout.html

Find your `<main>` component and add the flash code above the template. We want to show our message
on the top of the page, but within our `.container`.

```html
<main class="container mt-2">
    {{ if .Flashes }}
    {{ range $state, $msg := .Flashes }}
    <div class="alert alert-{{ $state }} alert-dismissible fade show" role="alert">
        {{ $msg }}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    </div>
    {{ end }}
    {{ end }}
    
    {{ template "content" . }}
</main>
```

#### handler/product.go

Go to both `products()` and `productByID()` and add the following attribute to your `pageData` struct.

```go
Flashes map[string]string
```

Make sure to fill the actual object you're sending in `Execute()`. This can be an empty `map`, but `nil`
also works. It's probably a good idea to extract this code out of the function into a variable, 
because it's becoming a little unwieldy.

```go
	data := pageData{
		User:     false,
		Flashes:  nil,
		Products: products,
	}
	// render the templates
	if err := tmpl.Execute(w, data); err != nil {
		// omitted for brevity
	}
```

In our 404 handler we're sending no data, so we don't need to add it there. But we eventually might,
if we want to send people to a 404 page **and** show them an error message.

#### handler/web.go

Let's use a library to make this easier for us. We're going to use `github.com/gorilla/sessions`.

```go
package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const cookieName = "flashes"

var store = sessions.NewCookieStore([]byte("this-really-should-be-a-secure-token-thats-not-stored-in-code"))

func storeAndSaveFlash(r *http.Request, w http.ResponseWriter, msg string) error {
	session, _ := store.Get(r, cookieName)
	session.AddFlash(msg)
	return session.Save(r, w)
}

func getFlashes(r *http.Request, w http.ResponseWriter) (map[string]string, error) {
	session, _ := store.Get(r, cookieName)
	flashes := session.Flashes()

	m := map[string]string{}
	for f := range flashes {
		fs := strings.SplitN(flashes[f].(string), "|", 2)
		if len(fs) == 2 {
			m[fs[0]] = fs[1]
		}
	}
	return m, session.Save(r, w)
}
```

Gorilla needs a `store`, which in our case is a `CookieStore`, then we can use that store to
get and save sessions. A `session` has methods for `AddFlash` and `Flashes`.

The helper functions here are just to help us do multiple things in a one-liner,
but in case of `getFlashes` also deal with our special formatting.

In our `layout.html` we call this `{{ range $state, $msg := .Flashes }}` so that our key is actually
the 'state', which in turn becomes the CSS class `alert-$state` and gives us 
[nice colours](https://getbootstrap.com/docs/5.3/components/alerts/#examples).

#### handler/product.go

Alright, back to our product handler.

```go
	flashes, err := getFlashes(r, w)
	if err != nil {
		h.logger.Warn("failed to get flashes", "error", err)
	}
	data := pageData{
		User:    false,
		Flashes: flashes,
		Product: product,
	}
```

For now I'm logging the errors, until we see it working. But we could probably ignore the error,
once we're certain it works.

So now we can read flashes, but nobody is setting one..

Let's scroll up to our `productID` validation code. Instead of throwing a 400 error, let's just 
create a flash message and redirect back to our homepage.

```go
	// validate our product ID
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "couldn't convert product ID to int", "error", err)
		_ = storeAndSaveFlash(r, w, "warning|Invalid product ID given")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
```

So we create a message, call `storeAndSaveFlash` which fetches our store, adds the flash and
persists the generated cookie to `w`, our ResponseWriter.

Then we send a redirect request with `http.Redirect` and a `307 Temporary Redirect`.
Because eventually this URL might work for this ID 
([see which redirect status code to use when](https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections)).

Then we `return` to stop the rest of the code from executing.

_Note: We're calling this from `productByID()`, but the flashes get rendered by `products()`. So you have to make sure 
that method also has the `getFlashes()` call and isn't still looking like `Flashes: nil`._

### Finally..

Alright, so now we redirect 404's to our dedicated 404 page and we can handle other type of
errors with our new flashes setup.

This means we're ready for more work.. or.. maybe we should create an API first?
