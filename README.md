# qshr.tn
[![CircleCI](https://circleci.com/gh/VolticFroogo/QShrtn.svg?style=svg)](https://circleci.com/gh/VolticFroogo/QShrtn)

 The server component for an anonymous, open-source URL shortening service.

## README Index
- [API](#api)
	- [Terminology](#terminology)
	- [Get redirect](#get-redirect)
		- [Request](#request)
		- [Response](#response)
	- [Create redirect](#create-redirect)
	    - [Request](#request-1)
    	- [Response](#response-1)
- [How to host](#how-to-host)
    - [Installing Docker (and Docker Compose)](#installing-docker-and-docker-compose)
    - [Cloning the project](#cloning-the-project)
    - [Setting up environment variables](#setting-up-environment-variables)
        - [production.env](#productionenv)
        - [db.env](#dbenv)
    - [Updating](#updating)

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
url   | string | URL to redirect to      | When the code is success
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

# How to host
This will assume you're using a Ubuntu 18.04 server; this isn't necessary.

As we will be using Docker to run this application, any Docker capable OS is fine.

Note: we'll be using the [DB Docker Compose file](docker-compose.db.yaml) for this as it fully encapsulates everything to make hosting easier.
This isn't used on the official qshr.tn system as replica sets and load balancers are used to ensure redundancy.

## Installing Docker (and Docker Compose)
I won't go over this as there are hundreds of guides on this online, but I will link what I recommend for Ubuntu 18.04.
Keep in mind, this is operating system specific and you may need a different guide for your case.
Just Google "install docker YOUR OS", and find a trusted guide.

[Install Docker on Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-18-04)

[Install Docker Compose on Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-18-04)

## Cloning the project
Assuming you have Git installed, go the directory in which you want to install qshr.tn, and execute this command:

`git clone https://github.com/VolticFroogo/QShrtn`

If you don't have Git installed, [here's a guide for Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-install-git-on-ubuntu-18-04).

## Setting up environment variables
Two environment files are used which you will need to create and configure: `production.env` and `db.env`.

These files must be in the root directory of the project.

### production.env
Name   | Description
------ | -----------
DB     | [A MongoDB connection URI](https://docs.mongodb.com/manual/reference/connection-string/)

Note: the `root:password` in this string must match the `MONGO_INITDB_ROOT_USERNAME` and `MONGO_INITDB_ROOT_PASSWORD` in the [db.env](#dbenv) config.

Example file:
```
DB=mongodb://root:password@mongo/?authSource=admin&appname=qshrtn&ssl=false
```

### db.env
Name                       | Description
-------------------------- | -----------
MONGO_INITDB_ROOT_USERNAME | The root username
MONGO_INITDB_ROOT_PASSWORD | The root password

Example file:
```
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=password
```

## Building and running
Change directory to the [directory that you cloned](#cloning-the-project) earlier.

Then simply execute this command:

`docker-compose -f docker-compose.db.yaml up --build`

As the [DB Docker Compose file](docker-compose.db.yaml) states the application should restart unless stopped,
this application will run after server restarts, updates, crashes, etc.

## Customise site
To customise the site, simply edit any files in the `static/` directory and rebuild:

`docker-compose -f docker-compose.db.yaml up --build`

Keep in mind, although everything by default uses a CDN for increased performance, all js, css, and img files are statically hosted.

So you can change them to local links like `/css/main.min.css` with no issues.

## Updating
If you ever want to update the application, follow this:

Change directory to the [directory that you cloned](#cloning-the-project) earlier.

Download the updates:

`git pull`

Rebuild and run:

`docker-compose -f docker-compose.db.yaml up --build`
