package main

import (
    "log"
    "net/http"
    "io"
    "encoding/json"
)

type Player struct {
	Id     int `json:"id"` 
	Names  []string `json:"names"`
	Buyin  float64 `json:"buyInSum"`
	Buyout float64 `json:"buyOutSum"`
    InGame float64 `json:"inGame"`
	Net    float64 `json:"net"`
}
type Game struct {
    BuyInTotal    float64  `json:"buyInTotal"`
    InGameTotal   float64  `json:"inGameTotal"`
    BuyOutTotal   float64  `json:"buyOutTotal"`
	/*Players  []Player `json:"playersInfos"`*/
}
type Payment struct {
	Payer    string `json:"payer"`
	Reciever string `json:"reciever"`
	Amount   float64 `json:"amount"`
}

func GetInfo(w http.ResponseWriter, r *http.Request) {

    var game Game

    body1, err := io.ReadAll(r.Body)
    if err != nil {
        log.Fatal("err in request", err)
    }
    url := string(body1)

    url = url + "/players_sessions"
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal("err in reaching pokernow", err)
    }
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal("err on line 41", err)
    }

    defer res.Body.Close()
    
    body, _ := io.ReadAll(res.Body)

    err = json.Unmarshal(body, &game)
    if err != nil {
        log.Fatal("err in unmarshal", err)
    }
    
    log.Println(game)

}  
func main() {
   
    http.HandleFunc("/GetInfo", GetInfo)

    http.ListenAndServe(":8080", nil)
    
}
