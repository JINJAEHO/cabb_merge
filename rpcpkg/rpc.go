package main

// gRPC server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	_ "errors"
	"fmt"
	"net"
	"net/rpc"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// type Return struct{} // 월렛 주소 반환하는 값을 담은 바구니

// type Response struct {
// 	Address    string
// 	PublicKey  []byte
// 	PrivateKey []byte
// }

// type Wallet struct {
// 	PublicKey  []byte
// 	PrivateKey ecdsa.PrivateKey
// 	Address    string
// 	Alias      string
// }
// type Wallets struct {
// 	Wts map[string]*Wallet
// }

// // ---------------------------------------------------------------- Functions --------------------------------
// func newKeyPair() (ecdsa.PrivateKey, []byte) {
// 	curve := elliptic.P256()
// 	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
// 	publicKey := privateKey.PublicKey
// 	bPublicKey := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
// 	return *privateKey, bPublicKey
// }
// func HashPublicKey(publicKey []byte) []byte {
// 	publicSHA256 := sha256.Sum256(publicKey)
// 	RIPEMD160Hasher := ripemd160.New()
// 	_, _ = RIPEMD160Hasher.Write(publicSHA256[:])

// 	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
// 	return publicRIPEMD160
// }
// func encodeAddress(publicREPEMD160 []byte) string {
// 	version := byte(0x00)
// 	s := base58.CheckEncode(publicREPEMD160, version)
// 	return s
// }
// func MakeWallet(Alias string) *Wallet {
// 	ws := new(Wallets)
// 	w := new(Wallet)
// 	PrivateKey, PublicKey := newKeyPair()
// 	w.PrivateKey = PrivateKey
// 	w.PublicKey = PublicKey
// 	// 유효성검사 Address 가 존재하지 않는다면
// 	if ws.CheckWalletAddr(encodeAddress(HashPublicKey(PublicKey))) {
// 		w.Address = encodeAddress(HashPublicKey(PublicKey))
// 		//지갑을 Wallets에 저장
// 		ws.PutWallet(w, w.Address)
// 	} else {
// 		// 만약 이미 존재하는 Address를 만들었다면 , 다시 키쌍을 만들기
// 		PrivateKey, PublicKey = newKeyPair()
// 		w.PrivateKey = PrivateKey
// 		w.PublicKey = PublicKey
// 		w.Address = encodeAddress(HashPublicKey(PublicKey))
// 	}
// 	w.Alias = Alias
// 	return w
// }
// func NewWallets() {
// 	Ws := &Wallets{}
// 	Ws.Wts = make(map[string]*Wallet)

// }
// func (ws *Wallets) PutWallet(Wallet *Wallet, encodedAddress string) {
// 	// Wallet을 Wallets에 저장하기
// 	ws.Wts[encodedAddress] = Wallet

// }

// func (ws *Wallets) CheckWalletAddr(encodedAddress string) bool {
// 	if ws.Wts[encodedAddress] != nil {
// 		return false
// 	}
// 	return true
// }

// // -- 실제 Response 해주는 Functions
// func (r *Return) SendWallet(Alias string, response *Response) error {
// 	Wallet := MakeWallet(Alias)
// 	response.Address = Wallet.Address
// 	response.PublicKey = Wallet.PublicKey
// 	response.PrivateKey = Wallet.PrivateKey.D.Bytes()
// 	return nil
// }

var wallets = make(map[string]*Wallet)

type RpcServer struct{}

type Reply struct {
	Alias      string
	Address    string
	PublicKey  []byte
	PrivateKey []byte
	Check      bool
	SignValue  []byte
}

// ------- Actual Responsable Functions -------------------------------

func (wRPC *RpcServer) MakeNewWallet(Alias string, reply *Reply) error {
	prvKey, pubKey := NewKeyPair()
	w := MakeWallet(&prvKey, pubKey, Alias)
	reply.Address = w.Address
	reply.PrivateKey = w.PrivateKey
	fmt.Println(reply.PrivateKey, "reply.PrvKey 입니다")
	reply.PublicKey = w.PublicKey
	fmt.Println(reply.PublicKey, "reply.PubKey 입니다")
	reply.Alias = w.Alias
	return nil
}
func (wRPC *RpcServer) CheckAddress(Address string, reply *Reply) error {
	// 주소가 존재한다면
	if wallets[Address] != nil {
		reply.Check = true
	} else {
		reply.Check = false
	}
	return nil
}

func (wRPC *RpcServer) GetWallet(Address string, reply *Reply) error {

	w := wallets[Address]
	reply.PrivateKey = w.PrivateKey
	reply.PublicKey = w.PublicKey
	return nil
}

func (wRCP *RpcServer) Signature(request *Request, reply *Reply) error {
	wallet := wallets[request.Address]
	Txid := []byte(request.Txid)
	SignValue, _ := ecdsa.SignASN1(rand.Reader, wallet.ecdsaPrviateKey, Txid)
	reply.SignValue = SignValue
	return nil
}

// ----------------- End of Actual Responsable Functions ----------------

type Args struct {
	Alias   string
	Address string
}

type Wallet struct {
	PrivateKey      []byte
	PublicKey       []byte
	Address         string
	Alias           string
	ecdsaPrviateKey *ecdsa.PrivateKey
}

type Request struct {
	Txid    string
	Address string
}

func MakeWallet(prvkey *ecdsa.PrivateKey, pubkey []byte, alias string) *Wallet {
	w := &Wallet{}
	publicRIPEMD160 := HashPubKey(pubkey)
	version := byte(0x00)
	Address := base58.CheckEncode(publicRIPEMD160, version)
	w.PrivateKey = prvkey.D.Bytes()
	w.ecdsaPrviateKey = prvkey
	w.PublicKey = pubkey
	w.Address = Address
	w.Alias = alias
	// walltes 에 방금 만들어진 wallet을 넣기
	wallets[w.Address] = w
	return w
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	prvKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := prvKey.PublicKey
	bpubKey := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)

	return *prvKey, bpubKey
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func (Wallet *Wallet) PrintInfo() {
	fmt.Printf("Alias : %s\n", Wallet.Alias)
	fmt.Printf("Address : %s\n", Wallet.Address)
	fmt.Printf("PublicKey : %x\n", Wallet.PublicKey)
	fmt.Printf("PrivateKey : %s\n", Wallet.PrivateKey)
}

// -------------------- main ----------------------------------------------------

func main() {
	rpc.Register(new(RpcServer))
	In, err := net.Listen("tcp", ":9000")
	fmt.Println(In, "In 입니다")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer In.Close()
	for {
		conn, err := In.Accept()
		fmt.Println(conn, err, "In.Accept 입니다")
		if err != nil {
			continue
		}
		defer conn.Close()

		go rpc.ServeConn(conn)
	}
}
