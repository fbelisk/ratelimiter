# ratelimiter

This package provides a Golang lock-free implementation of the token-bucket rate limit algorithm. 

```
func main() {
	tb := New(100, 1, 1*time.Second) 

	for i := 0; i < 10; i++ {
		tokens, waitTime := tb.TakeWait(2, 5*time.Second)
		fmt.Println(tokens, waitTime)
	}

	// Output:
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
	//2 2.000000001s
}

```
