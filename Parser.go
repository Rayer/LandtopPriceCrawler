package TelecamFeeAnalyzer

import (
	_ "embed"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	reader := strings.NewReader(target)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}
	var providers []string
	doc.Find(".brand-tags span").Each(func(i int, selection *goquery.Selection) {
		providers = append(providers, selection.Text())
	})

	//order of providers will be used in next parse
	//id will be "tab-index-{order}"

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

	for i, provider := range providers {
		providerNode := doc.Find(fmt.Sprintf("#tab-index-%d", i))
		for planTypeIndex, _ := range planTypes {
			providerNode.Find("div .tab-content.no-border.no-shadow").ChildrenFiltered("div").Each(func(i int, priceTable *goquery.Selection) {
				signalType := signalTypes[0]
				priceTable.Find("tr").Each(func(i int, priceRow *goquery.Selection) {
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
						Provider:      provider,
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

		}

	}

	for _, r := range ret {
		if r.Limit == 0 || r.Limit > 36 {
			fmt.Println(r)
		}
	}
	return
}
