To run the forecast service, check out the repo and run the main file:

```
git clone https://github.com/avrekhm/forecast.git
cd forecast/cmd
go build main.go
./main
```

Test the server with curl by varying the latitude and longitude parameters, e.g.

```
curl 'http://localhost:8080/forecast/44.3864,-73.5163'
```

TODOS:

* Better / more descriptive error handling: the server currently only returns the HTTP status code on error, with an empty body
* Deal with temperature scales other than F. The NWS API allows toggling between US and SI units via the `units=us|si` parameter
* Tests
* Better logging
