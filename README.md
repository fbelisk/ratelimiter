# ratelimiter

基于令牌桶算法的限流器，使用lock-free实现

### 示例
```
func main() {
	tb := New(10000, 1000, 1*time.Second) 

	for i := 0; i < 10; i++ {
		tokens, waitTime := tb.TakeWait(200, 1*time.Second)
	}
}

```
### 初始化
初始化入参分别是：令牌桶容量、单位间隔新增令牌数、间隔时长
```
tb := New(10000, 1000, 1*time.Second)
```
### Take
向限流器申请指定数量的请求令牌，非阻塞

响应为实际可用令牌数量
```
tokens := tb.Take(10)
```

### TakeWait
向限流器申请指定数量的请求许可，如果令牌桶内令牌数量不足，则阻塞，直到数量足够

入参分别为请求令牌数量和最长等待时间

响应为实际可用令牌数量和阻塞等待的时间
```
tokens, waitTime := tb.TakeWait(200, 1*time.Second)
```
