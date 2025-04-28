package service

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/nemo-protocol/nemo-go-sdk/utils"
	"sync"
)

var (
	instance    *SuiService
	onlyEndpointInstance *SuiService
	onlyOnce sync.Once
	once sync.Once
	SuiMainNetEndpoint = "https://fullnode.mainnet.sui.io"
	servMap     *sync.Map
)

type SuiService struct {
	SuiApi *client.Client
	BlockApi *sui.ISuiAPI
}

func InitSuiService(endpointList ...string) *SuiService{
	if servMap == nil{
		servMap = &sync.Map{}
	}
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
		servMap.Store(SuiMainNetEndpoint, instance)
	})

	for _,endpoint := range endpointList{
		_, ok := servMap.Load(endpoint)
		if !ok{
			c, err := client.Dial(endpoint)
			if err != nil {
				errorMsg := fmt.Sprintf("connect sui main net error:%v", err)
				fmt.Printf("\n==errorMsg:%v==\n",errorMsg)
				continue
			}
			blockSuiApi := sui.NewSuiClient(endpoint)
			instance = &SuiService{
				SuiApi: c,
				BlockApi: &blockSuiApi,
			}
			servMap.Store(endpoint, instance)
		}
	}


	instanceValue, ok := utils.GetRandomValueFromSyncMap(servMap)
	fmt.Printf("instanceValue:%v, ok:%v",instanceValue, ok)
	if ok {
		if suiService, typeOk := instanceValue.(*SuiService); typeOk {
			return suiService
		}
	}
	return nil
}

func InitSuiServiceByOnlyEndpoint(endpoint string) *SuiService{
	onlyOnce.Do(func() {
		c, err := client.Dial(endpoint)
		if err != nil {
			errorMsg := fmt.Sprintf("connect sui main net error:%v", err)
			panic(errorMsg)
		}
		blockSuiApi := sui.NewSuiClient(endpoint)
		onlyEndpointInstance = &SuiService{
			SuiApi: c,
			BlockApi: &blockSuiApi,
		}
	})

	return onlyEndpointInstance
}