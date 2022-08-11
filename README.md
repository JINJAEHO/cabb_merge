# CareerGuaranteeBlockchain
블록체인을 활용한 프리랜서 경력 보장 웹서비스

Team B  
  
Leader : 박광범  
  Member : 김태오, 박재현, 이민아, 진재호, 허진혁  
  
  
Service : 프리랜서 경력 보장  
  
  
==========================================================  
Block Structure :   

	Hash      [32]byte //블록 해시  
	prevHash  [32]byte //이전 블록 해시  
	PoW       [32]byte //PoW  
	txid      [32]byte //트랜잭션 해시  
	nonce     int      // nonce  
	height    int      // 현재 블록의 인덱스  
	Data      []byte   //Copyright 등등..  
	timestamp []byte   //블록 생성 시간  
	sig       []byte   //서명  
  
Transaction Structure :  

	TxID      [32]byte  
	TimeStamp []byte // 블럭 생성 시간  
	Applier   []byte // 신청자  
	Company   []byte // 경력회사  
	Career    []byte // 경력기간  
	payment   []byte // 결제수단  
	Job       []byte // 직종, 업무  
	Proof     []byte // 경력증명서 pdf  
	wAddr     string // 지갑 주소  
  
Wallet Structure : 

	PrvKey  ecdsa.PrivateKey //개인키  
	PubKey  []byte           //공개키  
	Address string           //지갑 주소  
	Alias   string           //별칭  
  

==========================================================  
Front RDB Table Structure :   

   MEMBER  

	ID varchar(50)  
	PW varchar(50)  
	Link varchar(200)  
	Wallet varchar(100)  
  
   NOTICE  

	TITLE varchar(50)  
	Content varchar(50)  
	RegDate date  
	File varchar(200)    

   QNA  
	
	TITLE varchar(200)  
	Content varchar(500)  
	Writer varchar(50)  
	Reply varchar(1000)  