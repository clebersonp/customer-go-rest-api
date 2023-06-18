# Customer rest api
Simple rest API written in golang without any third-party library.
Just for study purposes only.

## Rest API Endpoints

* [x] - `GET /customers`. List all customers.
* [x] - `GET /customers/:id`. Returns the customer by id.
* [x] - `POST /customers`. Create a new customer.
* [x] - `POST /customers/:id/orders`. Create a new customer order.
* [x] - `GET /customers/:id/orders`. List all customer orders.
* [x] - `GET /customers/:id/orders/:id`. Returns a customer order by ids.
* [x] - `GET /admin`. Returns information only for the admin user with basic authentication.
* [x] - Validate the media type for `POST /customers`. It needs to be `application/json`.
* [x] - Validate the media type for `POST /customers/:id/orders`. It needs to be `application/json`.
* [x] - Validate any request for an `http method` other than those specified above, and returns the `header` value with `405 Method Not Allowed`.

## Run the application
In the terminal just type the following command to start the server in the root folder of this project:
```
$ ADMIN_NAME=admin ADMIN_PASSWORD=admin go run server.go
```
