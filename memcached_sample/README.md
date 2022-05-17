# Elasticache for Memcached

## Memcached Go Library 이용하기 

### Local Memcached with Docker

- DOCKER를 이용하여 Memcached 실행하기. 

```py
docker run -d -p 11211:11211 memcached
```

### go module 가져오기 

- 모듈 초기화 하기 
  
```py
go mod init github.com/schooldevops/memcached
```

- memcached 라이브러리 획득하기 
  
```py
go get github.com/bradfitz/gomemcache/memcache
```

### 샘플 코드 작성 

```go
package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// Memcached에 커넥션을 연결한다. 
	mc := memcache.New("localhost:11211")

	// key/value 을 Memcached에 저장한다. 
	mc.Set(&memcache.Item{Key: "greeting", Value: []byte("Hello World")})

	// key를 이용하여 값을 조회한다. 결과값은 byte이다.
	item, err := mc.Get("greeting")

	if err != nil {
		panic("Error go with Panic")
	}

	// key, value, flags, expiration 값을 각각 조회한다.
	fmt.Println("Hello Greeting Key: ", item.Key)
	fmt.Println("Hello Greeting Value: ", string(item.Value))
	fmt.Println("Hello Greeting Flags: ", item.Flags)
	fmt.Println("Hello Greeting Expiration: ", item.Expiration)

}
```

- memcache.New: 
  - memcached 커넥션을 생성한다. 
- mc.Set:
  - 키/값 을 저장한다. 
- mc.Get:
  - 키에 대해서 값을 조회한다. 
  - item으로 반환되며, (키, 값, 플래그, 유효시간) 을 가지고 있다. 

## 관련 라이브러리 

- https://pkg.go.dev/github.com/bradfitz/gomemcache/memcache
- 상기 라이브러리를 활용하여 추가 작업을 수행한다. 


