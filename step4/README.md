# Step 4: JSON REST API

We've set up our basic webshop, with HTML and Go. Business is booming and our boss
has decided to hire a dedicated frontend team, allowing us to work on the backend.

For this however, we need to supply them with a JSON API. 
([Learn more about REST API design](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/))

## API endpoints

We are going to make REST-ish endpoints. If we were to fully honour the
Representational State Transfer (REST), we'd have to do way more. So let's consider it RESTful
and leave pedantics behind us. We have JSON to serve!

### handler/handler.go

Let's add to our routes by accepting some API endpoints.

```go
// API routes
r.GET("/api/products", h.apiProducts)
r.GET("/api/products/:id", h.apiProductByID)
```

We'll be trying to convert our current service to this API, so `ListProducts()` becomes
`GET /api/products` and `ShowProduct()` becomes `GET /api/products/:id`.

The frontend team really wants to test the "Add to cart" button and asked us to also introduce baskets.

We need a way to 'get a basket, or create if there's none', 'add product to basket' and 'remove product from basket'.

```go
r.GET("/api/basket", h.apiBasket)
r.POST("/api/basket/add", h.apiAddToBasket)
r.POST("/api/basket/remove", h.apiRemoveFromBasket)
```

### handler/api.go

Let's implement our first API route.

To have a Go webserver respond with JSON, we really only need to set the Content-Type 
(although most likely most clients will auto-detect this) and return the JSON representation of our products.

```go
func (h *Handler) apiProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	products, err := h.Product.ListProducts(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch products", "error", err)
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		h.logger.Error("failed to write products JSON", "error", err)
		http.Error(w, "failed to write products JSON", http.StatusInternalServerError)
	}
}
```

So here we fetch our products, using the Product Service we already have in our Dependencies.

Then we set the correct `Content-Type` header on the ResponseWriter, initiate a new JSON Encoder, who's `io.Writer`
is actually our `http.ResponseWriter` (in other words: we directly write our converted JSON to the browser)
and call the `Encode()` method with our products.

__Task__: Comment out the other API routes and start your server. Open http://localhost:8000/api/products in your browser

If we then use `curl` to check our new endpoint, we see this.

```shell
$ curl http://localhost:8000/api/products
[{"ID":1,"Name":"Gopher plushie","Description":"A small purple Gophier plushie, perfect for kids and adults alike.","Image":"","Price":12.99},{"ID":2,"Name":"PHP Elephant plushie","Description":"An elephant with the PHP logo, available in blue and pink","Image":"","Price":20}]
```

Depending on what browser you use, you might see a nice readable JSON array, 
but CURL doesn't do this and the actual data we're sending also doesn't 'prettify' our JSON.

We could use the method `SetIndent()` on our `Encoder`, but that means we need to remove our one-liner, and initiate them separately.
If you have `jq` installed, you can also use that. 

#### Errors in APIs

You can see when we handle our errors, we return for example the following:

```go
http.Error(w, "failed to fetch products", http.StatusInternalServerError)
```

If this gets triggered, your client sees a 500 error with plain-text "failed to fetch products".
This could be fine, but if your API is always supposed to return JSON, then it's not nice.

This is really a design question. For now, we're just going to use our `http.Error()` helper function.

#### JSON Struct Tags

The frontend team got back to us and told us they had already a lot of code, and they used the following parameters:
`id, name, desc, img, price`.

Unless you want to escalate this with your PO, this means we need to change our `Product` model by adding 'struct tags' 
([read more](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)).
You can consider these like Annotations in Java, as in meta programming.

```go
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	...
}
```

_Struct tags are started by backtics and use quotes for the name. There's some special cases for `-` and `,omitempty` but we won't touch those._

**We have a design decision of our own to make.**
Are we going to update our `/product.go` (our domain model) or do we create a new product specifically for our API responses?

According to the 'svc' mantra we should create a new model, but we could try to just get away with updating our existing model. 😉

__Task__: Update the rest of your `/product.go` with JSON struct tags, according the 'specification' of the frontend team.

Start your server again and either browse or `curl` to http://localhost:8000/api/products. 

Do you have correctly name attributes now?

```json
{
	"id": 1,
	"name": "Gopher plushie",
	"desc": "A small purple Gophier plushie, perfect for kids and adults alike.",
	"img": "",
	"price": 12.99
}
```

#### Finish the product API

__Task__: Implement the `r.GET("/api/products/:id", h.apiProductByID)` route (don't forget to uncomment)

You should already have most of this code. For sending the JSON part, look at your newly created `h.apiProducts()` method.

For getting the correct product (including fetching the `:id` parameter) see your code `handler/product.go`, specifically
the `productByID()` method.

## Baskets

We seem to have a new requirement with regards to having 'basket functionality'.

We have none of this, so we need to set up: domain models, storages, services.

### /basket.go

Let's start by creating our domain model.

```go
type Basket struct {
	UserID int
	Items  []BasketItem
}

type BasketItem struct {
	ProductID int
	Quantity  int
}
```

We're going to pretend we have Users and they can have a Basket, the Basket consists of BasketItems.

A BasketItem has a Product ID and a Quantity.

### services/services.go

Time to update your service interfaces. Add the following interface next to your existing ProductService.

```go
type BasketService interface {
	GetBasket(ctx context.Context, userID int) (app.Basket, error)
	AddToBasket(ctx context.Context, userID, productID, quantity int) error
	RemoveFromBasket(ctx context.Context, userID, productID, quantity int) error
}
```

These methods seem like enough to fulfil our API endpoints. 

- We can look up a basket by User ID
- We can add a product to a User's basket, with specified quantity
- We can remove a product from a User's basket, specifying again quantity

What do you think of the names? Are they clear? When we say "add to basket" is it always a product, or could it be a coupon?

### storage/storage.go

__Task__: Add a similar interface to the storage layer.

### storage/memory_basket.go

Since our tech lead still hasn't decided on a database, we're going to have to rely on our memory database. Yes, that `map`.

Create a new struct in `storage/memory_basket.go` (or something that matches your naming structure).

```go
type BasketRepo struct {
	baskets map[int]app.Basket
}
```

This map (which has User IDs as keys) should be enough to store our baskets in memory.

__Task__: Create a constructor that returns a `*BasketRepo`, make sure to initialize the map to prevent compile errors.

You can have a look at your Product Repository, they're quite similar.

__Task__: Implement the missing methods from your `BasketRepository` interface. (See notes below)

_Note 1_: You're most likely going to do a map lookup, but if the basket is not found, we should probably
return an error: `app.ErrBasketNotFound`. 

Create this domain model error in `/basket.go`, if you're unsure you can lookup how you did it in `/product.go`.

_Note 2_: For our `GetBasket()` method we actually want to create an entry in the map and not return an error. Consider it a `PUT`.
Don't forget about manually setting the `UserID`.

_Note 3_: Also it's going to be a bit of work to properly deal with `quantity` because our `Basket.Items` is a slice
and so we're going to have the 'challenge' of having to loop over the items to consolidate. The same thing 
needs to be taken care of when removing items, especially if you remove for example 'all but one'.

So for now ignore `quantity` and consider the value always to be `1`.

_Note 4_: For dealing with `Basket.Items`, you can have a look at https://go.dev/wiki/SliceTricks.


### services/basket.go

Time to implement our `BasketSvc` service. We can refresh our mind by looking at `services/product.go`.\

```go
type BasketSvc struct {
	repo storage.BasketRepository
}

func NewBasketService(repo storage.BasketRepository) *BasketSvc {
	return &BasketSvc{repo: repo}
}
```

Luckily we only have a storage dependency and no `CouponService`, nor do we need to talk to `UserService`, 
so our dependencies and constructor are quite straightforward.

__Task__: Implement the missing `BasketService` interface methods.

## Finishing up the API

Alright, seems like our _Basket Backend_ is in order. Time to connect the API endpoints to it.

So we just load `h.Basket.GetBasket..`.. wait a minute; we have not added our Dependency to our app!

### Dependency inject our Basket

Alright, so first we need to go to `handler/handler.go` and add `Basket services.BasketService` to our `Dependencies` struct.

This means we then need to go to `cmd/app/main.go` and supply it here as well. Our `Handler` should now know about `BasketService`
and we can call its methods in our API handlers.


### handler/api.go

Let's start with the most straightforward one. Retrieving a user's basket.

```go
func (h *Handler) apiBasket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := 1
	basket, err := h.Basket.GetBasket(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to fetch basket", "error", err)
		http.Error(w, "failed to fetch basket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(basket); err != nil {
		h.logger.Error("failed to write basket JSON", "error", err)
		http.Error(w, "failed to write basket JSON", http.StatusInternalServerError)
	}
}
```

So we don't have a user system yet and we kind of painted ourselves in a corner with our API design.
We have no way of getting the User ID. Normally this would be fine, because that information comes out
of the context, either directly or via an access token.

For now we're just going to say that everyone shares a basket and you have to be really quick to press "Order"!

_Don't forget to uncomment your Get Basket route in `handler/handler.go`.._

Restart your server and browse or `curl` to http://localhost:8000/api/basket. Are you seeing your empty basket?

```json
{"UserID":1,"Items":[]}
``` 

#### Adding a product

Our "add product" route is a POST, so we need to start dealing with form data. Although we could also JSON input.

This is something that you need to design into your API.

```go
r.ParseForm()
productIDParam := r.Form.Get("product_id")
```

`ParseForm()` is a helper method on the `http.Request` that does the following.

```text
// ParseForm populates r.Form and r.PostForm.
//
// For all requests, ParseForm parses the raw query from the URL and updates
// r.Form.
//
// For POST, PUT, and PATCH requests, it also reads the request body, parses it
// as a form and puts the results into both r.PostForm and r.Form. Request body
// parameters take precedence over URL query string values in r.Form.
//
// If the request Body's size has not already been limited by [MaxBytesReader],
// the size is capped at 10MB.
```

Which allows us to use some nice methods like our `r.Form.Get("product_id")` example.

__Task__: Implement the `apiAddToBasket` method with the above code (you can hardcode `1` for both `userID` and `quantity`)

Don't forget to uncomment the specific route and restart your server.

- Use the following `curl` command: `curl -i -d "product_id=1" http://localhost:8000/api/basket/add` 
	- (by adding the `-d` flag you are saying this is a POST request, you don't need `-X POST`)
- Notice you get an error, create your basket first: `curl -i http://localhost:8000/api/basket`
- Now run the first `curl` command again (you can use Up-arrow twice in your CLI)
- If you receive no error, you can now check your basket again: `curl -i http://localhost:8000/api/basket`

```text
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 17 Jan 2025 22:58:38 GMT
Content-Length: 52

{"UserID":1,"Items":[{"ProductID":1,"Quantity":1}]}
```

__Note__: Everytime you restart your server, you'll have to create a basket, since our in-memory repository is temporary.

### (Optional) Implement the 'delete item from basket' endpoint

If you still have some time, implement the `r.POST("/api/basket/remove", h.apiRemoveFromBasket)` route.