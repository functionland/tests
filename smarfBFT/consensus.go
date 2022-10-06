package main

import (
	"github.com/SmartBFT-Go/consensus/internal/bft"
	"github.com/SmartBFT-Go/consensus/pkg/types"
	protos "github.com/SmartBFT-Go/consensus/smartbftprotos"
	"go.uber.org/zap"
)

func RestoreConsensusState() {
	prePrepare := &protos.PrePrepare{
		Proposal: &protos.Proposal{
			Header:  []byte{1},
			Payload: []byte{1},
			Metadata: bft.MarshalOrPanic(&protos.ViewMetadata{
				DecisionsInView: 0,
				LatestSequence:  0,
				ViewId:          1,
			}),
			VerificationSequence: 100,
		},
		Seq:  200,
		View: 300,
	}

	expectedInFlightProposal := &types.Proposal{
		VerificationSequence: int64(prePrepare.Proposal.VerificationSequence),
		Metadata:             prePrepare.Proposal.Metadata,
		Payload:              prePrepare.Proposal.Payload,
		Header:               prePrepare.Proposal.Header,
	}

	proposedRecord := &protos.SavedMessage{
		Content: &protos.SavedMessage_ProposedRecord{
			ProposedRecord: &protos.ProposedRecord{
				PrePrepare: prePrepare,
				Prepare: &protos.Prepare{
					Seq:  200,
					View: 300,
				},
			},
		},
	}

	preparedProof := &protos.SavedMessage{
		Content: &protos.SavedMessage_Commit{
			Commit: &protos.Message{
				Content: &protos.Message_Commit{
					Commit: &protos.Commit{
						Seq:  200,
						View: 300,
						Signature: &protos.Signature{
							Signer: 11,
						},
					},
				},
			},
		},
	}

	basicLog, err := zap.NewDevelopment()
	log := basicLog.Sugar()

	//proposed
	var (
		expectedPhase       = bft.PROPOSED
		expectedViewNumber  = 300
		expectedProposalSeq = 200
		WALContent          = [][]byte{bft.MarshalOrPanic(proposedRecord)}
	)

	state := &bft.PersistedState{
		Entries:          testCase.WALContent,
		Logger:           log,
		InFlightProposal: &bft.InFlightData{},
	}

	view := &bft.View{
		Number:           300,
		ProposalSequence: testCase.proposalSeqViewInitializedWith,
	}

	err = state.Restore(view)

}
