# Gophers [![GoDoc](https://godoc.org/github.com/go-gophers/gophers?status.svg)](https://godoc.org/github.com/go-gophers/gophers) [![Build Status](https://travis-ci.org/go-gophers/gophers.svg?branch=master)](https://travis-ci.org/go-gophers/gophers)

<img align="right" src="https://github.com/go-gophers/gophers/wiki/logo.png" />

Gophers is a tool for API testing. It covers:
* unit testing of individual endpoints;
* functional testing of broader scenarios;
* generation of up-to-date examples for documentation from scenarios.

> Note: For now it's focused on HTTP APIs. Support for other protocols is planned.

Gophers allows you to write test scenarios in full-power programming languages, not by using
limited pesky UI. Those languages are Go and (in the future) Lua.

Go package contains a lot of helpers tailored just for that task. In particular, sometimes they
sacrifice idiomatic approach for brevity and simplicity of usage in test scenarios. For example,
many methods explicitly fail test or panic instead of returning error which should be checked
in test manually.

For example this code can be used to
[create repository on Github via API](https://developer.github.com/v3/repos/#create)
and check result:
```go
// Client contains base URL with host, path prefix, default headers and query parameters
// t is *testing.T or compatible interface

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

// check repository exists via other API
j := Client.Get(t, "/repos/"+Login+"/"+repo, 200).JSON(t)
assert.Equal(t, jsons.Parse(`{"login": %q}`, Login), j.Get("/owner").KeepFields("login"))
```

Running this scenario with `go test` and combining recorded request and response with
[API Blueprint template](examples/testing/github/github.apib) will produce
[documentation with accurate and up-to-date examples](https://rawgit.com/go-gophers/gophers/master/examples/testing/github/github.html).

Lua bindings would allow making tests even simpler while using the whole power and speed of Go
networking stack. They are work-in-progress.


## Usage

Get the package as usual with Go 1.6+:
```
go get -u github.com/go-gophers/gophers
```

Then use it for writing your tests, see [examples](examples/) directory.


## Future work

First version was hacked during [Gopher Gala 2016](http://gophergala.com). Now development happens at
https://github.com/go-gophers/gophers. Our plans include:

* allow to remove extra headers from requests and responses for documentation (Github, _why_ you send so much of them?)
* better ideomatic Lua bindings (already drafted)
* fuzz testing (?)
* support for other protocols and API types
* mruby bindings (?)


## License

Code is licensed under [MIT-style license](LICENSE).

Gopher artwork is taken from [gophericons](https://github.com/hackraft/gophericons).
Created by [Olga Shalakhina](https://www.facebook.com/olga.shalakhina), based on original work
by [Renée French](http://reneefrench.blogspot.com). Licensed under
[Creative Commons 3.0 Attributions](http://creativecommons.org/licenses/by/3.0/).
