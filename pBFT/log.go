package main

import (
	"fmt"
	"log"
)

func LogMsg(msg interface{}) {
	switch msg.(type) {
	case *RequestMsg:
		reqMsg := msg.(*RequestMsg)
		fmt.Printf("[REQUEST] TxID : %x, TimeStamp : %s, WAdr : %s, Applier : %s , Company : %s, Career : %s, Job : %s, Sign : %s \n",
			reqMsg.TxID[:], string(reqMsg.TimeStamp), reqMsg.WAddr, reqMsg.Applier, reqMsg.Company, reqMsg.Career, reqMsg.Job, reqMsg.Sign[:])
	case *PrePrepareMsg:
		prePrepareMsg := msg.(*PrePrepareMsg)
		fmt.Printf("[PREPREPARE] TxID : %x, ViewID : %d, SequenceID : %d, Digest : %s\n", prePrepareMsg.RequestMsg.TxID[:], prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest)

	case *VoteMsg:
		voteMsg := msg.(*VoteMsg)
		if voteMsg.MsgType == PrepareMsg {
			// log.Printf("[PREPARE] NodeID: %s\n", voteMsg.NodeID)
		} else if voteMsg.MsgType == CommitMsg {
			// log.Printf("[COMMIT] NodeID: %s\n", voteMsg.NodeID)
		}
	}
}

func LogStage(stage string, isDone bool) {
	if isDone {
		// log.Printf("[STAGE-DONE] %s\n", stage)
		if stage == "Reply" {
			log.Println("consensus count : ", j-1)
		}
	} else {
		// log.Printf("[STAGE-BEGIN] %s\n", stage)
	}

}
