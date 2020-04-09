# shale

## Description
Shale is a simple RESTful API used for the management of a to-do list.  The backend functionality is written in go and database functionality comes from mysql.
Both the go API as well the accompanying database are dockerized and sit in their own container.  A `docker-compose` file can be used to build and run
both containers.

## Endpoints
Shale currently includes the following endpoints.  All of the below endpoint concern data for a single specified user:

Get Todos:  Returns a list of all todo item
    `GET: /todo/<username>`
    `username:  string`

Get Active/Inactive: Returns a list of all todo items that have either an active or inactive status, as specified
    `GET: /todo/<username>/active/<active>`
    `username: string`
    `active:  boolean`

Get by Priority:  Return a list of todo items with the provided priority level
    `GET: /todo/<username>/highs/<priority>`
    `username: string`
    `priority: integer`

Get by Category:  Return a list of todo items with the provided category
    `GET: /todo/<username>/cat/<category>`
    `username: string`
    `category: string`

Get by ID:  Return a list of todo items with the provided id
    `GET: /todo/<username>/id/<id>`
    `username: string`
    `id: integer`

Add Item: Add a new todo item to the list
    `POST: /todo/<username>/add --data { <types.TodoData> }`
    `username: string`


Change Title:  Change the title of a todo item based on its id
    `POST: /todo/<username>/ctitle/<id> --data { <types.TodoData>}`
    `username: string`
    `id: integer`

Change Priority:  Change the priority of a todo item based on its id
    `POST: /todo/<username>/cpri/<id> --data { <types.TodoData> }`
    `username: string`
    `id integer`

Change Active:  Change whether a todo item is active or inactive based on its id
    `POST: /todo/<username>/cactive/<id> --data { <types.TodoData>}`

Remove by Title: Remove a todo item from the list based on its title
    `DELETE: /todo/<username>/rmtitle --data { <types.TodoData>}`
    `username: string`

Remove by Priority: Remove all todo items from the list that have a given priority level
    `DELETE: /todo/<username>/rmpri --data { <types.TodoData>}`
    `username: string`

Remove by ID: Remove a todo item from the list based on its id
    `DELETE: /todo/<username>/rmid --data "{ <types.TodoData>}`

For the above endpoints that include a data payload, the types.TodoData is a go struct with the following attributes.  It is onyl necessary to return the individual
values of concern for a given endpoint:

`Name        string         json:"acct_name"`
`Title       string         json:"title"`
`Body        string         json:"body"`
`Category    string         json:"category"`
`Priority    int            json:"item_priority"`
`PublishDate mysql.NullTime json:"publish_date"`
`Active      bool           json:"active"`
`ID          int            json:"id"`

Note the following conditions:
1.  ID is assigned byu the database and is immutable
2.  Priority can be any integer.  Negative values are acceptable.  A priority of 0 is treated is if it does not have a priority.
3.  Todo notes are added with a default active value of true


## Example Usage
The API is currently being hosted at `73.78.155.49:8080`.  Below are usage examples for each endpoint using curl.  The examples assume a user `tom` and that there is already data in the database that satisfies the request.

Get todos: `curl -vv 73.78.155.49:8080/todo/tom`

Get Active todos: `curl -vv 73.78.155.49:8080/todo/tom/active/1`

Get Inactive todos: `curl -vv 73.78.155.49:8080/todo/tom/active/0`

Get todos by Priority of 3 or lower: `curl -vv 73.78.155.49:8080/todo/tom/highs/3`

Get todos by Category: `curl -vv 73.78.155.49:8080/todo/tom/cat/pandemic`

Get todos by ID: `curl -vv 73.78.155.49:8080/todo/tom/id/4`

ADD a new todo: `curl -vv -X POST localhost:8080/todo/tom/add --data {"Title": "Cure covid-19", "item_priority": 5, "category": "pandemic"}`

Change todo Title: `curl -vv -X POST localhost:8080/todo/tom/ctitle/4 --data {"title": "Research covid-19 first"}`

Change todo Priority: `curl -vv -X POST localhost:8080/todo/tom/cpri/4 --data {"item_priority": 2}`

Change to Active Status:  `curl -vv -X POST localhost:8080/todo/tom/cactive/2 --data {"active": false}`

Remove todo Based on Title: `curl -vv -X DELETE localhost:8080/todo/tom/rmtitle --data {"title": "Research covid-19 first"}`

Remove todos with Given Priority: `curl -vv -X DELETE localhost:8080/todo/tom/rmpri --data {"item_priority": 0}`

Remove todo with id: `curl -vv -X DELETE localhost:8080/todo/tom/rmid --data {"id": 8}`


## Potential Pitfalls/Assumptions
1.  This API was developed in a Windows environment running Docker Desktop for Windows.  While the environment is relatively simple,
it should be noted that there may be potential differences in operation in Windows environments.  Please notify me of any environment-related
issues and I will see about addressing them.

2.  In addition to Docker, curl for Windows requires the use of escape characters in curl commands for characters such as quotes.  As an example, during testing, the curl
command actually executed in place of `curl -vv -X POST localhost:8080/todo/tom/ctitle/4 --data {"title": "Research covid-19 first"}` would actually look like:
`curl -v -X POST localhost:8080/todo/tom/ctitle/4 --data "{\"title\": \"Changed TITLE\"}"\"`

Therefore, while the above curl commands are trivial, it is the Windows counterpart of each curl command that was used for testing.