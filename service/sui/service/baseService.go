package service

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/sui"
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
	BlockApi *sui.ISuiAPI
}

func InitSuiService() *SuiService{
	once.Do(func() {
		c, err := client.Dial(SuiMainNetEndpoint)
		if err != nil {
			errorMsg := fmt.Sprintf("connect sui main net error:%v", err)
			panic(errorMsg)
		}
		blockSuiApi := sui.NewSuiClient(SuiMainNetEndpoint)
		instance = &SuiService{
			SuiApi: c,
			BlockApi: &blockSuiApi,
		}

	})
	return instance
}
