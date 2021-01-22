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
http://localhost:8080/proxy?client-id=1234&url=http://localhost:8081/random&method=GET&headers={"A": ["1"]}&body={"B": 2}
```
