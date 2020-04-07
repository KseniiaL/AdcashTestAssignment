# AdcashTestAssignment
A simple REST API for the catalog of products.
<br/><br/>The application should contain:
<br/>● Product categories;
<br/>● Products which belong to some category (one product may belong to one category);
<br/><br/>The following actions should be implemented:
<br/>● Getting the list of all categories;
<br/>● Getting a category by ID;
<br/>● Getting the list of all products;
<br/>● Getting the list of products of the concrete category;
<br/>● Getting a product by ID;
<br/>● Create/update/delete of category;
<br/>● Create/update/delete of product;

To deploy and run the application Go, "github.com/gorilla/mux", "github.com/stretchr/testify/assert", and "github.com/rs/xid"(for generating unique IDs) should be installed.

To install dependencies run:
<br/>```go get -u github.com/gorilla/mux```
<br/>```go get -u github.com/stretchr/testify/assert```
<br/>```go get -u github.com/stretchr/testify```

Run the following commands to run/test application:
<br/>```go run main.go```
<br/>```go test ./... -cover```
<br/>Be sure to run the commands while being in the project's directory.
