package main

import (
	"testing"

	"github.com/ApiRequestLimiter"
)

func TestLimiterAgent_HandleRequest(t *testing.T) {
	ttt := 50
	user := "1234567890"
	re := make(chan bool)
	for i := 0; i < ttt; i++ {
		go request(i, t, re, user)
	}
	for i := 0; i < ttt; i++ {
		<-re
	}

}

func request(i int, t *testing.T, re chan bool, user string) {
	result, err := limiter.GLimiterAgent().HandleRequest(user, 1)
	re <- result
	t.Log(i, result, err)
}
