

# Introductions

This programm supports following endpoints and actions:

| Method | URL Pattern        | Action                                           |
| :-------| :-------------------| :-------------------------------------------------|
| GET    | /healthcheck       | Show application health and version information. |
| GET    | /persons           | Show the details of all persons.                 |
| POST   | /persons           | Create a new person.                             |
| GET    | /persons/:id       | Show the details of a specific person.           |
| GET    | /persons/color/:id | Shows all people with the same favorite color.   |

## Prerequisites

### Go 1.24

The current program was created with Go version 1.24.12.

```
$ go version
go version go1.24.12 linux/amd64
```

```
$ go version
go version version go1.24.12 windows/amd64
```

## Building

### Linux x64
```
$ go build -o bin/assecor_api ./api
```

### Windows

You must have the correct version of gcc and the necessary runtime libraries installed on Windows.
One method to do this is using msys64. To begin, install msys64 using their installer.
Once you installed msys64, open a msys64 shell and run:

```
pacman -S mingw-w64-ucrt-x86_64-gcc
```

Select "yes" when necessary; it is okay if the shell closes. Then, add gcc to the path using
whatever method you prefer.

```
$ go build -o .\bin\assecor_api.exe .\api
```

# Background knowledge

## Project setup and structure

* The `bin` directory will contains the compiled application binaries, ready for deployment
to a production server.
* The `api` directory will contain the application-specific code for the API
application. This will include the code for running the server, reading and writing HTTP
requests.
* The `internal` directory will contain various ancillary packages used by the API. It will
contain the code for interacting with our database, doing data validation and so on.
The Go code under `api` will import the packages in the internal directory (but never the other way around).
* The `go.mod` file will declare the project dependencies, versions and module path.
* The `go.sum` is an append-only log of checksums, used to verify the integrity of modules downloaded during builds.
* All current test files (with the attachment _test) are located in the api folder.

## Logging

All information and error messages are displayed via a logging output.
The log entries are simple free-form messages prefixed by the current
date and time, written to the standard out stream. A example of these are the log
entries that we see when we start the API:

```
$ go run ./api
2026/02/02 11:26:37 database connection established
2026/02/02 11:26:37 starting server :4000
```

## Program arguments

The program currently supports two parameters:

* The argument `port` specifies the port for the server; 4000 is configured by default.
* The argument `dsn` specifies where the data to be loaded is located. No file is specified by default.


# Testing

## Unit-tests

```
$ go test -v ./api
=== RUN   TestPing
--- PASS: TestPing (0.00s)
=== RUN   TestShowPerson
=== RUN   TestShowPerson/Valid_ID
=== RUN   TestShowPerson/Non-existent_ID
=== RUN   TestShowPerson/Negative_ID
=== RUN   TestShowPerson/Decimal_ID
=== RUN   TestShowPerson/String_ID
=== RUN   TestShowPerson/Trailing_slash
--- PASS: TestShowPerson (0.00s)
    --- PASS: TestShowPerson/Valid_ID (0.00s)
    --- PASS: TestShowPerson/Non-existent_ID (0.00s)
    --- PASS: TestShowPerson/Negative_ID (0.00s)
    --- PASS: TestShowPerson/Decimal_ID (0.00s)
    --- PASS: TestShowPerson/String_ID (0.00s)
    --- PASS: TestShowPerson/Trailing_slash (0.00s)
=== RUN   TestShowPersons
=== RUN   TestShowPersons/No_persons
=== RUN   TestShowPersons/Valid
=== RUN   TestShowPersons/Trailing_slash
--- PASS: TestShowPersons (0.00s)
    --- PASS: TestShowPersons/No_persons (0.00s)
    --- PASS: TestShowPersons/Valid (0.00s)
    --- PASS: TestShowPersons/Trailing_slash (0.00s)
=== RUN   TestFilterByPersonsFavoriteColor
=== RUN   TestFilterByPersonsFavoriteColor/Valid_ColorID(blue)
=== RUN   TestFilterByPersonsFavoriteColor/Valid_ColorID(green)
=== RUN   TestFilterByPersonsFavoriteColor/Valid_ColorID(yellow)
=== RUN   TestFilterByPersonsFavoriteColor/Non-existent_ColorID
=== RUN   TestFilterByPersonsFavoriteColor/Negative_ColorID
=== RUN   TestFilterByPersonsFavoriteColor/Trailing_slash
--- PASS: TestFilterByPersonsFavoriteColor (0.00s)
    --- PASS: TestFilterByPersonsFavoriteColor/Valid_ColorID(blue) (0.00s)
    --- PASS: TestFilterByPersonsFavoriteColor/Valid_ColorID(green) (0.00s)
    --- PASS: TestFilterByPersonsFavoriteColor/Valid_ColorID(yellow) (0.00s)
    --- PASS: TestFilterByPersonsFavoriteColor/Non-existent_ColorID (0.00s)
    --- PASS: TestFilterByPersonsFavoriteColor/Negative_ColorID (0.00s)
    --- PASS: TestFilterByPersonsFavoriteColor/Trailing_slash (0.00s)
=== RUN   TestCreatePersons
=== RUN   TestCreatePersons/Valid_person(ID=1)
=== RUN   TestCreatePersons/Valid_person(ID=2)
=== RUN   TestCreatePersons/Empty_name
=== RUN   TestCreatePersons/Empty_lastname
=== RUN   TestCreatePersons/Empty_zipcode
=== RUN   TestCreatePersons/Zipcode_xxx_zipcode
=== RUN   TestCreatePersons/Empty_city
=== RUN   TestCreatePersons/Negative_colorID
=== RUN   TestCreatePersons/Out_of_range_colorID
=== RUN   TestCreatePersons/Trailing_slash
--- PASS: TestCreatePersons (0.00s)
    --- PASS: TestCreatePersons/Valid_person(ID=1) (0.00s)
    --- PASS: TestCreatePersons/Valid_person(ID=2) (0.00s)
    --- PASS: TestCreatePersons/Empty_name (0.00s)
    --- PASS: TestCreatePersons/Empty_lastname (0.00s)
    --- PASS: TestCreatePersons/Empty_zipcode (0.00s)
    --- PASS: TestCreatePersons/Zipcode_xxx_zipcode (0.00s)
    --- PASS: TestCreatePersons/Empty_city (0.00s)
    --- PASS: TestCreatePersons/Negative_colorID (0.00s)
    --- PASS: TestCreatePersons/Out_of_range_colorID (0.00s)
    --- PASS: TestCreatePersons/Trailing_slash (0.00s)
PASS
ok      assecor.assessment.test/api     (cached)
```

## REST-API

```
$ go run ./api -dsn sample-input.csv
2026/02/02 13:24:11 database connection established
2026/02/02 13:24:11 [csv]: record on line 8: wrong number of fields
2026/02/02 13:24:11 starting server :4000
```

### GET /persons

```
$ curl -i localhost:4000/persons
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 02 Feb 2026 12:24:58 GMT
Content-Length: 1274

[
  {
    "id": 1,
    "name": "Hans",
    "lastname": "Müller",
    "zipcode": "67742",
    "city": "Lauterecken",
    "color": "blau"
  },
  {
    "id": 2,
    "name": "Peter",
    "lastname": "Petersen",
    "zipcode": "18439",
    "city": "Stralsund",
    "color": "grün"
  },
  {
    "id": 3,
    "name": "Johnny",
    "lastname": "Johnson",
    "zipcode": "88888",
    "city": "made up",
    "color": "violett"
  },
  {
    "id": 4,
    "name": "Milly",
    "lastname": "Millenium",
    "zipcode": "77777",
    "city": "made up too",
    "color": "rot"
  },
  {
    "id": 5,
    "name": "Jonas",
    "lastname": "Müller",
    "zipcode": "32323",
    "city": "Hansstadt",
    "color": "gelb"
  },
  {
    "id": 6,
    "name": "Tastatur",
    "lastname": "Fujitsu",
    "zipcode": "42342",
    "city": "Japan",
    "color": "türkis"
  },
  {
    "id": 7,
    "name": "Anders",
    "lastname": "Andersson",
    "zipcode": "32132",
    "city": "Schweden - ☀",
    "color": "grün"
  },
  {
    "id": 8,
    "name": "Gerda",
    "lastname": "Gerber",
    "zipcode": "76535",
    "city": "Woanders",
    "color": "violett"
  },
  {
    "id": 9,
    "name": "Klaus",
    "lastname": "Klaussen",
    "zipcode": "43246",
    "city": "Hierach",
    "color": "grün"
  }
]
```

### GET /persons/:id

```
$ curl -i localhost:4000/persons/1
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 02 Feb 2026 12:26:31 GMT
Content-Length: 123

{
  "id": 1,
  "name": "Hans",
  "lastname": "Müller",
  "zipcode": "67742",
  "city": "Lauterecken",
  "color": "blau"
}
```

### GET /persons/color/:id

```
$ curl -i localhost:4000/persons/color/2
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 02 Feb 2026 12:27:41 GMT
Content-Length: 431

[
  {
    "id": 2,
    "name": "Peter",
    "lastname": "Petersen",
    "zipcode": "18439",
    "city": "Stralsund",
    "color": "grün"
  },
  {
    "id": 7,
    "name": "Anders",
    "lastname": "Andersson",
    "zipcode": "32132",
    "city": "Schweden - ☀",
    "color": "grün"
  },
  {
    "id": 9,
    "name": "Klaus",
    "lastname": "Klaussen",
    "zipcode": "43246",
    "city": "Hierach",
    "color": "grün"
  }
]
```

### POST /persons

```
$ curl -i -d '{"name":"Max", "lastname":"Mustermann","zipcode":"55555","city":"Musterstadt","color":5}' localhost:4000/persons
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 02 Feb 2026 12:32:09 GMT
Content-Length: 121

{
  "id": 10,
  "name": "Max",
  "lastname": "Mustermann",
  "zipcode": "55555",
  "city": "Musterstadt",
  "color": 5
}
```
