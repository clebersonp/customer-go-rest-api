# Customer rest api
Simple rest API written in golang without any third-party library.
Just for study purposes only.

## Rest API Endpoints

* [ ] - `GET /customers`. List all customers.
* [ ] - `GET /customers/:id`. Returns the customer by id.
* [ ] - `POST /customers`. Create a new customer.
* [ ] - `POST /customers/:id/orders`. Create a new customer order.
* [ ] - `GET /customers/:id/orders`. List all customer orders.
* [ ] - `GET /customers/:id/orders/:id`. Returns a customer order by ids.
* [ ] - `GET /customers/random`. Returns a random customer.
* [ ] - `GET /customers/random/orders`. List orders from a random customer.
* [ ] - `GET /admin`. Returns information only for the admin user with basic authentication.
* [ ] - Validate the media type for `POST /customers`. It needs to be `application/json`.
* [ ] - Validate the media type for `POST /customers/:id/orders`. It needs to be `application/json`.
* [ ] - Returns the `location` for the `POST /customers` resource right after its creation.
* [ ] - Returns the `location` for the `POST /customers/:id/orders` resource right after its creation.
* [ ] - Validate any request for an `http method` other than those specified above, and returns the `header` value with `405 Method Not Allowed`.

## Run the application
In the terminal just type the following command to start the server in the root folder of this project:
```
$ go run server.go
```
