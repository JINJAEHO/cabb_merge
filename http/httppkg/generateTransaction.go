package httppkg

import (
	"bytes"
	_ "bytes"
	"encoding/json"
	"fmt"
	"http/txpkg"
	"net/http"
	_ "net/http"
)

//Json 타입으로 리턴해주기 위한 구조체
type JsonResponse struct {
	Txid [32]byte `json:"txid"`
}

type ForSign struct {
	Address string
	Txid    string
}

type ResSing struct {
	SignValue []byte
}

var Txs *txpkg.Txs

// Generate Transaction
func ApplyCareer(w http.ResponseWriter, req *http.Request) {
	var body Request
	decoder := json.NewDecoder(req.Body)
	//decoder.DisallowUnknownFields()
	decoder.Decode(&body)

	//트랜잭션 생성
	T := txpkg.NewTx(body.Applier, body.Company, body.Career, body.Payment, body.Job, body.Proof, body.Address)

	//전자서명 생성
	signBody := ForSign{Address: body.Address, Txid: string(T.TxID[:])}
	jsonSign, _ := json.Marshal(signBody)
	SignRes, err := http.Post("http://192.168.10.99:3000/DigitalSignature", "application/json", bytes.NewBuffer(jsonSign))
	if err != nil {
		fmt.Println(err)
		return
	}

	var HashedTxid ResSing
	json.NewDecoder(SignRes.Body).Decode(&HashedTxid)
	fmt.Printf("전자서명: %x\n", HashedTxid.SignValue)

	T.Sign = HashedTxid.SignValue
	Txs.AddTx(T)
	T.PrintTx()
	fmt.Println("Tx-TxID: ", T.TxID)

	jsonForPBFT, _ := json.Marshal(T)
	PBFT_Res, err := http.Post("http://192.168.10.35:10000/req", "application/json", bytes.NewBuffer(jsonForPBFT))
	if err != nil {
		fmt.Println("합의 통신 실패 ", err)
	}

	var PBFT_check bool
	json.NewDecoder(PBFT_Res.Body).Decode(&PBFT_check)

	if !PBFT_check {
		fmt.Println("합의 실패")
		return
	}

	jsonResponse := JsonResponse{Txid: T.TxID}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

	// value := map[string]string{
	// 	"txID": hex.EncodeToString(T.TxID[:]),
	// 	"data": body.Data}
	// json_data, _ := json.Marshal(value)
	// resp, err := http.Post("http://localhost:3000/newBlk", "application/json", bytes.NewBuffer(json_data))
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// // Response 체크.
	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err == nil {
	// 	str := string(respBody)
	// 	println(str)
	// }
	// var response = JsonResponse{Address: body.Address, Txid: Txid}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}
