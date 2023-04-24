package types

import proto "code.zetaprotocol.io/vega/protos/vega/commands/v1"

type NodeVoteType = proto.NodeVote_Type

const (
	NodeVoteTypeUnspecified             NodeVoteType = proto.NodeVote_TYPE_UNSPECIFIED
	NodeVoteTypeStakeDeposited          NodeVoteType = proto.NodeVote_TYPE_STAKE_DEPOSITED
	NodeVoteTypeStakeRemoved            NodeVoteType = proto.NodeVote_TYPE_STAKE_REMOVED
	NodeVoteTypeFundsDeposited          NodeVoteType = proto.NodeVote_TYPE_FUNDS_DEPOSITED
	NodeVoteTypeSignerAdded             NodeVoteType = proto.NodeVote_TYPE_SIGNER_ADDED
	NodeVoteTypeSignerRemoved           NodeVoteType = proto.NodeVote_TYPE_SIGNER_REMOVED
	NodeVoteTypeBridgeStopped           NodeVoteType = proto.NodeVote_TYPE_BRIDGE_STOPPED
	NodeVoteTypeBridgeResumed           NodeVoteType = proto.NodeVote_TYPE_BRIDGE_RESUMED
	NodeVoteTypeAssetListed             NodeVoteType = proto.NodeVote_TYPE_ASSET_LISTED
	NodeVoteTypeAssetLimitsUpdated      NodeVoteType = proto.NodeVote_TYPE_LIMITS_UPDATED
	NodeVoteTypeStakeTotalSupply        NodeVoteType = proto.NodeVote_TYPE_STAKE_TOTAL_SUPPLY
	NodeVoteTypeSignerThresholdSet      NodeVoteType = proto.NodeVote_TYPE_SIGNER_THRESHOLD_SET
	NodeVoteTypeGovernanceValidateAsset NodeVoteType = proto.NodeVote_TYPE_GOVERNANCE_VALIDATE_ASSET
)
