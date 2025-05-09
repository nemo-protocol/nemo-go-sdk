package service

import (
	"fmt"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/nemo-protocol/nemo-go-sdk/utils"
	"math/rand"
	"reflect"
	"sync"
)

var (
	instance    *SuiService
	onlyEndpointInstance *SuiService
	onlyOnce sync.Once
	once sync.Once
	SuiMainNetEndpoint = constant.SuiMainnetEndpoint
	servMap     *sync.Map
	onlyServMap *sync.Map
)

type SuiService struct {
	SuiApi *client.Client
	BlockApi *sui.ISuiAPI
}

func InitSuiService(params ...map[string]interface{}) *SuiService{
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

	priority := 1
	if len(params) > 0 {
		endpointList,ok := params[0]["endpointList"]
		if ok && reflect.TypeOf(endpointList).Kind() == reflect.Slice{
			for _, endpoint := range endpointList.([]string) {
				_, ok = servMap.Load(endpoint)
				if !ok {
					inst := createInstance(endpoint)
					if inst == nil{
						continue
					}
					servMap.Store(endpoint, inst)
				}
			}
		}

		priorityElement,ok := params[0]["priority"]
		if ok && reflect.TypeOf(priority).Kind() == reflect.Int{
			priority = priorityElement.(int)
		}
	}

	mainNetInstance, _ := servMap.Load(SuiMainNetEndpoint)

	switch priority {
	case 3:
		var candidates []*SuiService
		servMap.Range(func(key, value interface{}) bool {
			if key != SuiMainNetEndpoint {
				if inst, ok := value.(*SuiService); ok {
					candidates = append(candidates, inst)
				}
			}
			return true
		})
		if len(candidates) > 0 {
			return candidates[rand.Intn(len(candidates))]
		}
		return mainNetInstance.(*SuiService)

	case 2:
		if rand.Intn(2) == 0 {
			return mainNetInstance.(*SuiService)
		}
		instanceValue, ok := utils.GetRandomValueFromSyncMap(servMap)
		if ok {
			if suiService, typeOk := instanceValue.(*SuiService); typeOk {
				return suiService
			}
		}
		return mainNetInstance.(*SuiService)

	case 1:
		fallthrough
	default:
		return mainNetInstance.(*SuiService)
	}
}

func createInstance(endpoint string) *SuiService{
	c, err := client.Dial(endpoint)
	if err != nil {
		errorMsg := fmt.Sprintf("connect sui main net error:%v", err)
		fmt.Printf("\n==errorMsg:%v==\n", errorMsg)
		return nil
	}
	blockSuiApi := sui.NewSuiClient(endpoint)
	instance = &SuiService{
		SuiApi:   c,
		BlockApi: &blockSuiApi,
	}
	return instance
}

func InitSuiServiceByOnlyEndpoint(endpoint ...string) *SuiService{
	if onlyServMap == nil{
		onlyServMap = &sync.Map{}
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
	})
	_, ok := onlyServMap.Load(SuiMainNetEndpoint)
	if !ok{
		onlyServMap.Store(SuiMainNetEndpoint, instance)
	}

	for _,v := range endpoint{
		c, err := client.Dial(v)
		fmt.Printf("c:%+v",c)
		if err != nil {
			errorMsg := fmt.Sprintf("connect sui main net error:%v", err)
			panic(errorMsg)
		}
		blockSuiApi := sui.NewSuiClient(v)
		fmt.Printf("blockSuiApi:%+v",blockSuiApi)
		onlyEndpointInstance = &SuiService{
			SuiApi: c,
			BlockApi: &blockSuiApi,
		}
		onlyServMap.Store(v, onlyEndpointInstance)
	}

	instanceValue, ok := utils.GetRandomValueFromSyncMap(onlyServMap)
	fmt.Printf("instanceValue:%v, ok:%v",instanceValue, ok)
	if ok {
		if suiService, typeOk := instanceValue.(*SuiService); typeOk {
			return suiService
		}
	}
	return nil
}