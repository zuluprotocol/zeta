package v1

import types "code.zetaprotocol.io/vega/protos/vega"

func ProposalSubmissionFromProposal(p *types.Proposal) *ProposalSubmission {
	return &ProposalSubmission{
		Reference: p.Reference,
		Terms:     p.Terms,
	}
}
