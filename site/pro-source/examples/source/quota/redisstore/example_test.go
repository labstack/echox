// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package redisstore_test

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"echo.labstack.com/pro/quota/redisstore"
)

// docs:start ExampleNew
func ExampleNew() {
	client := redis.NewClient(redisstore.Options(&redis.Options{Addr: "localhost:6379"}))
	defer func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}()
	store, err := redisstore.New(client, redisstore.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(store.Remote())
	// Output: true
}

// docs:end ExampleNew

// docs:start ExampleMustNew
func ExampleMustNew() {
	client := redis.NewClient(redisstore.Options(nil))
	defer func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println(redisstore.MustNew(client, redisstore.Config{}).Validate())
	// Output: <nil>
}

// docs:end ExampleMustNew

// docs:start ExampleOptions
func ExampleOptions() {
	o := redisstore.Options(&redis.Options{Addr: "localhost:6379"})
	fmt.Println(o.ContextTimeoutEnabled, o.MaxRetries)
	// Output: true -1
}

// docs:end ExampleOptions

// docs:start ExampleFailoverOptions
func ExampleFailoverOptions() {
	o := redisstore.FailoverOptions(&redis.FailoverOptions{
		MasterName:    "primary",
		SentinelAddrs: []string{"localhost:26379"},
	})
	fmt.Println(o.MasterName, o.ContextTimeoutEnabled, o.MaxRetries)
	// Output: primary true -1
}

// docs:end ExampleFailoverOptions

// docs:start ExampleClusterOptions
func ExampleClusterOptions() {
	o := redisstore.ClusterOptions(&redis.ClusterOptions{Addrs: []string{
		"localhost:7001",
		"localhost:7002",
		"localhost:7003",
	}})
	fmt.Println(len(o.Addrs), o.ContextTimeoutEnabled, o.MaxRetries)
	// Output: 3 true -1
}

// docs:end ExampleClusterOptions
