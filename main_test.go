package myokx

import (
	"fmt"
	"testing"
)

var c *Client

var maria *MyMariaDBClass

func TestMain(m *testing.M) {
	fmt.Println("Testing ================= begin")

	// okex client:
	config := NewConfig("config.json")
	c = NewClient(config)

	// mariadb:
	maria = &MyMariaDBClass{}
	maria.Init("root", "rF111222k", "127.0.0.1", "3306", "testdb")

	m.Run()

	fmt.Println("Testing ================== end.")
}
