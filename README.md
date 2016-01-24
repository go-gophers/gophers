# Gophers [![GoDoc](https://godoc.org/github.com/gophergala2016/gophers?status.svg)](https://godoc.org/github.com/gophergala2016/gophers)

Hacked during [Gopher Gala 2016](http://gophergala.com).

<img align="right" src="https://github.com/gophergala2016/gophers/wiki/logo.png" />

Gophers is a tool for API testing. It covers:
* unit testing of individual endpoints;
* functional testing of broader scenarios;
* generation of up-to-date examples for documentation from scenarios.

Gophers allows one to write tests scenarios in real programming languages, not in pesky UI.
Those languages are Go and Lua.

Go package contains a lot of helpers tailed just for that task. In particular, sometimes they
sacrifice idiomatic approach to brevity and simplicity of usage in test scenarios. For example,
many methods explicitly fail test or panic instead of returning error which should be checked manually.

For example this code can be used to
[create repository on Github via API](https://developer.github.com/v3/repos/#create)
and check result:
```go
// Client contains base URL with host, path prefix, headers and query parameters
// create new request with JSON body
req := Client.NewRequest(t, "POST", "/user/repos", jsons.Parse(`{"name": %q}`, repo))

// enable recording of request and response for documentation
req.EnableRecording("repo_create.apib")

// make request and check response status code
j := Client.Do(t, req, 201).JSON(t)

// check create repository
assert.Equal(t, jsons.Parse(`{"name": %q, "full_name": %q}`, repo, Login+"/"+repo),
	j.KeepFields("name", "full_name"))

// check repository is owned by authenticated user
assert.Equal(t, jsons.Parse(`{"login": %q}`, Login), j.Get("/owner").KeepFields("login"))
```

Running this scenario with `go test` and combining recorded request and response with
[API Blueprint template](examples/github-go/github.apib) will produce
[documentation with accurate examples](https://rawgit.com/gophergala2016/gophers/master/examples/github-go/github.html).

Lua bindings would allow making tests even simpler while using the whole power and speed of Go.
They are work-in-progress.


## Usage

Enable vendor experiment (`GO15VENDOREXPERIMENT=1`) and get package as usual:
```
go get github.com/gophergala2016/gophers
```

Then use it for writing your tests, see [examples](examples/) directory.


## Future work

After Gopher Gala development will happen at https://github.com/go-gophers/gophers. Plans include:

* allow to remove extra headers from requests and responses (Github, _why_ you send so much of them?)
* better ideomatic Lua bindings
* support for other test frameworks (`testing` wasn't the the best choice due to logging issues
  and panic handling)
* load testing
* fuzz testing (?)
* support for other protocols and API types


## License

Code is licensed under [MIT-style license](LICENSE).

Gopher artwork is taken from [gophericons](https://github.com/hackraft/gophericons).
Created by [Olga Shalakhina](https://www.facebook.com/olga.shalakhina), based on original work
by [Renée French](http://reneefrench.blogspot.com). Licensed under
[Creative Commons 3.0 Attributions](http://creativecommons.org/licenses/by/3.0/).
