package unitpay

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

var testRequestUrl = "https://test.ru/unitpayhandler?method=check&params[account]=test&params[date]=2020-11-08 17:44:21&params[ip]=100.200.10.100&params[isPreauth]=0&params[operator]=mts&params[orderCurrency]=RUB&params[orderSum]=10.00&params[payerCurrency]=RUB&params[payerSum]=11.43&params[paymentType]=mc&params[phone]=0&params[profit]=9.3&params[projectId]=11111&params[signature]=b89e3b5e6d38451fa2aabc14598eb1f1dfd41379aa0a675338d3534a76450169&params[sum]=10&params[test]=1&params[unitpayId]=1"


func TestInit(t *testing.T) {
	Init("test", "111-111")
}

func TestGetSignature(t *testing.T)  {
	params := NewParams(100, "1", "Buy 100 for 1 user")

	signature := params.GetSignature()
	if signature != "e79ea8c952f8abbecd71d19714b46f9de4a84906e185d94a98902c04fedae5c4" {
		t.Error("Expect e79ea8c952f8abbecd71d19714b46f9de4a84906e185d94a98902c04fedae5c4, got ", signature)
	}

	r, _ := http.NewRequest(http.MethodGet, testRequestUrl, nil)
	hr := GetParamsFromRequest(r)

	signature1 := hr.getSignature()
	if signature1 != "b89e3b5e6d38451fa2aabc14598eb1f1dfd41379aa0a675338d3534a76450169"  {
		t.Error("Expect b89e3b5e6d38451fa2aabc14598eb1f1dfd41379aa0a675338d3534a76450169, got ", hr.Signature)
	}
}

func TestForm(t *testing.T)  {
	params := NewParams(100, "1", "Buy 100 for 1 user")

	url, _ := params.Form()
	if url != "https://unitpay.money/pay/111-111/card?account=1&currency=RUB&desc=Buy+100+for+1+user&locale=ru&signature=e79ea8c952f8abbecd71d19714b46f9de4a84906e185d94a98902c04fedae5c4&sum=100" {
		t.Error("Expect https://unitpay.money/pay/111-111/card?account=1&currency=RUB&desc=Buy+100+for+1+user&locale=ru&signature=e79ea8c952f8abbecd71d19714b46f9de4a84906e185d94a98902c04fedae5c4&sum=100, got ", url)
	}
}

func TestFormGetRequest(t *testing.T) {
	SetDefaultPayBaseUrl("https://unitpay.money/pay/demo")
	params := NewParams(100, "1", "Buy 100 for 1 user")

	url, _ := params.Form()

	resp, err := http.Get(url)
	var s string
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	for true {

		bs := make([]byte, 1014)
		n, err := resp.Body.Read(bs)
		s = s + string(bs[:n])

		if n == 0 || err != nil {
			break
		}
	}
	if strings.Contains(s, "Ошибка 400. Не переданы требуемые параметры запроса (desc, account, sum)") || strings.Contains(s, "Ошибка 404 &quot;Страница не найдена&quot;") {
		t.Error("Unitpay response with error")
	}
}

func TestHandleRequest(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, testRequestUrl, nil)
	_, err := HandleRequest(r)

	if err != nil {
		t.Error(err)
	}
}

func TestGetParamsFromRequest(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, testRequestUrl, nil)
	hr := GetParamsFromRequest(r)

	if hr.Account != "test" {
		t.Error("Expect test, got ", hr.Account)
	}
}