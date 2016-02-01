# golangTest
=========
Go code must be kept inside a workspace. A workspace is a directory hierarchy with three directories at its root:

1. src contains Go source files organized into packages (one package per directory),
2. pkg contains package objects, and
3. bin contains executable commands.

The src subdirectory will contain my git repository: git@github.com:mujibur-rahman/golangTest.git

The GOPATH environment variable
========
The GOPATH environment variable specifies the location of your workspace.

* $ mkdir $HOME/work
* $ export GOPATH=$HOME/work

Now cd to work/src directory and the use following comand

* git clone  git@github.com:mujibur-rahman/golangTest.git

Add new third part packages which is needed to run.
* go get github.com/gin-gonic/gin
* go get github.com/dgrijalva/jwt-go

Run application
=======
To run the application run the follow commands
* cd golangTest
* go run main.go

To use flag while running the `main.go`

	* go run main.go [It will start default port: 8000, ginMode: release, logToFile:true and useGinLogger: false
	* go run main.go -port=9000 -logToFile=false [It will start port 9000 and logToFile=false(It supports to write log on stdout)]


it will start go web server. It will log all to `logs/web-server.log` command to see that: `tail -f  logs/web-server.log` from current directory


Testing
=======
Open a new command line.

1. export GOPATH=$HOME/work
2. cd cd golangTest/auth
3. go test -v

Applications specific
===========
Following commands needed to run and test the apps.
I did not use any web html template, but everything I tested in command line tool by using `curl`
I used GIN framework for RESTFUL API supports

1. curl http://localhost:8000/

Response: {"message":"Its an entry point :)"}

2. curl http://localhost:8000/token

Response: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6Im11amlidXJUZXN0IiwiZXhwIjoxNDU0MjU3MzcyfQ.4RnySbetK3SBFeSjJ7f6FHF7Ji63pCdES8dZAVXGVRE"}

Now you need to copy the token and paste while testing the auth and others api. Replace `paste-me-here-token-keys` with the token

3. curl --form authkey=abc123456abc -H "Authorization: bearer paste-me-here-token-keys" http://localhost:8000/auth
following things happend if:

### expired keys:
>$ `curl --form authkey=abc123456abc -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6Im11amlidXJUZXN0IiwiZXhwIjoxNDU0MTc4ODIyfQ.n-BrOHIX2ns--X52BA2eF1PyeD77-SPsUGpyftEAD_w" http://localhost:8000/auth`

#Log
>2016-01-31T23:25:58GMT 127.0.0.1 401 316.142µs "POST /auth HTTP/1.1" -1

>Error #01: token is expired

### No token provided
>$ `curl --form authkey=abc123456abc http://localhost:8000/auth`

#Log
>2016-01-31T23:27:22GMT 127.0.0.1 401 734.016µs "POST /auth HTTP/1.1" -1

> Error #01: no token present in request

### Wrong auth key will give 401
>$ `curl -v --form authkey=abc123456abc1 -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6Im11amlidXJUZXN0IiwiZXhwIjoxNDU0MjU3MzcyfQ.4RnySbetK3SBFeSjJ7f6FHF7Ji63pCdES8dZAVXGVRE" http://localhost:8000/auth`

>Response: "Wrong key"

### Without Token
> $ `curl -v http://localhost:8000/user`

>Response: 401

### With token return correct result. If matched then returned user data
>$ `curl -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6Im11amlidXJUZXN0IiwiZXhwIjoxNDU0MjU3MzcyfQ.4RnySbetK3SBFeSjJ7f6FHF7Ji63pCdES8dZAVXGVRE" http://localhost:8000/user`

> Response: {"user":{"country":"malaysia","id":9001,"mobile":"+019378646","username":"Mujibur"}}

### For todo list I added a find(`id`) handler
> $ `curl -H "Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6Im11amlidXJUZXN0IiwiZXhwIjoxNDU0MjU3MzcyfQ.4RnySbetK3SBFeSjJ7f6FHF7Ji63pCdES8dZAVXGVRE" http://localhost:8000/user/9001`

> Response: {"id":9001,"username":"Mujibur","country":"malaysia","mobile":"+019378646"}

#### DB Accessor (According to docs I did not make it under auth middleware)
> $ `curl http://localhost:8000/dbaccessor`

> {"object":{"Databases":[{"Tables":[{"Types":[{"country":"string","id":"int","mobile":"string","name":"string"}],"name":"user"}],"host":"localhost","name":"golang","pass":"root","port":3306,"type":"mysql","user":"root"}],"buffer_size":10}}

Documentation
=========
$ `godoc golangTest/auth`
