package part

import (
	"fmt"
	"testing"
)

func TestPart(t *testing.T) {
	// forever := make(chan int, 1)
	// ServerPart()
	// <-forever

	servers := []string{"server0", "server1", "server2"}
	s := ServerIndex(servers)
	fmt.Println(s)
}
