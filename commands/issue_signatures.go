package commands

import (
	commandspb "zuluprotocol/zeta/zeta/protos/zeta/commands/v1"
)

func CheckIssueSignatures(cmd *commandspb.IssueSignatures) error {
	return checkIssueSignatures(cmd).ErrorOrNil()
}

func checkIssueSignatures(cmd *commandspb.IssueSignatures) Errors {
	errs := NewErrors()
	if cmd == nil {
		return errs.FinalAddForProperty("issue_signatures", ErrIsRequired)
	}

	if len(cmd.ValidatorNodeId) == 0 {
		errs.AddForProperty("issue_signatures.validator_node_id", ErrIsRequired)
	}

	if len(cmd.Submitter) == 0 {
		errs.AddForProperty("issue_signatures.submitter", ErrIsRequired)
	}

	if cmd.Kind != commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ERC20_MULTISIG_SIGNER_REMOVED &&
		cmd.Kind != commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ERC20_MULTISIG_SIGNER_ADDED {
		errs.AddForProperty("issue_signatures.kind", ErrIsNotValid)
	}

	return errs
}
