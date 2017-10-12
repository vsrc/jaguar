# Jaguar
CRUD mock server
------

Free and open surced (MIT license) mock server with built-in persistent storage (database). 

## Purpose

When writing client side software we need something that will pretend to be back-end server side service and that something will idealy sit on local machine to avoid conectivity and latency issues while developing. Most of us usually write mock json files whcih we put somewhere on our computer and invoke from within our code. Beside slowing down our client side app we are doing work on, this approach is clumsy and not flexible. This way we can't mock create functionality or we need to do additional code that will mimic other back-end features like limit and offset. That mock code can sometimes be even source of wierd bugs as it dosen't truly represent back-end behaviour.

Ideally we need something that will be:
- easy to setup
- easy to run
- easy to add inital data
- cross platform/OS compatible
- faaaaast

This is where Jaguar kicks in.

## Features

As mentioned before Jaguar server have built in persistent storage solution which will save all information locally in same folder where is Jaguar running from. To run Jaguar one should only invoke from command line binary file for his/her/their operating system.

Once running Jaguar server will start listening on port defined in conf.json file and will provide following endpoints:

#### [GET] /
Returns all records in database with default offset=0 and limit = 100. At the begining it will return empty json as there is nothing stored. Default limit can be changed in conf.json. Offset and limit can also be changed on each request by adding query parameters, for example: `http://localhost:8765/?limit=10&offset=5`

#### [GET] /:id
Returns record with that id if it exists or empty json if it doesn't.

#### [POST] /
Accepts json formatted payload in body of request and saves it in database.

#### [POST] /bulk
Accepts an array of json objects and saves them as individual records. Useful if we want to add more than a few records if we want to test performance, filtering, sorting and other functionality in our client side software that requires large scale of data.

#### [PATCH] /:id
Accepts json formatted payload in body of request. If record with provided id is found, it will merge already existing data with data provided.

#### [PUT] /:id
Accepts json formatted payload in body of request. If record with provided id is found, it will completely replace already existing data with data provided.

#### [DELETE] /single/:id
If it finds record with that id, it will be deleted.

#### [DELETE] /bulk
Accepts an array of ids in a payload of body of request. For each record, if found in database it will be deleted. Useful for large number of records at once and for testing that functionality from client side if it is implemented.

## Download

## Technology used

This project is 100% written in golang. It uses two open source packages:
- [Gin](https://github.com/gin-gonic/gin)
- [Tiedot](https://github.com/HouzuoGuo/tiedot)
