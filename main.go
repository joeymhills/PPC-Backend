package main

import (
    "log"
    "errors"
    "net/http"
    "io"
    "encoding/json"
    "sort"
    "os"
)

type Player struct {
    Id     string `json:"id"` 
    Names  []string `json:"names"`
    Buyin  int32 `json:"buyInSum"`
    Buyout int32 `json:"buyOutSum"`
    InGame int32 `json:"inGame"`
    Net    int32 `json:"net"`
}
type Game struct {
    BuyInTotal    int32  `json:"buyInTotal"`
    InGameTotal   int32  `json:"inGameTotal"`
    BuyOutTotal   int32  `json:"buyOutTotal"`
    Players       *json.RawMessage   `json:"playersInfos"`
}

type Players struct {
    TagName string
    players []Player
}

type Payment struct {
	Payer    string `json:"payer"`
	Reciever string `json:"reciever"`
	Amount   int32 `json:"amount"`
}
func abs(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}

func minInt(a, b int32) int32 {
  if a < b {
    return a
  }
  return b
}

func sortPlayers(p []Player) []Player {
sort.Slice(p, func(i, j int) bool {
  return p[i].Net < p[j].Net
})
return p
}

func reverseSortPlayers(p []Player) []Player {
sort.Slice(p, func(i, j int) bool {
  return p[i].Net < p[j].Net
})
return p
}
/*
//write a function that takes in playerArray and maps it to PlayersInfos
func GetPlayersInfos() {
    var playersInfos PlayersInfos
    err := json.Unmarshal([]byte(playerArray), &playersInfos)
    if err != nil {
        log.Fatal("err in unmarshal", err)
    }
    log.Println(playersInfos)
}
*/

func CalculatePayments(playersarray []Player) ([]Payment, error) {
    var payments []Payment
    var err error = nil
    var payers []Player
    var recievers []Player

    for _, player := range playersarray {
        if player.Net < 0 {
            payers = append(payers, player)
        } else if player.Net > 0 {
            recievers = append(recievers, player)
        }
    }
    payers = sortPlayers(payers)
    recievers = reverseSortPlayers(recievers)

    var payer *Player
    var reciever *Player
    var payment Payment
    amount := int32(0)

    for (len(payers) > 0) && (len(recievers) > 0) {
        payer = &payers[0]
        reciever = &recievers[0]
        amount = minInt(abs(payer.Net), abs(reciever.Net))

        payment = Payment{payer.Names[-0], reciever.Names[-0], amount}
        
        payer.Net = payer.Net + amount
        reciever.Net = reciever.Net - amount
        
        if payer.Net == 0 {
            payers = payers[1:]
        }
        if reciever.Net == 0 {
            recievers = recievers[1:]
        }
        payments = append(payments, payment)

    }

    if len(payers) > 0 || len(recievers) > 0 {
        err = errors.New("Error in calculating payments")
    }


    return payments, err 
}

func GetInfo(w http.ResponseWriter, r *http.Request) {

    var game Game
    body1, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    url := string(body1)
    url = url + "/players_sessions"
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer res.Body.Close()
    
    body, _ := io.ReadAll(res.Body)

    err = json.Unmarshal(body, &game)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    

    var playermap map[string]Player
    var playerarray []Player
    err = json.Unmarshal(*game.Players, &playermap)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
     for _, value := range playermap {
        playerarray = append(playerarray, value)
    }

    if len(playerarray) == 0 {
        err = errors.New("Game does not exist.")
        http.Error(w, err.Error(), 500)
        return
    } 

    payments, err := CalculatePayments(playerarray)
    
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(payments)

}  
func main() {

    log.Println("Listening and serving CONCURRENTLY")
    
    http.HandleFunc("/GetInfo", GetInfo)
    
    PORT := os.Getenv("PORT")
    if PORT == "" {
        PORT = "3333"
    }
    http.ListenAndServe("0.0.0.0:"+PORT, nil)

}
