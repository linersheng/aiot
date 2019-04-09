# devstat 设备是否在线，上线时间、下线时间、在线时长等
***
* 通过设备心跳来判断设备是否在线，超过指定的时间没收心跳就判断设备离线
* 设备离线后，存储一条日志，包括设备上线时间、下线时间
* 使用redis存储在线设备的信息
* 使用ElasticSearch存储设备登录日志信息

**Config**
```json
{
    "RedisAddr":"127.0.0.1:6379",
    "HeartTime":60,
    "TimeoutCount":3
}
```

**redis key 设计**
>
> | key | value |
> |:--:|:--:|
> |flag:devid:logintime | lastheartbeattime |
> | dev:1234567890:1551234567| 1551234577 |
>  
> flag : key的开头标志   
> devid : 设备的唯一id   
> login : 设备第一次心跳时间   
> lastheartbeattime : 设备最后一次心跳时间
>
> 1 采用dev开头，用于统计当前在线设备的总数   
> 2 用redis expire做超时，在key超时后，把设备信息写入日志

**ElasticSearch**  
```json
{
    "devid":"1234567890",
    "logintime":1551234567,
    "lasttime":1551234577,
    "online":10
}
```

**API**  
> UpdateHeartbeat(id string) 
> ```go
>UpdateHeartbeat(id string) error {
>
>}
>```
>
>CheckOnline(id string)
>```go
>CheckOnline(id string) (bool, error) {
>
>}
>````

**Docker**
```shell
docker pull redis
docker run -d -p 6379:6379 --name="myredis" -v /root/redis/redis.conf:/etc/redis/redis.conf -v /root/redis/data:/data redis  redis-server /etc/redis/redis.conf --appendonly yes

docker pull elasticsearch:6.70
docker run -d -p 9200:9200 --name="myelasticsearch" elasticsearch:6.7.0
```
