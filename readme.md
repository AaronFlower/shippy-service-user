## Problem

We've added our password hashing functionality, and we set it as our password before saving
a new user. * You should never, ever store plain-text passwords *

Also, on authentication, we check against the hashed password.

Now, traditionally, we can securely authenticate a user against the database, however, we need
a mechanism in which we can do this across our user interfaces and distributed services. There
are many ways in which to do this, but the simplest solution I've come across, which we can use
across our services and web, is JWT.

- JWT vs OAuth

### JWT
JWT stands for JSON web tokens and is a distributed security protocol. Similar to OAuth.
The concept is simple, you use an algorithm to generate a unique hash for a user, 
which can be compared and validated against.

But not only that, the token itself can contain and be made up of our users' metadata.
In other words, their data can itself become a part of the token. So let's look at an example
of a JWT:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

The token is separated into three by .'s. Each segment has a significance. 

- The first segment is made up of some metadata about the token itself. Such as the type of token and the algorithm used to create the token. This allows clients to understand how to decode the token. 
- The second segment is made up of user-defined metadata. This can be your user's details, an expiration time, anything you wish. 
- The final segment is the verification signature, which is information on how to hash the token and what data to use.


### micro

使用 micro 提供的 api-gateway 来管理 `shippy` 命名空间中的服务。 默认启动 `go.micro.api` 可以发现我们的服务
并且代理我们的服务。

```bash
$ docker run -p 8080:8080 \ 
    -e MICRO_REGISTRY=mdns \
    microhq/micro api \
    --handler=rpc \
    --address=:8080 \
    --namespace=shippy
```

Register a new user, use httpie
```bash
shippy-service-user master ✗ 1h2m △ ◒ ➜ http POST :8080/rpc <<<'{ "service": "shippy.auth", "method": "Auth.Create", "request":  { "email": "Bar@gmail.com", "password": "Foo" } }'
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PATCH, PUT, DELETE
Content-Length: 152
Content-Type: application/json
Date: Thu, 20 Dec 2018 09:11:44 GMT

{
    "user": {
        "email": "Bar@gmail.com",
        "id": "606ec98d-aa6b-47ad-99d1-f9ab49cb4a0b",
        "password": "$2a$10$f/rAalzrZ0hIZxEn5KlRcuoz9/9g/3TB/8kpHiO.i6.mMMPINNQDK"
    }
}
```

Authenticate a new user.

```bash
shippy-service-user master ✗ 1h2m △ ◒ ➜ http POST :8080/rpc <<<'{ "service": "shippy.auth", "method": "Auth.Auth", "request":  { "email": "Bar@gmail.com", "password": "Foo" } }'
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PATCH, PUT, DELETE
Content-Length: 353
Content-Type: application/json
Date: Thu, 20 Dec 2018 09:11:53 GMT

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7ImlkIjoiNjA2ZWM5OGQtYWE2Yi00N2FkLTk5ZDEtZjlhYjQ5Y2I0YTBiIiwiZW1haWwiOiJCYXJAZ21haWwuY29tIiwicGFzc3dvcmQiOiIkMmEkMTAkZi9yQWFsenJaMGhJWnhFbjVLbFJjdW96OS85Zy8zVEIvOGtwSGlPLmk2Lm1NTVBJTk5RREsifSwiZXhwIjoxNTQ1NTU2MzEzLCJpc3MiOiJnby5taWNyby5zcnYudXNlciJ9.FpxJ-8N6EURmv3eVHN1ADpeVzHrxowO_si87aOBgOnY"
}
```

Create a consignment
```bash
http POST :8080/rpc <<< '{"service": "shippy.consignment", "method": "ShippingService.CreateConsignment", "request": {"description": "This is a test", "weight": 500, "containers": [] } }'

HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PATCH, PUT, DELETE
Content-Length: 101
Content-Type: application/json
Date: Thu, 20 Dec 2018 10:55:05 GMT

{
    "consignment": {
        "description": "This is a test",
        "vessel_id": "vessel0001",
        "weight": 500
    },
    "created": true
}
```
