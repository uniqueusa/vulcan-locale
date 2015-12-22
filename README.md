# Locale Middleware for Vulcan Proxy
The Locale Middleware for the Unique USA API Gateway will set headers on the requests before passing them on to their final desitnation.

These headers are:

    Accept-Language
    Accept-Currency

The values for the headers will be driven by the confiuration, which is a YAML file.

## Install
```
go get github.com/uniqueusa/vulcan-locale
```

## Usage
This presumes you have built new `vulcand` and `vctl` binaries per [the instructions](http://vulcanproxy.com/middlewares.html#example-auth-middleware). Basically, you should be able to add `github.com/skookum/vulcan-cors` to your registry and build your `vulcand` and `vctl` binaries.

1. Create a YAML file of your allowed hosts and methods:
```
http://google.com:
  - GET
  - POST
http://balls.com:
  - "*"
"*":
  - GET

```
(Notice that to allow anything use `"*"`. The quotes are necessary. Probably another caveat.)

2. Add the middleware
```
vctl cors upsert -id=locale_middleware-f someFrontend -configFile=yourYaml.yml --vulcan=http://yourvulcanhost
```
(`-id` can be whatever you want to call the instance of the middleware)

3. Make CORS enabled requests!

### Remove
```
vctl cors rm -id locale_middeware -f someFrontend --vulcan=http://yourvulcanhost
```

## Contributing
1. Write tests
2. Write code
3. Run tests until they pass
4. Run `codeclimate analyze` and fix suggestions
5. Issue PR
