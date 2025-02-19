# nemo-go-sdk
Initialize wallet and access endpoints
```bigquery
sender := account2.NewAccountPrivateKey("")
s := service.InitSuiService()
```

MintPy
```bigquery
s.MintPy(0.001, sender, models.InitConfig())
```

RedeemPy
```bigquery
s.RedeemPy(0.001, sender, models.InitConfig())
```

source Coin swap to yt Coin 
```bigquery
s.SwapToPy(0.000000001, 0.005, models.InitConfig().CoinType, constant.YTTYPE, sender, models.InitConfig())
```

yt Coin swap to source Coin 
```bigquery
s.SwapByPy(0.000000001, 0.005, constant.YTTYPE, models.InitConfig().CoinType, sender, models.InitConfig())	
```
