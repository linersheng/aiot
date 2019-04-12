package devstat

type Config struct {
	RedisAddr string `json:"RedisAddr"`
	RedisPasswd string `json:"RedisAddr"`
	HeartTime int    `json:"HeartTime"`
	Port      int    `json:"Port"`
}
