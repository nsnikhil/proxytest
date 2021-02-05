## ProxyTest

#### Setup
```
make copy-config
make test
```

#### Build
```
make build
```

#### Local Run
```
make local-http-serve
```

#### Sample Request
```
curl --location --request GET 'http://localhost:8089/proxy?client-id=1234&method=GET' \
--header 'client-id: 1234' \
--header 'method: GET' \
--data-raw '{
    "url":"http://localhost:8081/random",
    "headers":{"A": ["1"]},
    "body":{"B": 2}
}'
```
