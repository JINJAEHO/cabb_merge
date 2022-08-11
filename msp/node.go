package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Addr struct {
	NewNode string `json:"node"`
	Address string `json:"address"`
}

var table map[string]string
var drop bool = false

const primary string = "10000"

func main() {
	const sleepDuration = time.Millisecond * 10
	var addr Addr
	table = make(map[string]string)
	table[primary] = "192.168.10.35:10000"

	http.HandleFunc("/newNode", addr.newNodeAccept)
	http.HandleFunc("/tableUpdateAlarm", addr.tableUpdateAlarm)
	go func() {
		http.ListenAndServe(":9999", nil)
	}()

	//일정 시간마다 핑을 날려 노드 연결상태 확인
	for {
		check := false
		for k, v := range table {
			if k == primary {
				continue
			}
			_, err := http.Post("http://"+v+"/pingReq", "text/plain", nil)
			if err != nil {
				delete(table, k)
				fmt.Printf("[Node:%s] is deleted\n", k)
				check = true
				drop = true
			}
		}
		if check {
			_, err := http.Post("http://"+table[primary]+"/newNodeAlarm", "text/plain", nil)
			fmt.Println("리더노드에게 테이블 업데이트 알림")
			if err != nil {
				fmt.Println("업데이트 알림 실패  ", err)
			}
		}
		time.Sleep(sleepDuration)
		//fmt.Println("sleep")
	}
}

func (addr *Addr) newNodeAccept(writer http.ResponseWriter, req *http.Request) {
	json.NewDecoder(req.Body).Decode(addr)
	// 새로운 노드 접속을 리더 노드에게 알림
	_, err := http.Post("http://"+table[primary]+"/newNodeAlarm", "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}
	table[addr.NewNode] = addr.Address + ":" + addr.NewNode
	fmt.Printf("[Node:%s] : %s 추가\n", addr.NewNode, table[addr.NewNode])
	drop = false
}

func (addr *Addr) tableUpdateAlarm(writer http.ResponseWriter, req *http.Request) {
	//리더 노드로 부터 ViewID 값 받아서 1 증가시킴
	var view map[string]int64
	json.NewDecoder(req.Body).Decode(&view)
	for _, v := range view {
		fmt.Printf("recieved ViewID: %d\n", v)
	}
	viewJ, _ := json.Marshal(view)

	//ping을 날려 연결이 끊긴 노드가 있는지 확인
	for k, v := range table {
		if k == primary {
			continue
		}
		_, err := http.Post("http://"+v+"/pingReq", "text/plain", nil)
		if err != nil {
			delete(table, k)
			fmt.Printf("[Node:%s] is deleted\n", k)
		}
	}
	fmt.Println(table)

	//어드레스 테이블을 브로드캐스팅 해줌(리더 노드 제외)
	address, _ := json.Marshal(table)
	for k, v := range table {
		if k == primary {
			continue
		}
		_, err := http.Post("http://"+v+"/getTable", "application/json", bytes.NewBuffer(address))
		//fmt.Println("주소 테이블 브로드캐스팅")
		if err != nil {
			fmt.Println("브로캐스팅 에러  ", err)
		}
	}

	//새로 참가한 노드에게 ViewID 전송
	if _, v := table[addr.NewNode]; v && !drop {
		_, err := http.Post("http://"+table[addr.NewNode]+"/getView", "application/json", bytes.NewBuffer(viewJ))
		fmt.Printf("[%s]에게 viewID 전송\n", addr.NewNode)
		if err != nil {
			fmt.Println("viewID전송 에러  ", err)
		}
	}

	//위 과정이 모두 끝나면 리더노드에게 어드레스 테이블을 응답데이터로 전송
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(table)
	if err != nil {
		fmt.Println("Response 에러  ", err)
	}
}
