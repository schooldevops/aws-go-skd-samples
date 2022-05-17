package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// Connect to Memcached
	mc := memcache.New("localhost:11211")

	// Set key/Value to Memcached
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
