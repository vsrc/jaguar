# Jaguar
superfast, easy to use mock server
------

Have you ever tried to implement front-end feature only to find out that back-end team hasn't implemented it yet? Have you ever tried to test the performance of your client-side software only to find out your database have only 3 records and your DBA is on vacation for next 14 days?

Meet Jaguar - free and open source (MIT license) mock server with out-of-the-box Create, Read, Update and Delete functionality.

But wait! There is more! If you download now, you will also get:

## Built-in zero setup database without structure 

Since the underlying database is a NoSql type of database, there is no need to do any database tables/structure setup. It accepts any data as long as it is in JSON format.

## Purpose

When writing client-side software we need something that will pretend to be back-end server-side service and that something will ideally sit on a local machine to avoid connectivity and latency issues while developing. Most of us usually write mock JSON files which we put somewhere on our computer and invoke from within our code. Besides slowing down our client-side app we are doing work on, this approach is clumsy and not flexible. This way we can't mock create functionality or we need to do additional code that will mimic other back-end features like limit and offset. That mock code can sometimes even be a source of weird bugs as it doesn't truly represent back-end behavior.

Ideally, we need something that will be:
- easy to setup
- easy to run
- easy to add initial data
- cross-platform/OS compatible
- faaaaast

This is where Jaguar kicks in.

### Example 1

Your back-end provides a response on endpoint `/hero` like this:
```json
{
  "alias": "Batman",
  "name": "Bruce"
}
```
But you need urgently implement in your web and mobile app surname as well. Back-end developer left for a lunch and won't be back for at least 2 hours. What will you do?

1. Run jaguar
2. Do a post request
  ```json
    {
      "alias": "Batman",
      "name": "Bruce",
      "surname": "Wayne"
    }
  ```
3. Do a get request and be ready to be served with your data while your back-end friend is still waiting for a waiter to pick up his/her/their lunch order.

### Example 2

Your back-end developer returned from lunch but he/she/they went now to get a coffee. He/she/they won't be back for at least 45 minutes. You need to test your web/mobile app performance with a large amount of data. Your database has only three records. What will you do?

1. Copy those 3 results
2. Paste them into your favorite or at least 2nd favorite text editor
3. Continue pasting until in front of your eyes 3 objects become 6, 6 becomes 12, 12 becomes 24...
4. Once you are happy, copy all objects
5. Run Jaguar
6. Do a post request on /bulk endpoint with all your copied objects
7. Do a get request and be ready to be served with your data while your back-end friend is still explaining to the waiter how to write his/her/their name on a coffee cup

## Features

As mentioned before Jaguar server has built-in persistent storage solution which will save all information locally in the same folder where is Jaguar running from. To run Jaguar one should only invoke from command line binary file for his/her/their operating system.

Once running Jaguar server will start listening on the port defined in conf.json file and will provide following endpoints:

#### [GET] /
Returns all records in a database with default offset=0 and limit = 100. At the beginning, it will return empty JSON as there is nothing stored. Default limit can be changed in conf.json. Offset and limit can also be changed on each request by adding query parameters, for example: `http://localhost:8765/?limit=10&offset=5`

#### [GET] /:id
Returns record with that id if it exists or empty json if it doesn't.

#### [POST] /
Accepts json formatted payload in body of request and saves it in database.

#### [POST] /bulk
Accepts an array of JSON objects and saves them as individual records. Useful if we want to add more than a few records if we want to test performance, filtering, sorting and other functionality in our client-side software that requires a large scale of data.

#### [PATCH] /:id
Accepts JSON formatted payload in the body of the request. If a record with provided id is found, it will merge already existing data with data provided.

#### [PUT] /:id
Accepts JSON formatted payload in the body of the request. If a record with provided id is found, it will completely replace already existing data with data provided.

#### [DELETE] /single/:id
If it finds record with that id, it will be deleted.

#### [DELETE] /bulk
Accepts an array of ids in a payload of the body of the request. For each record, if found in the database it will be deleted. Useful for a large number of records at once and for testing that functionality from client side if it is implemented.

## Can I use it for more permanent purposes?
Like as a backend for a website, blog, or my app? 

Probably yes. But since this project is done in a few hours there is no guarantee how it will behave in a more permanent role.

## Download
Binary file with the latest version compiled for your system can be downloaded from [releases page](https://github.com/vsrc/jaguar/releases)

## Who is using this project?
- me

_(end of list)_

If you would like to be mentioned here, let me know. :)

## Technology used

This project is 100% written in golang. It uses two open source packages:
- [Gin](https://github.com/gin-gonic/gin)
- [Tiedot](https://github.com/HouzuoGuo/tiedot)

## Code
As already mentioned above, this project is open sourced under MIT license, so you can easily read and examine the code, borrow, edit, rewrite, use and do whatever you want with it.

Code is consisting of just one `.go` file with under 300 lines of code, 257 lines of code to be precise. It consists of 9 functions. One of those 9 functions is the main bootstrap function, rest are endpoint handlers - 8 functions for 8 endpoints.
