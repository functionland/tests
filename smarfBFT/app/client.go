package app

import (
	"github.com/SmartBFT-Go/consensus/pkg/types"
	"github.com/SmartBFT-Go/consensus/smartbftprotos"
	"go.uber.org/zap"
	"time"
)

var fastConfig = types.Configuration{
	RequestBatchMaxCount:          10,
	RequestBatchMaxBytes:          10 * 1024 * 1024,
	RequestBatchMaxInterval:       10 * time.Millisecond,
	IncomingMessageBufferSize:     200,
	RequestPoolSize:               40,
	RequestForwardTimeout:         500 * time.Millisecond,
	RequestComplainTimeout:        2 * time.Second,
	RequestAutoRemoveTimeout:      3 * time.Minute,
	ViewChangeResendInterval:      5 * time.Second,
	ViewChangeTimeout:             1 * time.Minute,
	LeaderHeartbeatTimeout:        1 * time.Minute,
	LeaderHeartbeatCount:          10,
	NumOfTicksBehindBeforeSyncing: 10,
	CollectTimeout:                200 * time.Millisecond,
	LeaderRotation:                false,
	RequestMaxBytes:               10 * 1024,
	RequestPoolSubmitTimeout:      5 * time.Second,
}

type BFTClient struct {
	ID              uint64
	Delivered       chan *AppRecord
	Consensus       *consensus.Consensus
	Setup           func()
	Node            *Node
	logLevel        zap.AtomicLevel
	latestMD        *smartbftprotos.ViewMetadata
	lastDecision    *types.Decision
	clock           *time.Ticker
	heartbeatTime   chan time.Time
	viewChangeTime  chan time.Time
	secondClock     *time.Ticker
	logger          *zap.SugaredLogger
	lastRecord      lastRecord
	verificationSeq uint64
	messageLost     func(*smartbftprotos.Message) bool
}

func (*BFTClient) Put(key, val string) {

}

func (*BFTClient) Get() {

}

func (*BFTClient) Remove() {

}

func (*BFTClient) Size() {

}

func (*BFTClient) keySet() {

}
