package main

import (
	"time"
	"net/http"
	"encoding/json"
	"fmt"
	"net"
	"github.com/gomodule/redigo/redis"
)

const (
	SERVER_PORT = "3000"
	REDIS_HOST = "localhost"
	REDIS_PORT = "6379"
)

var redisConn redis.Conn

func redisSet(requestParams map[string]interface{}, w http.ResponseWriter) {
	// handle params
	key, isExist := requestParams["key"]
	if !isExist {
		println("KEY param not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, isExist := requestParams["value"]
	if !isExist {
		println("KEY param not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// make request ot redis
	_, err := redisConn.Do("SET", key, value)
	if err != nil {
		println("redis command error: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func redisDel(requestParams map[string]interface{}, w http.ResponseWriter) {
	// handle params
	key, isExist := requestParams["key"]
	if !isExist {
		println("KEY param not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// make request to redis
	_, err := redisConn.Do("DEL", key)
	if err != nil {
		println("redis command error: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func redisGet(requestParams map[string]interface{}, w http.ResponseWriter) {
	// handle params
	key, isExist := requestParams["key"]
	if !isExist {
		println("KEY param not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// make request ot redis
	value, err := redisConn.Do("GET", key)
	if err != nil {
		println("redis command error: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// it means no such key in redis
	if value == nil {
		return
	}

	// handle response from redis
	valueUint8, ok := value.([]uint8)
	if !ok && value != nil {
		println("value is not UINT8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Printf("%#v %T\n", value, value)
		return
	}
	w.Write([]byte(valueUint8))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT,DELETE,POST,OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		println("options")
		w.WriteHeader(http.StatusOK)
		return
	}

	var requestParams map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		println("json decode failed: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("params: %#v\n", requestParams)

	if r.Method == "PUT" {
		println("PUT. It means client wants to set key/value")
		redisSet(requestParams, w)
	} else if r.Method == "DELETE" {
		println("DELETE. It means client wants to delete key/value")
		redisDel(requestParams, w)
	} else if r.Method == "POST" {
		println("POST. It means client wants to get value by the key")
		redisGet(requestParams, w)
	}
}

func router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)
	return mux
}

func redisInit() error {
	var timeout, err = time.ParseDuration("2s")
	if err != nil {
		println("time parse error: "+err.Error())
		return err
	}

	var tcpConn net.Conn

	tcpConn, err = net.Dial("tcp", REDIS_HOST+":"+REDIS_PORT)

	if err != nil {
		println("tcp connection error: "+err.Error())
		return err
	}

	redisConn = redis.NewConn(
		tcpConn,
		timeout,
		timeout,
	)
	return nil
}

func main() {
	err := redisInit()
	if err != nil {
		return
	}
	defer redisConn.Close()

	mux := router()
	println("starting server at :"+SERVER_PORT)
	http.ListenAndServe(":"+SERVER_PORT, mux)
	println("Порт " + SERVER_PORT + " занят")
}