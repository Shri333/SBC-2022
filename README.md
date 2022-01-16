# Shopify Backend Developer Intern Challenge - Summer 2022
This is an API for managing inventory items. The API is written in Go and uses Gin and sqlite3. You can view a live deployment [here](#).

## Building From Source
In order to build from source, you need to have a modern Go toolchain (1.15+), a modern C/C++ toolchain (to build go-sqlite3), and sqlite3 installed.

You can download sqlite3 [here](https://www.sqlite.org/download.html). Make sure that after installing, sqlite3 is present in PATH.

After cloning the repository and installing all required tools, in order to build the server, inside of the repository directory, run `go mod tidy && go build`.

## Running The Application
After building, on Linux and macOS, run `./SBC-2022` to launch the server. On Windows, run `SBC-2022.exe`. 

The server (by default) runs on port 8080. To change the port, set the `PORT` environmental variable to the appropriate port.

To view the web application (if PORT is 8080), visit `http://localhost:8080`.

## API Features
- CRUD operations for items and groups (a group can contain multiple items and not every item has to be in a group)
- Input validation (server responds with appropriate error message and code for bad input)
- Supports both JSON and XML for request bodies (but responds in JSON)

## API Routes
- **GET** `/api/items`
  - Responds with a list of items. Each item has an id, name, count, and (optional) group.
  - ```
    $ curl -X GET http://localhost:8080/api/items
    [{"id":1,"name":"apples","count":100,"group":{"id":1,"name":"fruits"}},{"id":2,"name":"bananas","count":150,"group":{"id":1,"name":"fruits"}},{"id":3,"name":"carrots","count":120}]
    ```
- **GET** `/api/items/:id`
  - Responds with an item with the given id. An error message/code (404) is returned if no item exists with the given id.
  - ```
    $ curl -X GET http://localhost:8080/api/items/1
    {"id":1,"name":"apples","count":100,"group":{"id":1,"name":"fruits"}}
    ```
- **GET** `/api/groups`
  - Responds with a list of groups. Each group has an id, name, and an (optional) list of items. Items that do not have a group are put in a pseudo-group (id: 0, name: "").
  - ```
    $ curl -X GET http://localhost:8080/api/groups
    [{"id":0,"name":"","items":[{"id":3,"name":"carrots","count":120}]},{"id":1,"name":"fruits","items":[{"id":1,"name":"apples","count":100},{"id":2,"name":"bananas","count":150}]}]
    ```
- **GET** `/api/groups/:id`
  - Responds with a group with the given id. An error message/code (404) is returned if no group exists with the given id.
  - ```
    $ curl -X GET http://localhost:8080/api/groups/1
    {"id":1,"name":"fruits","items":[{"id":1,"name":"apples","count":100},{"id":2,"name":"bananas","count":150}]}
    ```
- **POST** `/api/items`
  - Creates a new item and responds with the new item. The item name and count are required (the groupId is optional). The item name must be a unique non-empty string, and the item count must be a positive integer (greater than zero). Responds with an error message/code (404) if no group exists with the given groupId.
  - ```
    $ curl -X POST http://localhost:8080/api/items -H 'content-type: application/json' -d '{"name":"oranges","count":50,"groupId":1}'
    {"id":4,"name":"oranges","count":50,"group":{"id":1,"name":"fruits"}}
    ```
- **POST** `/api/groups`
  - Creates a new group and responds with the new group. The group name is required/unique and must be a non-empty string.
  - ```
    curl -X POST http://localhost:8080/api/groups -H 'content-type: application/json' -d '{"name":"vegetables"}'
    {"id":2,"name":"vegetables"}
    ```
- **PUT** `/api/items/:id`
  - Updates an existing item and responds with the updated item. All unchanged and changed fields of the existing item are required. If the existing item was in a group and the groupId is omitted, the item is removed from the group it was originally in. Responds with an error message/code (404) if no item exists with the given id or no group exists with the given groupId.
  - ```
    $ curl -X PUT http://localhost:8080/api/items/1 -H 'content-type: application/json' -d '{"name":"pears","count":50}'
    {"id":1,"name":"pears","count":50}
    ```
- **PUT** `/api/groups/:id`
  - Updates an existing group and responds with the updated group. The name of the group is required/unique and must be a non-empty string. Responds with an error message/code (404) if no group exists with the given id.
  - ```
    $ curl -X PUT http://localhost:8080/api/groups/1 -H 'content-type: application/json' -d '{"name":"Fresh Fruits"}'
    {"id":1,"name":"Fresh Fruits","items":[{"id":1,"name":"apples","count":100},{"id":2,"name":"bananas","count":150}]}
    ```
- **DELETE** `/api/items/:id`
  - Deletes an item with the given id and responds with an OK status code (200). Responds with an error message/code (404) if no item exists with the given id.
  - ```
    $ curl -X DELETE http://localhost:8080/api/items/1
    ```
- **DELETE** `/api/groups/:id`
  - Deletes a group with the given id and responds with an OK status code (200). If the group contains any items, those items are also deleted. Responds with an error message/code (404) if no group exists with the given id.
  - ```
    $ curl -X DELETE http://localhost:8080/api/groups/1
    ```
