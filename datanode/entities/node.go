// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package entities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	v2 "code.zetaprotocol.io/vega/protos/data-node/api/v2"
	"code.zetaprotocol.io/vega/protos/vega"
	eventspb "code.zetaprotocol.io/vega/protos/vega/events/v1"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/encoding/protojson"
)

type _Node struct{}

type NodeID = ID[_Node]

type Node struct {
	ID                NodeID
	PubKey            ZetaPublicKey       `db:"zeta_pub_key"`
	TmPubKey          TendermintPublicKey `db:"tendermint_pub_key"`
	EthereumAddress   EthereumAddress
	InfoURL           string
	Location          string
	StakedByOperator  decimal.Decimal
	StakedByDelegates decimal.Decimal
	StakedTotal       decimal.Decimal
	MaxIntendedStake  decimal.Decimal
	PendingStake      decimal.Decimal
	EpochData         EpochData
	Status            NodeStatus
	Delegations       []Delegation  `json:""`
	RewardScore       *RewardScore  `json:""`
	RankingScore      *RankingScore `json:""`
	Name              string
	AvatarURL         string
	TxHash            TxHash
	ZetaTime          time.Time
}

type ValidatorUpdateAux struct {
	Added           bool
	EpochSeq        uint64
	ZetaPubKeyIndex uint32
	TxHash          TxHash
}

type EpochData struct {
	*zeta.EpochData
}

type RewardScore struct {
	RawValidatorScore   decimal.Decimal     `json:"raw_validator_score"`
	PerformanceScore    decimal.Decimal     `json:"performance_score"`
	MultisigScore       decimal.Decimal     `json:"multisig_score"`
	ValidatorScore      decimal.Decimal     `json:"validator_score"`
	NormalisedScore     decimal.Decimal     `json:"normalised_score"`
	ValidatorNodeStatus ValidatorNodeStatus `json:"validator_node_status,string"`
	TxHash              TxHash              `json:"tx_hash"`
	ZetaTime            time.Time           `json:"zeta_time"`
	EpochSeq            uint64
}

type RewardScoreAux struct {
	NodeID   NodeID
	EpochSeq uint64
}

type RankingScore struct {
	StakeScore       decimal.Decimal     `json:"stake_score"`
	PerformanceScore decimal.Decimal     `json:"performance_score"`
	PreviousStatus   ValidatorNodeStatus `json:"previous_status,string"`
	Status           ValidatorNodeStatus `json:",string"`
	VotingPower      uint32              `json:"voting_power"`
	RankingScore     decimal.Decimal     `json:"ranking_score"`
	TxHash           TxHash              `json:"tx_hash"`
	ZetaTime         time.Time           `json:"zeta_time"`
	EpochSeq         uint64
}

type RankingScoreAux struct {
	NodeID   NodeID
	EpochSeq uint64
}

type NodeSet struct {
	Total    uint32
	Inactive uint32
	Promoted []string
	Demoted  []string
	Maximum  uint32
}

type NodeData struct {
	StakedTotal     decimal.Decimal
	TotalNodes      uint32
	InactiveNodes   uint32
	TendermintNodes NodeSet
	ErsatzNodes     NodeSet
	PendingNodes    NodeSet
	Uptime          float64
	ZetaTime        time.Time
}

func NodeFromValidatorUpdateEvent(evt eventspb.ValidatorUpdate, txHash TxHash, zetaTime time.Time) (Node, ValidatorUpdateAux, error) {
	return Node{
			ID:              NodeID(evt.NodeId),
			PubKey:          ZetaPublicKey(evt.VegaPubKey),
			TmPubKey:        TendermintPublicKey(evt.TmPubKey),
			EthereumAddress: EthereumAddress(evt.EthereumAddress),
			InfoURL:         evt.InfoUrl,
			Location:        evt.Country,
			Name:            evt.Name,
			AvatarURL:       evt.AvatarUrl,
			TxHash:          txHash,
			ZetaTime:        zetaTime,

			// Not present in the event
			Status:            NodeStatusValidator, // This was the default value in the legacy store code
			StakedByOperator:  decimal.Zero,
			StakedByDelegates: decimal.Zero,
			StakedTotal:       decimal.Zero,
			MaxIntendedStake:  decimal.Zero,
			PendingStake:      decimal.Zero,
			EpochData:         EpochData{},
			Delegations:       []Delegation{},
			RewardScore:       nil,
			RankingScore:      nil,
		}, ValidatorUpdateAux{
			Added:           evt.Added,
			EpochSeq:        evt.EpochSeq,
			ZetaPubKeyIndex: evt.VegaPubKeyIndex,
			TxHash:          txHash,
		}, nil
}

func ValidatorNodeStatusFromString(status string) ValidatorNodeStatus {
	switch status {
	case "tendermint":
		return ValidatorNodeStatusTendermint
	case "ersatz":
		return ValidatorNodeStatusErsatz
	case "pending":
		return ValidatorNodeStatusPending
	case "unspecified":
		fallthrough
	default: // Is this appropriate behavior? Should we error on the default case?
		return ValidatorNodeStatusUnspecified
	}
}

func RankingScoreFromRankingEvent(evt eventspb.ValidatorRankingEvent, txHash TxHash, zetaTime time.Time) (RankingScore, RankingScoreAux, error) {
	stakeScore, err := decimal.NewFromString(evt.StakeScore)
	if err != nil {
		return RankingScore{}, RankingScoreAux{}, err
	}

	performanceScore, err := decimal.NewFromString(evt.PerformanceScore)
	if err != nil {
		return RankingScore{}, RankingScoreAux{}, err
	}

	rankingScore, err := decimal.NewFromString(evt.RankingScore)
	if err != nil {
		return RankingScore{}, RankingScoreAux{}, err
	}

	epochSeq, err := strconv.ParseUint(evt.EpochSeq, 10, 64)
	if err != nil {
		return RankingScore{}, RankingScoreAux{}, err
	}

	return RankingScore{
			StakeScore:       stakeScore,
			PerformanceScore: performanceScore,
			PreviousStatus:   ValidatorNodeStatusFromString(evt.PreviousStatus),
			Status:           ValidatorNodeStatusFromString(evt.NextStatus),
			VotingPower:      evt.TmVotingPower,
			RankingScore:     rankingScore,
			TxHash:           txHash,
			ZetaTime:         zetaTime,
			EpochSeq:         epochSeq,
		}, RankingScoreAux{
			NodeID:   NodeID(evt.NodeId),
			EpochSeq: epochSeq,
		}, nil
}

func (rs *RankingScore) ToProto() *zeta.RankingScore {
	return &zeta.RankingScore{
		StakeScore:       rs.StakeScore.String(),
		PerformanceScore: rs.PerformanceScore.String(),
		PreviousStatus:   zeta.ValidatorNodeStatus(rs.PreviousStatus),
		Status:           zeta.ValidatorNodeStatus(rs.Status),
		VotingPower:      rs.VotingPower,
		RankingScore:     rs.RankingScore.String(),
	}
}

func RewardScoreFromScoreEvent(evt eventspb.ValidatorScoreEvent, txHash TxHash, zetaTime time.Time) (RewardScore, RewardScoreAux, error) {
	rawValidatorScore, err := decimal.NewFromString(evt.RawValidatorScore)
	if err != nil {
		return RewardScore{}, RewardScoreAux{}, err
	}

	performanceScore, err := decimal.NewFromString(evt.ValidatorPerformance)
	if err != nil {
		return RewardScore{}, RewardScoreAux{}, err
	}

	multisigScore, err := decimal.NewFromString(evt.MultisigScore)
	if err != nil {
		return RewardScore{}, RewardScoreAux{}, err
	}

	validatorScore, err := decimal.NewFromString(evt.ValidatorScore)
	if err != nil {
		return RewardScore{}, RewardScoreAux{}, err
	}

	normalisedScore, err := decimal.NewFromString(evt.NormalisedScore)
	if err != nil {
		return RewardScore{}, RewardScoreAux{}, err
	}

	epochSeq, err := strconv.ParseUint(evt.EpochSeq, 10, 64)
	if err != nil {
		return RewardScore{}, RewardScoreAux{}, err
	}

	return RewardScore{
			RawValidatorScore:   rawValidatorScore,
			PerformanceScore:    performanceScore,
			MultisigScore:       multisigScore,
			ValidatorScore:      validatorScore,
			NormalisedScore:     normalisedScore,
			ValidatorNodeStatus: ValidatorNodeStatusFromString(evt.ValidatorStatus),
			TxHash:              txHash,
			ZetaTime:            zetaTime,
			EpochSeq:            epochSeq,
		}, RewardScoreAux{
			NodeID:   NodeID(evt.NodeId),
			EpochSeq: epochSeq,
		}, nil
}

func (rs *RewardScore) ToProto() *zeta.RewardScore {
	return &zeta.RewardScore{
		RawValidatorScore: rs.RawValidatorScore.String(),
		PerformanceScore:  rs.PerformanceScore.String(),
		MultisigScore:     rs.MultisigScore.String(),
		ValidatorScore:    rs.ValidatorScore.String(),
		NormalisedScore:   rs.NormalisedScore.String(),
		ValidatorStatus:   zeta.ValidatorNodeStatus(rs.ValidatorNodeStatus),
	}
}

func NodeFromProto(node *zeta.Node, txHash TxHash, vegaTime time.Time) (Node, error) {
	stakedByOperator, err := decimal.NewFromString(node.StakedByOperator)
	if err != nil {
		return Node{}, err
	}

	stakedByDelegates, err := decimal.NewFromString(node.StakedByDelegates)
	if err != nil {
		return Node{}, err
	}

	stakedTotal, err := decimal.NewFromString(node.StakedTotal)
	if err != nil {
		return Node{}, err
	}

	maxIntendedStake, err := decimal.NewFromString(node.MaxIntendedStake)
	if err != nil {
		return Node{}, err
	}

	pendingStake, err := decimal.NewFromString(node.PendingStake)
	if err != nil {
		return Node{}, err
	}

	delegations := make([]Delegation, len(node.Delegations))
	for i, delegation := range node.Delegations {
		delegations[i], err = DelegationFromProto(delegation, txHash)
		if err != nil {
			return Node{}, err
		}
	}

	return Node{
		ID:                NodeID(node.Id),
		PubKey:            ZetaPublicKey(node.PubKey),
		TmPubKey:          TendermintPublicKey(node.TmPubKey),
		EthereumAddress:   EthereumAddress(node.EthereumAddress),
		InfoURL:           node.InfoUrl,
		Location:          node.Location,
		StakedByOperator:  stakedByOperator,
		StakedByDelegates: stakedByDelegates,
		StakedTotal:       stakedTotal,
		MaxIntendedStake:  maxIntendedStake,
		PendingStake:      pendingStake,
		EpochData:         EpochData{node.EpochData},
		Status:            NodeStatus(node.Status),
		Delegations:       delegations,
		// RewardScore:       RewardScore{node.RewardScore},
		// RankingScore:      RankingScore{node.RankingScore},
		Name:      node.Name,
		AvatarURL: node.AvatarUrl,
		TxHash:    txHash,
		ZetaTime:  zetaTime,
	}, nil
}

func (node *Node) ToProto() *zeta.Node {
	protoDelegations := make([]*zeta.Delegation, len(node.Delegations))
	for i, delegation := range node.Delegations {
		protoDelegations[i] = delegation.ToProto()
	}

	res := &zeta.Node{
		Id:                node.ID.String(),
		PubKey:            node.PubKey.String(),
		TmPubKey:          node.TmPubKey.String(),
		EthereumAddress:   node.EthereumAddress.String(),
		InfoUrl:           node.InfoURL,
		Location:          node.Location,
		StakedByOperator:  node.StakedByOperator.String(),
		StakedByDelegates: node.StakedByDelegates.String(),
		StakedTotal:       node.StakedTotal.String(),
		MaxIntendedStake:  node.MaxIntendedStake.String(),
		PendingStake:      node.PendingStake.String(),
		EpochData:         node.EpochData.EpochData,
		Status:            zeta.NodeStatus(node.Status),
		Delegations:       protoDelegations,
		Name:              node.Name,
		AvatarUrl:         node.AvatarURL,
	}

	if node.RewardScore != nil {
		res.RewardScore = node.RewardScore.ToProto()
	}

	if node.RankingScore != nil {
		res.RankingScore = node.RankingScore.ToProto()
	}

	return res
}

func (node Node) Cursor() *Cursor {
	return NewCursor(NodeCursor{ID: node.ID}.String())
}

func (node Node) ToProtoEdge(_ ...any) (*v2.NodeEdge, error) {
	return &v2.NodeEdge{
		Node:   node.ToProto(),
		Cursor: node.Cursor().Encode(),
	}, nil
}

func (ed EpochData) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(ed)
}

func (ed *EpochData) UnmarshalJSON(b []byte) error {
	ed.EpochData = &zeta.EpochData{}
	return protojson.Unmarshal(b, ed)
}

func (n *NodeSet) ToProto() *zeta.NodeSet {
	return &zeta.NodeSet{
		Total:    n.Total,
		Demoted:  n.Demoted,
		Promoted: n.Promoted,
		Inactive: n.Inactive,
	}
}

func (n *NodeData) ToProto() *zeta.NodeData {
	return &zeta.NodeData{
		StakedTotal:     n.StakedTotal.String(),
		TotalNodes:      n.TotalNodes,
		InactiveNodes:   n.InactiveNodes,
		Uptime:          float32(n.Uptime),
		TendermintNodes: n.TendermintNodes.ToProto(),
		ErsatzNodes:     n.ErsatzNodes.ToProto(),
		PendingNodes:    n.PendingNodes.ToProto(),
	}
}

type NodeCursor struct {
	ID NodeID `json:"id"`
}

func (nc NodeCursor) String() string {
	bs, err := json.Marshal(nc)
	if err != nil {
		panic(fmt.Errorf("could not marshal node cursor: %w", err))
	}
	return string(bs)
}

func (nc *NodeCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}
	return json.Unmarshal([]byte(cursorString), nc)
}
