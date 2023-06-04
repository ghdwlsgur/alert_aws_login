package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambdaEvent(ctx context.Context) (string, error) {

	// Telegram
	//==================================================================================
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	content := "AWS 루트 계정의 콘솔 로그인이 감지되었습니다."

	postData := map[string]interface{}{
		"chat_id":                  chatID,
		"disable_web_page_preview": true,
		"text":                     content,
	}

	jsonData, err := json.Marshal(postData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	str := string(b)
	fmt.Println(str)

	// NHN Cloud (Message - SMS)
	//==================================================================================
	appKey := os.Getenv("SMS_APP_KEY")
	senderNo := os.Getenv("SMS_SEND_NO")
	smsUrl := "https://api-sms.cloud.toast.com"
	smsUri := fmt.Sprintf("/sms/v3.0/appKeys/%s/sender/sms", appKey)
	SecretKey := os.Getenv("SMS_SECRET_KEY")

	postData = map[string]interface{}{
		"body":   "AWS 루트 계정의 콘솔 로그인이 감지되었습니다.",
		"sendNo": senderNo,
		"recipientList": []map[string]interface{}{
			{
				"recipientNo": senderNo,
				"countryCode": "82",
			},
		},
	}

	jsonData, err = json.Marshal(postData)
	if err != nil {
		panic(err)
	}

	smsTarget := fmt.Sprintf("%s%s", smsUrl, smsUri)

	req, err = http.NewRequest("POST", smsTarget, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Secret-Key", SecretKey)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ = ioutil.ReadAll(resp.Body)
	str = string(b)
	fmt.Println(str)

	return "send message", nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
