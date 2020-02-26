# qshr.tn
[![CircleCI](https://circleci.com/gh/VolticFroogo/QShrtn.svg?style=svg)](https://circleci.com/gh/VolticFroogo/QShrtn)

 The server component for an anonymous, open-source URL shortening service.

# API
## Terminology
Term | Description
---- | -----------
Code | A status code sent with all requests specifying the successfulness of a request
ID   | The path in which a redirect is from, `qshr.tn/{id}`
URL  | The URL in which a client will be redirected to if they send a `GET` request to `qshr.tn/{id}`

## Get redirect
The function that handles these requests can be found in [redirect.go](redirect/redirect.go) and the tests in [redirect_test.go](redirect/redirect_test.go).

### Request
To get a redirect via the API you must send a `GET` request with the path `/{id}/json`.

### Response
This endpoint will yield the following encoded as JSON:

Field | Type   | Description             | Included
----- | ------ | ----------------------- | ---------
code  | int    | Status code             | Always
url   | string | Redirected URL          | When the code is success
error | string | Fatal error description | When the code is internal server error

The status codes are as follows:

Code | Description
---- | -----------
0    | Success
1    | Not found
2    | Internal server error

Example success response:
```json
{
    "code": 0,
    "url":  "https://froogo.co.uk/"
}
```

Example not found response:
```json
{
    "code": 1
}
```

Example internal server error response (EOF caused by bad JSON request body):
```json
{
    "code":  2,
    "error": "EOF"
}
```

## Create redirect
The function that handles these requests can be found in [new.go](redirect/new.go) and the tests in [redirect_test.go](redirect/redirect_test.go).

### Request
To create a new redirect you must send a `POST` request to `/new/` with a JSON body.

The request body must be the following encoded as JSON:

Field | Type   | Description                | Required
----- | ------ | -------------------------- | ---------
url   | string | URL to redirect to         | Yes
id    | string | Path to be redirected from | No

If the ID is not included or empty, a random four character long ID will be generated.

Example request body with no id:
```json
{
    "url": "https://froogo.co.uk/"
}
```

Example request body with an id:
```json
{
    "url": "https://froogo.co.uk/",
    "id":  "frog"
}
```

### Response
A request to this endpoint will yield the following encoded as JSON:

Field | Type   | Description                | Included
----- | ------ | -------------------------- | ---------
code  | int    | Status code                | Always
id    | string | Path to be redirected from | When the code is success
error | string | Fatal error description    | When the code is internal server error

The status codes are as follows:

Code | Description
---- | -----------
0    | Success
1    | Internal server error
2    | Forbidden domain (URL contains hostname (qshr.tn))
3    | ID taken
4    | Invalid URL

Example success response:
```json
{
    "code": 0,
    "id":  "frog"
}
```

Example ID taken response:
```json
{
    "code": 3
}
```

Example internal server error response (EOF caused by bad JSON request body):
```json
{
    "code":  1,
    "error": "EOF"
}
```
