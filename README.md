# nemo-go-sdk
Initialize wallet and access endpoints
```bigquery
sender := account2.NewAccountPrivateKey("")
s := service.InitSuiService()
```

get support nemo pool
```bigquery
list := models.InitConfig()
nemoConfig := list[0]
```

MintPy
```bigquery
s.MintPy(0.001, sender, &nemoConfig)
```

RedeemPy
```bigquery
s.RedeemPy(0.001, sender, &nemoConfig)
```

source Coin swap to yt Coin 
```bigquery
s.SwapToPy(0.000000001, 0.005, nemoConfig.CoinType, constant.YTTYPE, sender, &nemoConfig)
```

yt Coin swap to source Coin 
```bigquery
s.SwapByPy(0.000000001, 0.005, constant.YTTYPE, nemoConfig.CoinType, sender, &nemoConfig)	
```

add Liquidity
```bigquery
s.AddLiquidity(0.00001, 0.005, sender, nemoConfig.UnderlyingCoinType, &nemoConfig)
```

redeem Liquidity
```bigquery
s.RedeemLiquidity(0.00001,  0.005, sender, nemoConfig.UnderlyingCoinType, &nemoConfig)
```

get pool apy
```js
s.QueryPoolApy(&nemoConfig)
```