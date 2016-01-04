# Locale Middleware for Vulcan Proxy
The Locale Middleware for the Unique USA API Gateway will set headers on the requests before passing them on to their final desitnation.

These headers are:

    Accept-Language
    Accept-Currency

The values for the headers will be driven by the configuration, which is a YAML file. These values will be mapped by domain, and that domain is expected to be in the `Origin` header of the request

## Install
```
go get github.com/uniqueusa/vulcan-locale
```

## Usage
This presumes you have built new `vulcand` and `vctl` binaries per [the instructions](http://vulcanproxy.com/middlewares.html#example-auth-middleware). Basically, you should be able to add `github.com/skookum/vulcan-cors` to your registry and build your `vulcand` and `vctl` binaries.

1. Create a YAML file of your domains with their locales and currencies:
```
esalerugs.com:
  locales:
    - en_US
  currencies:
    - usd
irugs.ch:
  locales:
    - en_GB
    - de_DE
    - fr_FR
    - nl_NL
  currencies:
    - eur
    - franc
irugs.sk:
  locales:
    - sk_SK
  currencies:
    - eur
irugs.ca:
  locales:
    - en_US
  currencies:
    - cad
```

2. Add the middleware
```
vctl locale upsert -id=locale_middleware-f someFrontend -configFile=yourYaml.yml --vulcan=http://yourvulcanhost
```
(`-id` can be whatever you want to call the instance of the middleware)

3. Make requests!

### Remove
```
vctl locale rm -id locale_middeware -f someFrontend --vulcan=http://yourvulcanhost
```

## Contributing
1. Write tests
2. Write code
3. Run tests until they pass
4. Run `codeclimate analyze` and fix suggestions
5. Issue PR
