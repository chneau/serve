# serve
A VERY simple http server, WITH BASIC AUTH AS DEFAULT !  
Will setup an upload server and a static server.

### Install
```
go get -u -v github.com/chneau/serve/...
```

### Usage
```
Usage of serve:
  -noauth
        do not ask for auth
  -path string
        path to directory to serve (default ".")
  -port string
        port to listen on (default "8888")
  -pwd string
        password for auth
  -usr string
        username for auth
```

### Example
```
Username: *******
Password: *********
Serving files from  .
Listening on http://c:8888/serve
	http://localhost:8888/serve
Listening on http://c:8888/upload
	http://localhost:8888/upload

```
