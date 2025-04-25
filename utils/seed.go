package utils

import (
	"math/rand"
	"sync"
	"time"
)

func GetRandomValueFromSyncMap(servMap *sync.Map) (interface{}, bool) {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 收集所有键
	var keys []interface{}
	servMap.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})

	// 如果没有键，返回false
	if len(keys) == 0 {
		return nil, false
	}

	// 随机选择一个键
	randomKey := keys[rand.Intn(len(keys))]

	// 获取对应的值
	value, ok := servMap.Load(randomKey)
	return value, ok
}
