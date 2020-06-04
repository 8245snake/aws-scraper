# aws-scraper

## 概要
docomoシェアサイクルのページをAWSでスクレイピングをするツール。  
Lambdaにデプロイしてcronなどで定期実行させる運用を想定する。

## 準備
### ビルド
Go言語で書かれており、以下のコマンドで依存ライブラリをインストールする必要がある。  
`go get -u github.com/aws/aws-lambda-go/lambda`
### デプロイ
2020/06/04現在、LambdaでGo言語をデプロイするためにはローカルでビルドしてzipをアップロードする方法しかない。  
build-lambda-zipでzip化することが推奨されているため以下のコマンドでインストールする。  
`go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip`
### 環境変数
設定をLambdaの環境変数に記載する。

| キー | 値 | 意味 |
|----|----|----|
|TZ |Asia/Tokyo |タイムゾーン |
|LoginID |*your login ID* |ログインID |
|Password |*your password* |パスワード|
|SendAddress |*https://hogehoge.com/callback* |スクレイピング結果を送信する独自のエンドポイント |
|ApiCert |*secret token* |スクレイピング結果を送信するときにヘッダーに付加する秘密の文字列 |
|AreaIdString |1,2,3,5,6,4,10,12,7,8 |スクレイピング対象とするエリアのID |

### DB登録
環境変数`SendAddress`に設定したエンドポイントにPOSTされるのでそれを適切に処理して保存するサーバが必要。  
スクレイピングしたデータは以下のJSONで来る。  
```
{
  "spotinfo": [
    {
      "area": "D1",
      "spot": "10",
      "time": "2020/02/16 15:43:20",
      "count": "29"
    },
    {
      "area": "D1",
      "spot": "10",
      "time": "2020/02/16 15:43:20",
      "count": "10"
    }
  ]
}
```

## 使い方
### 台数スクレイピング開始
上記の準備を整えたうえでLambdaにデプロイし、トリガーで以下の入力を投げる。  
 `{ "type" : "get_spotinfo"}`
 
### マスタ更新
スポットのコード（areaとspot）をと名前を紐付けるマスタ（spotmaster）の登録は以下の入力を投げる。  
 `{ "type" : "get_master"}`
 
 ### リカバリ
何らかの事情でスクレイピング結果の送信に失敗したとき（DBサーバが落ちてるなど）、/tmp フォルダにJSONファイルとして溜めておき、あとから送信するという仕組みがある。  
以下の入力を投げることで20件ずつ送信することができる。  
 `{ "type" : "recovery"}`
