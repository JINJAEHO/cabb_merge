package httppkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type findReqBody struct {
	Address string `json:"address"`
}

type resBody struct {
	TxID    [][32]byte `json:"txID"`
	Career  []string   `json:"career"`
	Company []string   `json:"company"`
}

func FindAllbyAddr(w http.ResponseWriter, req *http.Request) {
	var body findReqBody
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	//에러 체크
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Req Body Address: %s\n", body.Address)

	//테스트용 더미 데이터
	// txs := txpkg.CreateTxDB()
	// gb := blockpkg.GenesisBlock()
	// bs := blockpkg.NewBlockchain()
	// prev := gb.Hash
	// height := gb.Height + 1
	// address := "address123"
	// for j := 0; j < 10; j++ {
	// 	tx := txpkg.NewTx("user_"+fmt.Sprint(j), "company_"+fmt.Sprint(j), fmt.Sprint(j)+"개월", "card", "블록체인 개발자", "proof.png", address)
	// 	txs.AddTx(tx)
	// 	b := blockpkg.NewBlock(prev, height, tx.TxID, "data")
	// 	bs.AddBlock(b)
	// 	prev = b.Hash
	// 	height = b.Height + 1
	// }

	list := Txs.FindTxByAddr(body.Address, BlkChain)

	for _, v := range list {
		fmt.Printf("txID: %x\n", v.TxID)
	}

	res := &resBody{}
	for i := 0; i < len(list); i++ {
		res.TxID = append(res.TxID, list[i].TxID)
		res.Career = append(res.Career, string(list[i].Career))
		res.Company = append(res.Company, string(list[i].Company))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
