# 多用户限流器

## 调用方法

1. 调用GLimiterAgent()接口中的HandleRequest()函数.

2. GLimiterAgent()先初始化redis连接池及参数单例.

3. 函数请求参数为用户uid以及请求的次数numRequest（默认是1，小于maxPermits), 返回true/false以及error。

## 实现方法

1. 令牌桶思路.

2. lua脚本操作redis增加性能.

3. 通过conf设置maxPermits和Rate, 通过请求的时间差按照rate增加permits, 但不超过maxPermits.

4. 在redis中的key = "Limiter_uid", hash value为上次的请求时间LastNanoSec以及桶内的令牌数目CurrPermits.

5. 当permits有增加, 更新LastNanoSec到redis; 当CurrPermits数大于请求数, 返回true并更新CurrPermits到redis, 否则返回false.

6. 当用户第一次请求时初始化lastNanoSec为当前时间, currPermits = maxPermits-1, 并直接返回true.

7. 保证事务原子性（悲观锁）, 在并发的情况下先获得锁, 再对令牌桶进行操作, 最后解锁.

## 性能

1. 当单个用户并发请求时, QTS = 400

2. 当多个用户并发请求时, QTS与并发分散程度成正相关, 与并发数成反相关.

## TODO

1. 在redis中对于key的有效期设定.

2. 对于同一个用户的并发请求, 可以通过设置numRequest合并成一次函数调用做为优化.