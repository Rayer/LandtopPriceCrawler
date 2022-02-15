package TelecamFeeAnalyzer

import (
	_ "embed"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//go:embed resources/test/LandtopResponse.html
var target string

const (
	np      = "tab-np_plan"
	renew   = "tab-renew_plan"
	newUser = "tab-new_plan"
)

func ParseLandmark() (ret []FeeEntry) {
	client := http.Client{}
	providerQueryArg := []string{"cht", "fet", "twn", "tstar", "apgt"}
	providerTwNames := []string{"中華電信", "遠傳", "台灣大哥大", "台灣之星", "亞太電信"}

	planTypes := []string{
		np, renew, newUser,
	}

	planText := []string{
		"NP", "續約", "新申辦",
	}

	signalTypes := []string{
		"4G", "5G",
	}

	feeMonthRegex := regexp.MustCompile(`(\d+)+\((\d+)\)`)
	limitRegex := regexp.MustCompile(`(\d+)+`)

	for pi, providerArg := range providerQueryArg {
		providerTwName := providerTwNames[pi]
		resp, _ := client.Get("https://www.landtop.com.tw/discount?provider=" + providerArg)
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}
		providerNode := doc.Find(fmt.Sprintf("#tab-index-0"))
		for planTypeIndex, planType := range planTypes {
			providerNode.Find("div#" + planType).Each(func(planOrder int, planTable *goquery.Selection) {
				planTable.Find("table").Each(func(signalTypeOrder int, priceTable *goquery.Selection) {
					priceTable.Find("tr").Each(func(i int, priceRow *goquery.Selection) {
						signalType := signalTypes[signalTypeOrder]
						priceDetail := priceRow.Find("td")
						if len(priceDetail.Nodes) == 0 {
							return
						}
						//Handle entries
						index := 0
						feeAndMonth := feeMonthRegex.FindStringSubmatch(priceDetail.Eq(index).Text())
						index++
						fee, _ := strconv.Atoi(feeAndMonth[1])
						months, _ := strconv.Atoi(feeAndMonth[2])
						rebate, _ := strconv.Atoi(strings.Join(strings.Fields(priceDetail.Eq(index).Text()), ""))
						index++
						prepaid, _ := strconv.Atoi(strings.Join(strings.Fields(priceDetail.Eq(index).Text()), ""))
						index++

						//處理傳輸量
						limit := 0
						limitText := priceDetail.Eq(index).Text()
						index++
						limitParse := limitRegex.FindStringSubmatch(limitText)
						if len(limitParse) != 0 {
							limit, _ = strconv.Atoi(limitParse[1])
						}

						//處理hotspot傳輸量(5g only)
						hotspotLimit := 0
						if signalType == signalTypes[1] {
							hotspotLimitText := priceDetail.Eq(index).Text()
							index++
							hotspotLimitParse := limitRegex.FindStringSubmatch(hotspotLimitText)
							if len(hotspotLimitParse) != 0 {
								hotspotLimit, _ = strconv.Atoi(hotspotLimitParse[1])
							}
						}

						description := priceDetail.Eq(index).Text()

						ret = append(ret, FeeEntry{
							Provider:      providerTwName,
							SignalType:    signalType,
							Limit:         limit,
							HotspotLimit:  hotspotLimit,
							ContractMonth: months,
							PrePaid:       prepaid,
							MonthlyFee:    fee,
							Rebate:        rebate,
							Plan:          planText[planTypeIndex],
							Description:   description,
						})
					})
				})
			})
		}
		_ = resp.Body.Close()
	}

	for _, r := range ret {
		if (r.Limit == 0 || r.Limit > 40) && r.RealMonthlyCost() < 700 && r.Plan == "NP" {
			fmt.Println(r)
		}
	}
	return
}
