package main

import (
	"TelecamFeeAnalyzer"
	"fmt"
)

func main() {
	ret := TelecamFeeAnalyzer.ParseLandmark()
	for _, r := range ret {
		if (r.Limit == 0 || r.Limit > 40) && r.Plan == "NP" {
			fmt.Println(r)
		}
	}
}
