package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ant0ine/go-json-rest/rest"
)

//////////////////////////////////////////////////////////////////////////////////////
// 定数
//////////////////////////////////////////////////////////////////////////////////////

const (
	//TimeLayout 時刻フォーマット
	TimeLayout = "2006/01/02 15:04:05"
	//AllSpot 全スポット
	AllSpot = "1,2,3,5,6,4,10,12,7,8"
)

//////////////////////////////////////////////////////////////////////////////////////
// 変数
//////////////////////////////////////////////////////////////////////////////////////

//Httpでもらう設定値
var (
	UserID       string
	Password     string
	SendAddress  string
	AreaIdString string
	ApiCert      string
)

//SessionID セッション
var SessionID string

//client HTTPリクエストクライアント（使いまわした方がいいらしいのでグローバル化）
var client *http.Client

//////////////////////////////////////////////////////////////////////////////////////
// 関数
//////////////////////////////////////////////////////////////////////////////////////

//GetSessionID ログインしてセッションIDを取得する
func GetSessionID() (string, error) {
	//リクエストBody作成
	values := url.Values{}
	values.Set("EventNo", "21401")
	values.Add("GarblePrevention", "ＰＯＳＴデータ")
	values.Add("MemberID", UserID)
	values.Add("Password", Password)
	values.Add("MemAreaID", "1")

	req, err := http.NewRequest(
		"POST",
		"https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		fmt.Println("[Error]GetSessionID create NewRequest failed", err)
		return "", err
	}

	// リクエストHead作成
	ContentLength := strconv.FormatInt(req.ContentLength, 10)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8,pt;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", ContentLength)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "tcc.docomo-cycle.jp")
	req.Header.Set("Origin", "https://tcc.docomo-cycle.jp")
	req.Header.Set("Referer", "https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[Error]GetSessionID client.Do failed", err)
		return "", err
	}
	defer resp.Body.Close()

	doc, e := goquery.NewDocumentFromResponse(resp)
	if e != nil {
		fmt.Println("[Error]GetSessionID NewDocumentFromResponse failed", e)
		return "", e
	}

	sessionID, success := doc.Find("input[name='SessionID']").Attr("value")
	if !success {
		fmt.Println("[Error]GetSessionID Find SessionID failed")
		return "", fmt.Errorf("error")
	} else {
		fmt.Println("GetSessionID success ", sessionID)
		//成功したら待ち時間（1回目の検索に失敗するため）
		time.Sleep(3 * time.Second)
		return sessionID, nil
	}
}

//GetSpotInfoMain スクレイピングメイン関数
func GetSpotInfoMain(AreaID string, retry bool) ([]SpotInfo, error) {
	fmt.Printf("GetSpotInfoMain_start AreaID = %s \n", AreaID)
	var list []SpotInfo
	//リクエストBody作成
	values := url.Values{}
	values.Set("EventNo", "25706")
	values.Add("SessionID", SessionID)
	values.Add("UserID", "TYO")
	values.Add("MemberID", UserID)
	values.Add("GetInfoNum", "200")
	values.Add("GetInfoTopNum", "1")
	values.Add("MapType", "1")
	values.Add("MapCenterLat", "")
	values.Add("MapCenterLon", "")
	values.Add("MapZoom", "13")
	values.Add("EntServiceID", "TYO0001")
	values.Add("Location", "")
	values.Add("AreaID", AreaID)

	req, err := http.NewRequest(
		"POST",
		"https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		fmt.Println("[Error]GetSpotInfoMain create NewRequest failed", err)
		return nil, err
	}

	// リクエストHead作成
	ContentLength := strconv.FormatInt(req.ContentLength, 10)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8,pt;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", ContentLength)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "tcc.docomo-cycle.jp")
	req.Header.Set("Origin", "https://tcc.docomo-cycle.jp")
	req.Header.Set("Referer", "https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[Error]GetSpotInfoMain client.Do failed", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	doc, e := goquery.NewDocumentFromResponse(resp)
	if e != nil {
		fmt.Println("[Error]GetSpotInfoMain NewDocumentFromResponse failed", e)
		return nil, e
	}

	//エラーならログインし直して再チャレンジ
	if err := CheckErrorPage(doc); err != nil {
		if retry {
			SessionID, err = GetSessionID()
			if err != nil {
				return nil, err
			}
			//再帰呼び出し（次はリトライしない）
			return GetSpotInfoMain(AreaID, false)
		} else {
			//２回目は諦める
			return nil, err
		}
	}

	//スポットリスト解析
	doc.Find("form[name^=tab_]").Each(func(i int, s *goquery.Selection) {
		spotinfo := SpotInfo{Time: time.Now()}
		html, _ := s.Find("a").Html()
		err := ParseSpotInfoByText(html, &spotinfo)
		if err != nil {
			//メンテナンス中のスポットのエラーログは出力しない
			if strings.Index(err.Error(), "not cyclespot") < 0 {
				fmt.Println("[Error]GetSpotInfoMain ParseSpotInfoByText failed", err)
			}
			return
		}
		if val, exist := s.Find("input[name=ParkingLat]").Attr("value"); exist {
			spotinfo.Lat = val
		}
		if val, exist := s.Find("input[name=ParkingLon]").Attr("value"); exist {
			spotinfo.Lon = val
		}
		list = append(list, spotinfo)
	})

	fmt.Printf("GetSpotInfoMain_end AreaID = %s (%d件)\n", AreaID, len(list))
	return list, nil
}

//ParseSpotInfoByText テキスト解析
// "H1-43.東京イースト21<br/>H1-43.Tokyo East 21<br/>13台"の形式のテキストからarea,spot,name,countを取得する
func ParseSpotInfoByText(text string, s *SpotInfo) error {
	var codeAndName, cycleCount string
	if arr := strings.Split(text, "<br/>"); len(arr) == 3 {
		codeAndName = arr[0]
		cycleCount = arr[2]
	} else {
		return fmt.Errorf("[Error]ParseSpotInfoByText unexpected html : " + text)
	}

	// "H1-43"の部分
	indexDot := strings.Index(codeAndName, ".")
	if indexDot < 0 {
		return fmt.Errorf("[Error]ParseSpotInfoByText not cyclespot : " + text)
	}
	code := codeAndName[:indexDot]
	if arr := strings.Split(code, "-"); len(arr) == 2 {
		s.Area = arr[0]
		s.Spot = arr[1]
	} else {
		return fmt.Errorf("[Error]ParseSpotInfoByText unexpected code : " + text)
	}

	//名前
	s.Name = codeAndName[indexDot+1:]
	//台数
	if _, err := strconv.Atoi(cycleCount[:len(cycleCount)-3]); err == nil {
		s.Count = cycleCount[:len(cycleCount)-3]
	} else {
		return fmt.Errorf("[Error]ParseSpotInfoByText count not obtained : " + text)
	}

	//データサイズチェック
	if len(s.Area) > 3 || len(s.Spot) > 3 || len(s.Count) > 3 {
		fmt.Println("[Error]ParseSpotInfoByText data size obver : " + text)
	}

	return nil
}

//RegAllSpotInfo 全スポット登録関数
func RegAllSpotInfo() (err error) {
	//特に指定してない場合は全スポット
	if AreaIdString == "" {
		AreaIdString = AllSpot
	}
	fmt.Println("RegAllSpotInfo_Start AreaIdString =", AreaIdString)
	AreaIDs := strings.Split(AreaIdString, ",")
	for _, AreaID := range AreaIDs {
		if AreaID == "" {
			continue
		}
		//待ち時間いれる
		time.Sleep(5 * time.Second)
		//台数取得
		var list []SpotInfo
		list, err = GetSpotInfoMain(AreaID, true)
		if err != nil {
			fmt.Println("[Error]RegAllSpotInfo GetSpotInfoMain failed AreaID =", AreaID, err)
			continue
		}
		//負荷緩和のため100件ずつ送信
		max := 100
		jsondata := JSpotinfo{}
		for _, s := range list {
			jsondata.Add(s.Time, s.Area, s.Spot, s.Count)
			if jsondata.Size() >= max {
				SendSpotInfo(jsondata, false)
				jsondata = JSpotinfo{}
				time.Sleep(1 * time.Second)
			}
		}
		if jsondata.Size() >= 1 {
			SendSpotInfo(jsondata, false)
		}
	}
	fmt.Println("RegAllSpotInfo_End")
	return nil
}

//RegAllSpotMaster 全スポット登録関数（マスタメンテナンス）
func RegAllSpotMaster() (err error) {
	fmt.Println("RegAllSpotMaster_Start")
	//マスタメンテでは全スポットを対象とする
	AreaIDs := strings.Split(AllSpot, ",")
	for _, AreaID := range AreaIDs {
		//待ち時間いれる
		time.Sleep(5 * time.Second)
		//台数取得
		var list []SpotInfo
		list, err = GetSpotInfoMain(AreaID, true)
		if err != nil {
			fmt.Println("[Error]RegAllSpotMaster GetSpotInfoMain failed AreaID =", AreaID, err)
			continue
		}
		//負荷緩和のため100件ずつ送信
		max := 100
		jsondata := JSpotmaster{}
		for _, s := range list {
			jsondata.Add(s.Area, s.Spot, s.Name, s.Lat, s.Lon)
			if jsondata.Size() >= max {
				SendSpotMaster(jsondata)
				jsondata = JSpotmaster{}
				time.Sleep(1 * time.Second)
			}
		}
		if jsondata.Size() >= 1 {
			SendSpotMaster(jsondata)
		}
	}
	fmt.Println("RegAllSpotMaster_End")
	return nil
}

//CheckErrorPage エラーページかをチェックする
func CheckErrorPage(doc *goquery.Document) error {
	if title := doc.Find(".tittle_h1").Text(); strings.Index(title, "エラー") > -1 {
		fmt.Println(title)
		return fmt.Errorf(strings.TrimSpace(doc.Find(".main_inner_message").Text()))
	}
	return nil
}

//SendSpotInfo DBに送信する。JSONファイルからのリカバリの場合は失敗したらJSONを保存しないフラグ（第２引数）
func SendSpotInfo(jsonStruct JSpotinfo, fromRecovery bool) error {
	marshalized, _ := json.Marshal(jsonStruct)
	req, err := http.NewRequest(
		"POST",
		SendAddress,
		bytes.NewBuffer(marshalized),
	)
	if err != nil {
		fmt.Println("[Error]SendSpotInfo create NewRequest failed", err.Error())
		return err
	}

	// リクエストHead作成
	ContentLength := strconv.FormatInt(req.ContentLength, 10)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", ContentLength)
	req.Header.Set("cert", ApiCert)

	//送信
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[Error]SendSpotInfo client.Do failed", err.Error())
		if !fromRecovery {
			SaveJSON(jsonStruct)
		}
		return err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("[Error]SendSpotInfo StatusCode is not OK", resp.StatusCode, resp.Body)
		if !fromRecovery {
			SaveJSON(jsonStruct)
		}
		return fmt.Errorf("StatusCode is not OK : %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}

//SendSpotMaster マスタ情報をDBに送信する。
func SendSpotMaster(jsonStruct JSpotmaster) error {
	marshalized, _ := json.Marshal(jsonStruct)
	req, err := http.NewRequest(
		"POST",
		SendAddress,
		bytes.NewBuffer(marshalized),
	)
	if err != nil {
		fmt.Println("[Error]SendSpotMaster create NewRequest failed", err.Error())
		return err
	}

	// リクエストHead作成
	ContentLength := strconv.FormatInt(req.ContentLength, 10)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", ContentLength)
	req.Header.Set("cert", ApiCert)

	//送信
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[Error]SendSpotMaster client.Do failed", err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("[Error]SendSpotMaster StatusCode is not OK", resp.StatusCode, resp.Body)
		return fmt.Errorf("StatusCode is not OK : %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}

//PrepareScrayping スクレイピング準備（返り値がtrueの場合は実行しない）
func PrepareScrayping(w rest.ResponseWriter, r *rest.Request) (cancel bool) {
	//セッションIDを使いまわす
	var err error
	if SessionID == "" {
		SessionID, err = GetSessionID()
		if err != nil {
			fmt.Println("[Error]Start GetSessionID failed", err)
			w.WriteHeader(http.StatusBadRequest)
			w.WriteJson("login failed")
			return true
		}
	}
	return false
}

//SaveJSON JSONファイルに保存する
func SaveJSON(jsonStruct JSpotinfo) error {
	filePath := strconv.FormatInt(time.Now().Unix(), 10) + "_save.json"
	if runtime.GOOS != "windows" {
		filePath = "/tmp/" + filePath
	}
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer fp.Close()

	e := json.NewEncoder(fp)
	if err := e.Encode(jsonStruct); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//EnumTempFiles 一時ファイルを列挙する
func EnumTempFiles() (result []string) {
	pattern := "*_save.json"
	if runtime.GOOS != "windows" {
		pattern = "/tmp/" + pattern
	}
	if files, err := filepath.Glob(pattern); err != nil {
		fmt.Println("EnumTempFiles_Error :", err.Error())
	} else {
		result = files
	}
	return
}

//Recover JSONからリカバリー
func Recover() {
	//最大件数
	max := 20
	// if param := params.Get("max"); param != "" {
	// 	if val, err := strconv.Atoi(param); err == nil {
	// 		max = val
	// 	}
	// }

	//tmpファイルを列挙
	files := EnumTempFiles()
	if len(files) < 1 {
		fmt.Println("tmpにファイルがありません")
		return
	} else if max == 0 {
		msg := fmt.Sprintf("%d files found : %v \n", len(files), files)
		fmt.Printf(msg)
		return
	}
	for i, filename := range files {
		if i >= max {
			break
		}
		path := filename
		file, err := os.Open(path)
		if err != nil {
			msg := fmt.Sprintf("%s Open error : %v", path, err)
			fmt.Println(msg)
			return
		}
		defer file.Close()
		d := json.NewDecoder(file)
		d.DisallowUnknownFields() // エラーの場合 json: unknown field "JSONのフィールド名"
		var jsonstruct JSpotinfo
		if err := d.Decode(&jsonstruct); err != nil && err != io.EOF {
			msg := fmt.Sprintf("%s Decode error : %v", path, err)
			fmt.Println(msg)
			return
		}
		//DB登録処理
		if err := SendSpotInfo(jsonstruct, true); err != nil {
			//同じファイルで失敗し続けないようにしたいが何回かリトライのチャンスを与えたいのでMAX回数を引き上げる
			max++
			fmt.Printf("%s SendSpotInfo error : %v \n", path, err)
		} else {
			//成功したらファイル削除
			if err := os.Remove(path); err != nil {
				fmt.Printf("%s Remove error : %v \n", path, err)
				continue
			}
			fmt.Printf("%s Recover success \n", path)
		}
	}
}
