package myokx

import (
	"fmt"
	"testing"
)

var c *Client

func TestMain(m *testing.M) {
	fmt.Println("Testing ================= begin")

	// okex client:
	config := NewConfig("config.json")
	c = NewClient(config)

	m.Run()

	fmt.Println("Testing ================== end.")
}
