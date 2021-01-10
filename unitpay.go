package unitpay

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var SecretKey string
var PublicKey string

/*
Init unitpay with secret key, public key and default values.
*/
func Init(secret string, public string)  {
	SecretKey = secret
	PublicKey = public
}

// Allowed values for Optional Param currency
var currencyValues = [34]string{"RUB","EUR","USD","AUD", "AZN", "AMD", "BYN", "BGN", "BRL", "HUF", "KRW", "HKD", "DKK", "INR", "KZT", "CAD", "KGS", "CNY", "MDL", "TMT", "NOK", "PLN", "RON", "SGD", "TJS", "TRY", "UZS", "UAH", "GBP", "CZK", "SEK", "CHF", "ZAR", "JPY"}
// Allowed values for Optional Param locale
var localeValues = [2]string{"ru", "en"}
// Allowed values for Optional Param paymentMethods
var paymentMethodsValues = []string{"mc","card","webmoney","webmoneyWmr","yandex","qiwi","paypal","alfaClick","applepay","samsungpay"}
// Allowed methods
var allowedMethods = [3]string{"check","pay","error"}
// TrustedUnitpayIps list
var supportedUnitpayIps = [5]string{
	"31.186.100.49",
	"178.132.203.105",
	"52.29.152.23",
	"52.19.56.234",
	"127.0.0.1",
}
// Check if method is allowed
func checkMethod(method string) bool  {
	for _, v := range allowedMethods {
		if method == v {
			return true
		}
	}
	return false
}

// Check ip of request
func CheckIp(ip string) bool  {
	for _, v := range supportedUnitpayIps {
		if ip == v {
			return true
		}
	}
	return false
}


// Check if unitpay accept Locale/Currency values
func checkCurrencyValue(value string) bool {
	for _, v := range currencyValues {
		if v == value {
			return true
		}
	}
	return false
}

// Check if unitpay accept Locale/Currency values
func checkLocaleValue(value string) bool {
	for _, v := range localeValues {
		if v == value {
			return true
		}
	}
	return false
}

/*
Form
*/
// Type of default values for form
type DefaultValues struct {
	defaultCurrency string
	defaultLocale string
	defaultBackUrl string
	defaultPayBaseUrl string
	defaultPaymentMethod string
	defaultHideMenu string
}

// Preset default values for form
var defaultValues = DefaultValues{
	defaultCurrency:   "RUB",
	defaultLocale:     "ru",
	defaultBackUrl:    "",
	defaultPayBaseUrl: "https://unitpay.money/pay/",
	defaultPaymentMethod: "card",
	defaultHideMenu: "true",
}

// Functions to change some of default values for form
func SetDefaultCurrency(currency string) error {
	if !checkCurrencyValue(currency) {
		return errors.New("this currency value is not allowed, check allowed values in docs")
	}
	defaultValues.defaultCurrency = currency
	return nil
}
// Functions to change some of default values for form
func SetDefaultLocale(locale string) error  {
	if !checkLocaleValue(locale) {
		return errors.New("this locale value is not allowed, allowed values: ru, en")
	}
	defaultValues.defaultLocale = locale
	return nil
}
// Functions to change some of default values for form
func SetDefaultBackUrl(backurl string)  {
	defaultValues.defaultBackUrl = backurl
}
// Functions to change some of default values for form
func SetDefaultPayBaseUrl(defaultpaybaseurl string)  {
	defaultValues.defaultCurrency = defaultpaybaseurl
}
// Functions to change some of default values for form
func SetDefaultPaymentMethod(paymentmethod string) error  {
	if !checkPaymentMethodsValue(paymentmethod) {
		return errors.New("this payment method is not allowed, check allowed values in docs")
	}
	defaultValues.defaultPaymentMethod = paymentmethod
	return nil
}
// Functions to change some of default values for form
func SetDefaultHideMenu(defaulthidemenu bool)  {
	defaultValues.defaultHideMenu = strconv.FormatBool(defaulthidemenu)
}



// Check if unitpay accept PaymentMethods values
func checkPaymentMethodsValue(value string) bool {
	for _, v := range paymentMethodsValues {
		if v == value {
			return true
		}
	}
	return false
}

/*
Simple payment form Required Params
Used in func Form
*/
type RequiredParams struct {
	Sum int `json:"sum"`
	Account string `json:"account"`
	Desc string `json:"desc"`
	Signature string `json:"signature"`
}

/*
Other simple payment form Params (non-required).
Used in func Form
*/
type OptionalParams struct {
	Currency string `json:"currency"`
	Locale string `json:"locale"`
	BackUrl string `json:"backUrl"`
	PaymentMethod string `json:"payment_method"`
	HideMenu string `json:"hide_menu"`
}

// Params options summary
type Params struct {
	RequiredParams
	OptionalParams
}


// Group of functions to set Form params values
func (params *Params) SetSum(sum int) error {
	if sum <= 0 {
		return errors.New("sum must be greater than zero")
	}
	params.Sum = sum
	return nil
}
// Group of functions to set Form params values
func (params *Params) SetDesc(desc string)  {
	params.Desc = desc
}
// Group of functions to set Form params values
func (params *Params) SetAccount(account string) error  {
	if account == "" {
		return errors.New("account must be not null")
	}
	params.Account = account
	return nil
}
// Group of functions to set Form params values
func (params *Params) SetCurrency(currency string) error  {
	params.Desc = currency
	return nil
}
// Group of functions to set Form params values
func (params *Params) SetLocale(locale string) error  {
	params.Desc = locale
	return nil
}
// Group of functions to set Form params values
func (params *Params) SetBackUrl(backurl string) error  {
	params.Desc = backurl
	return nil
}

// Group of functions to set Form params values
func (params *Params) SetPaymentMethod(paymentmethod string) error  {
	if !checkPaymentMethodsValue(paymentmethod) {
		return errors.New("this payment method is not allowed, check allowed values in docs")
	}
	params.PaymentMethod = paymentmethod
	return nil
}

// Group of functions to set Form params values
func (params *Params) SetHideMenu(hidemenu bool) error  {
	params.HideMenu = strconv.FormatBool(hidemenu)
	return nil
}

// Return new params object with given required params and default optional params
func NewParams(sum int, account string, desc string) *Params {
	return &Params{
		RequiredParams: RequiredParams{
			Sum: sum,
			Account: account,
			Desc: desc,
		},
		OptionalParams: OptionalParams{
			Currency: defaultValues.defaultCurrency,
			Locale:   defaultValues.defaultLocale,
			BackUrl:  defaultValues.defaultBackUrl,
			PaymentMethod: defaultValues.defaultPaymentMethod,
			HideMenu: defaultValues.defaultHideMenu,
		},
	}
}
// return new params object with empty required params and default optional params
func NewEmptyParams() *Params {
	return &Params{
		RequiredParams: RequiredParams{},
		OptionalParams: OptionalParams{
			Currency: defaultValues.defaultCurrency,
			Locale:   defaultValues.defaultLocale,
			BackUrl:  defaultValues.defaultBackUrl,
			PaymentMethod: defaultValues.defaultPaymentMethod,
			HideMenu: defaultValues.defaultHideMenu,
		},
	}
}

// func to get signature of given params
func (params *Params) GetSignature() string {

	var p []string

	p = append(p, params.Account)
	if params.Currency != "" {
		p = append(p, params.Currency)
	}
	p = append(p, params.Desc, strconv.Itoa(params.Sum))

	// append secret to params
	p = append(p, SecretKey)

	// convert params to string and add separators {up}
	paramsString := strings.Join(p, "{up}")

	// hash params and get string of hash sum
	hash := sha256.New()
	hash.Write([]byte(paramsString))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	return mdStr
}

/*
Function to init simple payment form.
Returns url of payment form.
*/
func (params *Params) Form() (URL string, err error)  {
	if params.Account == "" || params.Desc == "" {
		return "", errors.New("required request params were not passed")
	}

	if params.Sum <= 0 {
		return "", errors.New("sum must be greater than zero")
	}
	signature := params.GetSignature()

	form := defaultValues.defaultPayBaseUrl

	var req *http.Request
	if params.PaymentMethod != "" {
		req, _ = http.NewRequest("GET", form +PublicKey+ "/" + params.PaymentMethod + "?", nil)
	} else {
		req, _ = http.NewRequest("GET", form +PublicKey+  "?", nil)
	}


	query := req.URL.Query()

	query.Add("sum", strconv.Itoa(params.Sum))
	query.Add("account", params.Account)
	query.Add("desc", params.Desc)

	if params.OptionalParams.BackUrl != "" {
		query.Add("backUrl", params.OptionalParams.BackUrl )
	}
	if params.OptionalParams.Locale != "" {
		query.Add("locale", params.OptionalParams.Locale )
	}
	if params.OptionalParams.Currency != "" {
		query.Add("currency", params.OptionalParams.Currency )
	}
	if params.OptionalParams.HideMenu != "" {
		query.Add("hideMenu", params.OptionalParams.HideMenu )
	}


	query.Add("signature", signature)

	req.URL.RawQuery = query.Encode()

	return req.URL.String(), nil
}

/*
Form end
*/


/*
Handler Request
*/
// All Request values from UnitpayDocs
type HandlerRequest struct {
	Method string `json:"method"`
	UnitpayId int `json:"unitpayId"`
	ProjectId int `json:"projectId"`
	Account string `json:"account"`
	PayerSum int `json:"payerSum"`
	PayerCurrency string `json:"payerCurrency"`
	Profit int `json:"profit"`
	Phone int `json:"phone"`
	PaymentType string `json:"paymentType"`
	OrderSum int `json:"orderSum"`
	OrderCurrency string `json:"orderCurrency"`
	Operator string `json:"operator"`
	Date string `json:"date"`
	ErrorMessage string `json:"errorMessage"`
	Test int `json:"test"`
	DS3 int `json:"3ds"`
	SubscriptionId int `json:"subscriptionId"`
	Signature string `json:"signature"`
	ParamsSlice []string `json:"params_slice"`
}

// Check method and currency for Handler Request
func (hr *HandlerRequest) checkHandleRequest() error  {
	if !checkMethod(hr.Method) {
		return errors.New("this method value is not allowed, check allowed values in docs")
	}

	if !checkCurrencyValue(hr.OrderCurrency) {
		return errors.New("this currency value is not allowed, check allowed values in docs")
	}

	return nil
}

// Local func to save all Params from Request, sort it, and append to slice
func getAllParamsFromRequest(r *http.Request)  []string {
	// Get all URL keys
	keys := r.URL.Query()

	// Delete signature from keys
	keys.Del("params[signature]")
	// Delete signature from keys with alternative key
	keys.Del("params[sign]")


	/* Array for all Keys of Request.
	It is necessary to sort Request values in correct sequence later.*/
	var paramsKeys []string

	for key := range keys {

		if strings.Contains(key, "params") {
			key = key[7:len(key)-1]

			key = strings.Replace(key, "[", "", -1)
			key = strings.Replace(key, "]", "", -1)
			key = strings.Replace(key, " ", "", -1)

			paramsKeys = append(paramsKeys, key)
		}

	}

	// Sorting keys in ASC
	sort.Strings(paramsKeys)

	// Array for params
	var params []string

	// Add params in ASC
	for _, key := range paramsKeys {
		key = fmt.Sprintf("params[%v]", key)
		key := keys[key]

		var key1 string
		for _, v := range key {
			key1 = v
		}

		params = append(params, key1)
	}

	params = append([]string{keys.Get("method")}, params...)

	return params
}

/* Function to get values from Request
You can access this function from outside of package to operate params manually
*/
func GetParamsFromRequest(r *http.Request) *HandlerRequest {
	var hr HandlerRequest
	hr.Method = r.URL.Query().Get("method")
	hr.UnitpayId, _ = strconv.Atoi(r.URL.Query().Get("params[unitpayId]"))
	hr.ProjectId, _ = strconv.Atoi(r.URL.Query().Get("params[projectId]"))
	hr.Account = r.URL.Query().Get("params[account]")
	hr.PayerSum, _ = strconv.Atoi(r.URL.Query().Get("params[payerSum]"))
	hr.PayerCurrency = r.URL.Query().Get("params[payerCurrency]")
	hr.Profit, _ = strconv.Atoi(r.URL.Query().Get("params[profit]"))
	hr.PaymentType = r.URL.Query().Get("params[paymentType]")
	hr.OrderCurrency = r.URL.Query().Get("params[orderCurrency]")
	hr.Operator = r.URL.Query().Get("params[operator]")
	hr.Date = r.URL.Query().Get("params[date]")
	hr.ErrorMessage = r.URL.Query().Get("params[errorMessage]")
	hr.Test, _ = strconv.Atoi(r.URL.Query().Get("params[test]"))
	hr.DS3, _ = strconv.Atoi(r.URL.Query().Get("params[3ds]"))
	hr.SubscriptionId, _ = strconv.Atoi(r.URL.Query().Get("params[subscriptionId]"))
	hr.Signature = r.URL.Query().Get("params[signature]")

	hr.ParamsSlice = getAllParamsFromRequest(r)

	return &hr
}

// Get signature of unitpay Handler Request
func (hr *HandlerRequest) getSignature() string {
	params := append(hr.ParamsSlice, SecretKey)
	paramsString := strings.Join(params, "{up}")

	// hash params and get string of hash sum
	hash := sha256.New()
	hash.Write([]byte(paramsString))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	return mdStr
}

/* Function to handle Unitpay request. Accept http.Request.
Returns Params of request and nil error if its ok.
Returns nil Params and error if there are errors.
Alternatively you can get params by yourself and handle it with HandleRequestWithParams function
*/
func HandleRequest(r *http.Request) (*HandlerRequest,error)  {
	hr := GetParamsFromRequest(r)

	err := hr.checkHandleRequest()

	if err != nil {
		return nil,err
	}

	signature := hr.getSignature()

	if signature != hr.Signature {
		return nil, errors.New("signature is not valid")
	}

	return hr, nil
}


/* Function to handle Unitpay request. Accept unitpay.HandleRequest.
Returns nil error if its ok.
Alternatively you can pass http.Request in HandleRequest function
*/
func HandleRequestWithParams(hr *HandlerRequest) error {
	err := hr.checkHandleRequest()

	if err != nil {
		return err
	}

	signature := hr.getSignature()

	if signature != hr.Signature {
		return errors.New("signature is not valid")
	}

	return nil
}