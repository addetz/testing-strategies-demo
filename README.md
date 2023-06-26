# testing-strategies-demo
This repo contains the demo code for the talk `Comprehensive testing strategies for modern microservice architectures`. Slides are also available [here](./slides/).

It contains code for the conference talk server, which exposes the following endpoints: 
```
GET /events
GET /events/{id}
GET /events/{id}?day=X
```
For your convenience, this repo contains a Postman collection with the requests you can make to the server. See [`Conference_Talks.postman_collection.json`](./Conference_Talks.postman_collection.json).

## Run tests 
Run unit tests: 
```
$ go test ./... -v
```

Run integration tests:
```
$ INTEGRATION=true go test ./... -v
```

Run E2E tests: 
```
$ E2E=true go test ./... -v
```

## Run in Docker
The server image has been pushed to a public repo `classicaddetz/conf-talks-server`.

You can run it using: 
```
$ $ docker run -dt -p 8000:8000/tcp classicaddetz/conf-talks-server
```

Alternatively, you can build your own and run it:
```
$ docker build -f Dockerfile -t conf-talks-server .
$ docker run -dt -p 8000:8000/tcp conf-talks-server
```

## Push image
You can tag and push the image to your own repo as well: 

```
$ docker tag conf-talks-server username/conf-talks-server
$ docker push username/conf-talks-server
```
