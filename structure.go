package main

import "time"

//////////////////////////////////////////////////////////////////////////////////////
// 構造体
//////////////////////////////////////////////////////////////////////////////////////

//SpotInfo スクレイピング結果を格納する構造体
type SpotInfo struct {
	Time                              time.Time
	Area, Spot, Count, Name, Lat, Lon string
}

//JSpotinfo JSONマーシャリング構造体
type JSpotinfo struct {
	Spotinfo []InnerSpotinfo `json:"spotinfo"`
}

//InnerSpotinfo 台数情報
type InnerSpotinfo struct {
	Time  string `json:"time"`
	Area  string `json:"area"`
	Spot  string `json:"spot"`
	Count string `json:"count"`
}

//JSpotmaster JSONマーシャリング構造体
type JSpotmaster struct {
	Spotmaster []InnerSpotmaster `json:"spotmaster"`
}

//InnerSpotmaster スポット情報
type InnerSpotmaster struct {
	Area string `json:"area"`
	Spot string `json:"spot"`
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

//////////////////////////////////////////////////////////////////////////////////////
// レシーバ
//////////////////////////////////////////////////////////////////////////////////////

//Add SpotInfo構造体をJSON用にパースして加える
func (s *JSpotinfo) Add(time time.Time, area string, spot string, count string) {
	s.Spotinfo = append(s.Spotinfo, InnerSpotinfo{Time: time.Format(TimeLayout), Area: area, Spot: spot, Count: count})
}

//Size SpotInfo構造体のサイズを返す
func (s *JSpotinfo) Size() int {
	return len(s.Spotinfo)
}

//Add SpotInfo構造体をJSON用にパースして加える
func (s *JSpotmaster) Add(area string, spot string, name string, lat string, lon string) {
	s.Spotmaster = append(s.Spotmaster, InnerSpotmaster{Area: area, Spot: spot, Name: name, Lat: lat, Lon: lon})
}

//Size SpotInfo構造体のサイズを返す
func (s *JSpotmaster) Size() int {
	return len(s.Spotmaster)
}
