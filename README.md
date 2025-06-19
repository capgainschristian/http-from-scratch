# http-from-scratch

Codecrafter's "Build your own HTTP server" course.

To start the server:

```
go run app/main.go
```

Some curl commands to run against the server:

```
curl -v http://localhost:4221
```

```
curl -v http://localhost:4221/abcdefg
```

```
curl -v http://localhost:4221/echo/abc
```

```
curl -v --header "User-Agent: foobar/1.2.3" http://localhost:4221/user-agent
```

```
echo -n 'Hello, World!' > /tmp/foo
curl -i http://localhost:4221/files/foo
```

```
curl -v --data "12345" -H "Content-Type: application/octet-stream" http://localhost:4221/files/file_123
```

Check out my backend API server here: https://github.com/capgainschristian/go-backend-api
