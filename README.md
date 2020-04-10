# shale

## Description
Shale is a simple RESTful API used for the management of a to-do list.  The backend functionality is written in go and database functionality comes from mysql.
Both the go API as well the accompanying database are dockerized and sit in their own container.  A `docker-compose` file can be used to build and run
both containers.

## Endpoints
Shale currently includes the following endpoints.  Every endpoint requires a `username` to select the necessary todo list.  In this way, the system allows for multiple lists.  All of the below endpoint concern data for a single specified user:

Get Todos:  Returns a list of all todo item<br>
    `GET: /todo/<username>`<br>
    `username:  string`<br>

Get Active/Inactive: Returns a list of all todo items that have either an active or inactive status, as specified<br>
    `GET: /todo/<username>/active/<active>`<br>
    `username: string`<br>
    `active:  boolean`<br>

Get by Priority:  Return a list of todo items with the provided priority level<br>
    `GET: /todo/<username>/highs/<priority>`<br>
    `username: string`<br>
    `priority: integer`<br>

Get by Category:  Return a list of todo items with the provided category<br>
    `GET: /todo/<username>/cat/<category>`<br>
    `username: string`<br>
    `category: string`<br>

Get by ID:  Return a list of todo items with the provided id<br>
    `GET: /todo/<username>/id/<id>`<br>
    `username: string`<br>
    `id: integer`<br>

Add Item: Add a new todo item to the list<br>
    `POST: /todo/<username>/add --data { <types.TodoData> }`<br>
    `username: string`<br>


Change Title:  Change the title of a todo item based on its id<br>
    `POST: /todo/<username>/ctitle/<id> --data { <types.TodoData>}`<br>
    `username: string`<br>
    `id: integer`<br>

Change Priority:  Change the priority of a todo item based on its id<br>
    `POST: /todo/<username>/cpri/<id> --data { <types.TodoData> }`<br>
    `username: string`<br>
    `id integer`<br>

Change Active:  Change whether a todo item is active or inactive based on its id<br>
    `POST: /todo/<username>/cactive/<id> --data { <types.TodoData>}`<br>

Remove by Title: Remove a todo item from the list based on its title<br>
    `DELETE: /todo/<username>/rmtitle --data { <types.TodoData>}`<br>
    `username: string`<br>

Remove by Priority: Remove all todo items from the list that have a given priority level<br>
    `DELETE: /todo/<username>/rmpri --data { <types.TodoData>}`<br>
    `username: string`<br>

Remove by ID: Remove a todo item from the list based on its id<br>
    `DELETE: /todo/<username>/rmid --data "{ <types.TodoData>}`<br>
    `username: string`<br>

For the above endpoints that include a data payload, the types.TodoData is a go struct with the following attributes.  It is onyl necessary to return the individual
values of concern for a given endpoint:

`Name        string         json:"acct_name"`<br>
`Title       string         json:"title"`<br>
`Body        string         json:"body"`<br>
`Category    string         json:"category"`<br>
`Priority    int            json:"item_priority"`<br>
`PublishDate mysql.NullTime json:"publish_date"`<br>
`Active      bool           json:"active"`<br>
`ID          int            json:"id"`<br>

Note the following conditions:
1.  ID is assigned byu the database and is immutable
2.  Priority can be any integer.  Negative values are acceptable.  A priority of 0 is treated is if it does not have a priority.
3.  Todo notes are added with a default active value of true


## Example Usage
The API is currently being hosted at `73.78.155.49:8080`.  Because authentication is not required for this project, a username can be simply chosen, and todos added or changed based on that username, with the format listed in the above section.

Below are usage examples for each endpoint using curl.  The examples assume a username of `tom` and that there is already data in the database that satisfies the request.  These examples are simply that; there is no telling whats actually in the database at any given moment.

Get todos: `curl -vv 73.78.155.49:8080/todo/tom`

Get Active todos: `curl -vv 73.78.155.49:8080/todo/tom/active/1`

Get Inactive todos: `curl -vv 73.78.155.49:8080/todo/tom/active/0`

Get todos by Priority of 3 or lower: `curl -vv 73.78.155.49:8080/todo/tom/highs/3`

Get todos by Category: `curl -vv 73.78.155.49:8080/todo/tom/cat/pandemic`

Get todos by ID: `curl -vv 73.78.155.49:8080/todo/tom/id/4`

ADD a new todo: `curl -vv -X POST 73.78.155.49:8080/todo/tom/add --data {"Title": "Cure covid-19", "item_priority": 5, "category": "pandemic"}`

Change todo Title: `curl -vv -X POST 73.78.155.49:8080/todo/tom/ctitle/4 --data {"title": "Research covid-19 first"}`

Change todo Priority: `curl -vv -X POST 73.78.155.49:8080/todo/tom/cpri/4 --data {"item_priority": 2}`

Change to Active Status:  `curl -vv -X POST 73.78.155.49:8080/todo/tom/cactive/2 --data {"active": false}`

Remove todo Based on Title: `curl -vv -X DELETE 73.78.155.49:8080/todo/tom/rmtitle --data {"title": "Research covid-19 first"}`

Remove todos with Given Priority: `curl -vv -X DELETE 73.78.155.49:8080/todo/tom/rmpri --data {"item_priority": 0}`

Remove todo with id: `curl -vv -X DELETE 73.78.155.49:8080/todo/tom/rmid --data {"id": 8}`


## Potential Pitfalls/Assumptions
1.  This API was developed in a Windows environment running Docker Desktop for Windows.  While the environment is relatively simple,
it should be noted that there may be potential differences in operation in Windows environments.  Please notify me of any environment-related
issues and I will see about addressing them.

2.  In addition to Docker, curl for Windows requires the use of escape characters in curl commands for characters such as quotes.  As an example, during testing, the curl
command actually executed in place of `curl -vv -X POST localhost:8080/todo/tom/ctitle/4 --data {"title": "Research covid-19 first"}` would actually look like:
`curl -v -X POST localhost:8080/todo/tom/ctitle/4 --data "{\"title\": \"Changed TITLE\"}"\"`

Therefore, while the above curl commands are trivial, it is the Windows counterpart of each curl command that was used for testing.
