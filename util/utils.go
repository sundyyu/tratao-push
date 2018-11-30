package util

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
)

func LogErrorM(err interface{}, msg string) {
	if err != nil {
		log.Printf("%s: %s \n", msg, err)
	}
}

func LogError(err ...interface{}) {
	if err != nil {
		log.Println(err)
	}
}

func LogInfo(msg ...interface{}) {
	log.Println(msg)
}

func LogInfoF(format string, msg ...interface{}) {
	log.Printf(format+" \n", msg)
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func RedisResult2Int64(result interface{}) (int64, error) {
	if val, ok := result.(string); ok {
		val_64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return -1, err
		}
		return val_64, nil
	}
	return -1, errors.New("interface is not string")
}

func RedisResult2Float64(result interface{}) (float64, error) {
	if val, ok := result.(string); ok {
		val_64, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return -1, err
		}
		return val_64, nil
	}
	return -1, errors.New("interface is not string")
}

func RedisResult2Bool(result interface{}) (bool, error) {
	if val, ok := result.(string); ok {
		val_bool, err := strconv.ParseBool(val)
		if err != nil {
			return false, err
		}
		return val_bool, nil
	}
	return false, errors.New("interface is not string")
}

func MapResult(status int, data interface{}, msg interface{}) map[string]interface{} {
	result := make(map[string]interface{}, 10)
	result["status"] = status
	result["data"] = data
	result["msg"] = msg

	return result
}

func ReadChan(r chan int) {
	<-r
}
func WriteChan(w chan int) {
	w <- 1
}

func ReadChanNum(r chan int) int {
	return <-r
}
func WriteChanNum(w chan int, n int) {
	w <- n
}

func ToJson(v interface{}) string {
	if byt, err := json.Marshal(v); err == nil {
		return string(byt)
	}
	return ""
}
