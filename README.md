# LandtopPriceCrawler

這是爬取地標網通的價格表 (範例: https://www.landtop.com.tw/discount?provider=cht)來獲得電信商的價格跟折扣。

## 使用方式

目前我沒打算花時間做出完整命令列，所以請自己改code去設定條件。

修改`cmd/main.go`

```go
ret := TelecamFeeAnalyzer.ParseLandmark()
for _, r := range ret {
    if (r.Limit == 0 || r.Limit > 40) && r.Plan == "NP" {
    fmt.Println(r)
}

```

以這範例來講就是擷取所有

- 吃掉飽或者40G以上的方案，且
- NP專案

Plan以及4G 5G字串可以參考

```go
planText := []string{
    "NP", "續約", "新申辦",
}

signalTypes := []string{
    "4G", "5G",
}
```

請直接執行這個main即可