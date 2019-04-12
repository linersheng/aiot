package devstat

import (
	"fmt"
	"github.com/linersheng/aiot/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type DevStat struct {
	redisAddr   string `json:"RedisAddr"`
	redisPasswd string `json:"RedisPasswd"`
	heartTime   int    `json:"HeartTime"`
	port        int    `json:"Port"`
	countOnline int64
	redis       *Redis
}

type Message struct {
	Id string `json:"id"`
}

// func updateHeartbeat(id string) error {
// 	fmt.Println("test")
// 	return nil
// }

// func checkDevOnline(id []string) (bool, error) {
// 	fmt.Println("test")
// 	return true, nil
// }

var (
	app *DevStat = &DevStat{
		redisAddr:   "localhost:6379",
		redisPasswd: "",
		heartTime:   60,
		port:        4455,
	}
)

func newKeyString(id string) string {
	return fmt.Sprintf("%s:%010d", id, time.Now().Unix())
}

func updateHeartbeat(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	//r.Header
	fmt.Println("recv:", string(body))

	m := &Message{}
	if err := utils.JsonUnmarshal(body, m); err != nil {
		fmt.Fprintf(w, `{"status":"json error"}`)
		return
	}

	fmt.Println("id:", m.Id)

	fmt.Fprintf(w, `{"status":"OK"}`)

	if key := app.redis.TryGetKeyInDB(m.Id); key != "" {
		app.redis.SetExpired(key, app.heartTime*3)
		return
	}
	app.redis.AddKeyWithTimeout(newKeyString(m.Id), app.heartTime*3)
}

func SetConfig(c Config) bool {
	app.redisAddr = c.RedisAddr
	app.redisPasswd = c.RedisPasswd
	app.heartTime = c.HeartTime
	app.port = c.Port

	return true
}

func Run() {

	app.redis = &Redis{
		Addr:      app.redisAddr,
		Passwd:    app.redisPasswd,
		Subscribe: 0x03,
	}

	app.redis.InitRedis()
	http.HandleFunc("/heart", updateHeartbeat)
	addr := ":" + strconv.Itoa(app.port)
	http.ListenAndServe(addr, nil)
}
