package commands_test

import (
	"testing"

	"zuluprotocol/zeta/commands"
	types "zuluprotocol/zeta/protos/zeta"
	commandspb "zuluprotocol/zeta/protos/zeta/commands/v1"
	"github.com/stretchr/testify/assert"
)

func TestCheckProposalSubmissionForNewFreeform(t *testing.T) {
	t.Run("Submitting a new freeform change without new freeform fails", testNewFreeformChangeSubmissionWithoutNewFreeformFails)
	t.Run("Submitting a new freeform proposal without rational URL and hash fails", testNewFreeformProposalSubmissionWithoutRationalURLandHashFails)
}

func testNewFreeformChangeSubmissionWithoutNewFreeformFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &types.ProposalTerms{
			Change: &types.ProposalTerms_NewFreeform{},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_freeform"), commands.ErrIsRequired)
}

func testNewFreeformProposalSubmissionWithoutRationalURLandHashFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &types.ProposalTerms{
			Change: &types.ProposalTerms_NewFreeform{},
		},
		Rationale: &types.ProposalRationale{},
	})

	assert.Contains(t, err.Get("proposal_submission.rationale.description"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.rationale.title"), commands.ErrIsRequired)
}
