package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PayOutRequestDemo() {
	fmt.Println("=====> step2 : Create Access Token. You need set your timestamp|clientKey|privateKey")

	accessToken := AccessToken()
	//get time
	timestamp := GetTimeStamp()
	//get merchantId from merchant platform
	merchantId := merchantIdSandBox
	baseUrl := baseUrlSanbox
	merchantSercet := merchantSecretSandBox
	//build string to sign
	stringToSign := merchantId + "|" + timestamp
	fmt.Println(stringToSign)

	money := Money{Currency: "IDR", Amount: 10000}
	merchant := Merchant{MerchantId: merchantId}

	payoutRequest := PayOutRequest{OrderNo: merchantId + "ddfd",
		Purpose:     "for test demo",
		Merchant:    merchant,
		Money:       money,
		CashAccount: "1231232132", Area: 10, PaymentMethod: "BCA"}
	requestJson, _ := json.Marshal(payoutRequest)

	lowerString := LowerHexSha256Body(string(requestJson))

	signString := "POST:" + "/v1.0/disbursement/cash-out" + ":" + accessToken + ":" + lowerString + ":" + timestamp
	//signature
	signatureString, _ := hmacSHA512(signString, merchantSercet)

	//postJson
	postPayOutRequestDemo(timestamp, merchantId, signatureString, baseUrl, accessToken, payoutRequest)
}

func postPayOutRequestDemo(timestamp string, merchantId string, signatureString string, baseUrl string, accessToken string, payoutRequest PayOutRequest) string {
	// Create the JSON payload
	requestJson, _ := json.Marshal(payoutRequest)

	// Send the POST request
	url := baseUrl + "/v1.0/disbursement/cash-out"
	fmt.Println("request path:" + url)
	fmt.Println("request request param:" + string(requestJson))
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJson))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// Add custom headers
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-TIMESTAMP", timestamp)
	request.Header.Add("X-PARTNER-ID", merchantId)
	request.Header.Add("X-SIGNATURE", signatureString)
	request.Header.Add("X-EXTERNAL-ID", timestamp)
	request.Header.Add("Authorization", "Bearer  "+accessToken)

	// Send the request
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""

	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	// 读取响应体
	if err != nil {
		fmt.Println("Error:", err)
		return ""

	}

	// 打印响应状态码
	fmt.Println("Status Code:", response.StatusCode)

	bodyString := string(body)
	// 打印响应体
	fmt.Println("Response Body:", bodyString)

	var accessTokenBean AccessTokenBean

	err = json.Unmarshal([]byte(bodyString), &accessTokenBean)
	if err != nil {
		fmt.Println("Error:", err)
		return ""

	}
	fmt.Println("Email:", accessTokenBean.AccessToken)
	return accessTokenBean.AccessToken

}
