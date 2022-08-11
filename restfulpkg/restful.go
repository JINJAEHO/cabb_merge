package main

// go version  go 1.18.4 window/amd64
//Restful API

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

type Args struct {
	Alias   string
	Address string
}
type Request struct {
	Alias   string
	Address string
	Txid    string
}

type Response struct {
	Address    string
	PublicKey  []byte
	PrivateKey []byte
	Check      bool
	Txid       []byte
	SignValue  []byte
}

func main() {
	r := &Request{}
	setRouter(r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func setRouter(r *Request) {
	// if /mdware/MakeWallet으로 요청이 들어오면 r.ConnectWallet 실행
	http.HandleFunc("/MakeWallet", r.GenerateWallet)
	// if /mdware/CheckAddress으로 요청이 들어오면 r.ConnectTransaction 실행
	http.HandleFunc("/CheckAddress", r.CheckAddress)
	// if /mdware/RegisterCareer 요청이 들어오면
	http.HandleFunc("/RegisterCareer", r.RegisterCareer)
	// if /mdware/FindAllTxByAddress 요청이 들어오면
	http.HandleFunc("/GetWalletInfo", r.GetWallet)
	// if /mdware/digitalSignature
	http.HandleFunc("/DigitalSignature", r.DigitalSigniture)

	http.HandleFunc("/refTx", r.findAllList)

	http.HandleFunc("/detailTx", r.findDetail)
}
func (r *Request) GenerateWallet(w http.ResponseWriter, re *http.Request) {
	//------------ Json 으로 들어온 Alias 확인 ( 서버에서 Send)
	headerContentType := re.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		// json 타입이 아니라면
		fmt.Println("Json 타입이 아닙니다!!")
	}
	decoder := json.NewDecoder(re.Body)
	var request Request
	err := decoder.Decode(&request) // request Body에 들어있는 json 데이터를 해독하고 저장
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(request.Alias, "요청받은 Alias 입니다.")

	// -----------------------------------------JSON 해독 끝 -------------------------
	// ------------------------- RPC 서버 연결 ---------------------
	Client, err := rpc.Dial("tcp", "192.168.10.158:9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Client.Close()
	response := new(Response) // 연결후 return을 받기 위해 빈 바구니 생성
	err = Client.Call("RpcServer.MakeNewWallet", request.Alias, response)
	if err != nil {
		fmt.Println(err, "Client.Call 에서 에러가 났음 ")
		return
	}
	fmt.Println(request.Alias, "님의 지갑의 Address 입니다 ", response.Address, response.PrivateKey, "님의 PrivateKey 입니다. 보안에 유의하세요", response.PublicKey, "님의 PublicKey 입니다")
	// Wallet.go에서 받아온 데이터 요청한 서비스로 다시 돌려주기
	// 돌려주기 위해서 Json Parsing
	PrivateKey := hex.EncodeToString(response.PrivateKey)
	fmt.Println(PrivateKey, "PrivateKey")
	PublicKey := hex.EncodeToString(response.PublicKey)
	fmt.Println(PublicKey, "PublicKey")
	// fmt.Println(PrivateKey, "PrvateKey")
	// fmt.Println(PublicKey, "PublicKey")
	value := map[string]interface{}{
		"Alias":      request.Alias,
		"Address":    response.Address,
		"PublicKey":  PublicKey,
		"PrivateKey": PrivateKey,
	}
	// json_data, err := json.Marshal(value) // Parsing 완료
	// fmt.Println(json_data, "json 파싱한 후 데이터 ")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// 주소 검증 후 findAllTxByAddress 요청 보내기
func (r *Request) CheckAddress(w http.ResponseWriter, re *http.Request) {
	// 주소 검증.
	headerContentType := re.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		// json 타입이 아니라면
		fmt.Println("Json 타입이 아닙니다!!")
	}
	decoder := json.NewDecoder(re.Body)
	var request Request
	err := decoder.Decode(&request) // request Body에 들어있는 json 데이터를 해독하고 저장
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(request.Address, "요청받은 Address 입니다.")
	Client, err := rpc.Dial("tcp", "192.168.10.158:9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Client.Close()
	response := new(Response)
	err = Client.Call("RpcServer.CheckAddress", request.Address, response)
	if err != nil {
		fmt.Println(err)
		return
	}
	if response.Check {
		fmt.Println("존재하는 지갑주소입니다")
		value := map[string]interface{}{
			"Address": request.Address,
		}
		json_data, err := json.Marshal(value)
		res, err := http.Post("http://192.168.10.239:3000/FindAllTxByAddress", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			fmt.Println(err)
			return
		} // 받아온 Txs를 돌려줌
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)

	} else {
		fmt.Println("존재하지 않는 지갑주소입니다.")
	}
}

// 지갑 주소를 주고 그 주소에 해당하는 지갑을 받아오기
func (r *Request) GetWallet(w http.ResponseWriter, re *http.Request) {
	headerContentType := re.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		// json 타입이 아니라면
		fmt.Println("Json 타입이 아닙니다!!")
	}
	decoder := json.NewDecoder(re.Body)
	var request Request
	err := decoder.Decode(&request) // request Body에 들어있는 json 데이터를 해독하고 저장
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(request.Address, "요청받은 Address 입니다.")
	Client, err := rpc.Dial("tcp", "192.168.10.158:9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Client.Close()
	response := new(Response)
	err = Client.Call("RpcServer.GetWallet", request.Address, response)
	if err != nil {
		fmt.Println(err)
		return
	}
	value := map[string]string{
		"Address":    response.Address,
		"PublicKey":  string(response.PublicKey),
		"PrivateKey": string(response.PrivateKey),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}

// Digital Signature Function
func (r *Request) DigitalSigniture(w http.ResponseWriter, re *http.Request) {
	headerContentType := re.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		// json 타입이 아니라면
		fmt.Println("Json 타입이 아닙니다!!")
	}
	decoder := json.NewDecoder(re.Body)
	var request Request
	err := decoder.Decode(&request) // request Body에 들어있는 json 데이터를 해독하고 저장
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("요청받은 Address 와 Txid 입니다.", request.Address, request.Txid)
	Client, err := rpc.Dial("tcp", "192.168.10.158:9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Client.Close()
	response := new(Response)
	err = Client.Call("RpcServer.Signature", request, response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sign Value: ", response.SignValue)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 경력 등록 요청을 받으면 Apply/Career로 보내기
func (r *Request) RegisterCareer(w http.ResponseWriter, re *http.Request) {
	Res, err := http.Post("http://192.168.10.239:5000/newBlk", "application/json", re.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := make(map[string][32]byte)
	json.NewDecoder(Res.Body).Decode(&s)

	// Txid 만 돌려줌
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

type listBody struct {
	TxID    [][32]byte `json:"txID"`
	Career  []string   `json:"career"`
	Company []string   `json:"company"`
}

func (r *Request) findAllList(w http.ResponseWriter, re *http.Request) {
	Res, err := http.Post("http://192.168.10.239:5000/refTx", "application/json", re.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var listStruct listBody
	json.NewDecoder(Res.Body).Decode(&listStruct)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listStruct)
}

type JsonDetailResponse struct {
	Hash      [32]byte `json:"blockID"`
	Data      string   `json:"Data"`
	Timestamp string   `json:"Timestamp"`
	Txid      [32]byte `json:"Txid"`
	Applier   string   `json:"Applier"`
	Company   string   `json:"Company"`
	Career    string   `json:"Career"`
	Job       string   `json:"Job"`
	Proof     string   `json:"Proof"`
}

func (r *Request) findDetail(w http.ResponseWriter, re *http.Request) {
	Res, err := http.Post("http://192.168.10.239:5000/detailTx", "application/json", re.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var detailStruct JsonDetailResponse
	json.NewDecoder(Res.Body).Decode(&detailStruct)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detailStruct)
}
