### IMPLEMENTATION
The code consists of a weather module, gateway, and a http helper.   

The weather module contains the request handlers for retrieving the weather, temperatures and windspeeds for a given range of dates. It uses the gateway implementation to requests the data and it does that using goroutines, so it could do several request in parallel.

The gateway is a generic code to connect to the temperatures or speeds api. It consists of a get request and basic request error handling.

The http helper contains the configuration for http client, http structs, and request implementation used by the entire application.

### TESTS
The provided tests coverages 94.0% of the code. There're two files for that, `weather_test.go` and`gateway_test.go`.

### HOW TO RUN
application: `docker-compose up`  
unit tests: `docker-compose -f docker-compose.test.yml up`


