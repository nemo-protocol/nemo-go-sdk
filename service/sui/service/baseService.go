package service

import (
	"fmt"
	"github.com/coming-chat/go-sui/v2/client"
	"sync"
)

var (
	instance    *SuiService
	once sync.Once
	SuiMainNetEndpoint = "https://fullnode.mainnet.sui.io"
)

type SuiService struct {
	SuiApi *client.Client
}

func InitSuiService() *SuiService{
	once.Do(func() {
		c, err := client.Dial(SuiMainNetEndpoint)
		if err != nil {
			errorMsg := fmt.Sprintf("connect sui main net error:%v", err)
			panic(errorMsg)
		}
		instance = &SuiService{
			SuiApi: c,
		}
	})
	return instance
}
