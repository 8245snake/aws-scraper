package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var invokeCount = 0

//MyEvent イベント振り分け
type MyEvent struct {
	Type string `json:"type"`
}

const (
	//TypeTest テスト用
	TypeTest = "test"
	//TypeStartScraping スクレイピング開始
	TypeStartScraping = "get_spotinfo"
	//TypeStartMaster マスタ取得開始
	TypeStartMaster = "get_master"
	//TypeStartRecovery レカバリ開始
	TypeStartRecovery = "recovery"
)

//LambdaHandler ハンドラ
func LambdaHandler(ctx context.Context, event MyEvent) (result string, err error) {
	switch event.Type {
	case TypeTest:
	case TypeStartScraping:
		err = RegAllSpotInfo()
	case TypeStartMaster:
		err = RegAllSpotMaster()
	case TypeStartRecovery:
		Recover()

	}
	result = fmt.Sprintf("time=%s type=%s, id=%s pw=%s addr=%s cert=%s area=%s session=%s",
		time.Now().Format(TimeLayout),
		event.Type,
		UserID,
		Password,
		SendAddress,
		ApiCert,
		AreaIdString,
		SessionID,
	)
	return
}

//InitClient クライアント初期化
func InitClient() {
	//SSL証明書を無視したクライアント作成
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
	}
}

func init() {
	InitClient()
	UserID = os.Getenv("LoginID")
	Password = os.Getenv("Password")
	SendAddress = os.Getenv("SendAddress")
	AreaIdString = os.Getenv("AreaIdString")
	ApiCert = os.Getenv("ApiCert")
	if sess, err := GetSessionID(); err == nil {
		SessionID = sess
	} else {
		panic(err)
	}
}

func main() {
	lambda.Start(LambdaHandler)
}
