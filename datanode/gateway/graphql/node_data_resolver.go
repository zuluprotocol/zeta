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

package gql

import (
	"context"

	"zuluprotocol/zeta/libs/ptr"
	proto "zuluprotocol/zeta/protos/zeta"
)

type nodeDataResolver ZetaResolverRoot

func toNodeSet(obj *proto.NodeSet) *NodeSet {
	ns := &NodeSet{
		Total:    int(obj.Total),
		Demoted:  obj.Demoted,
		Promoted: obj.Promoted,
		Inactive: int(obj.Inactive),
	}
	if obj.Maximum != nil {
		ns.Maximum = ptr.From(int(*obj.Maximum))
	}
	return ns
}

func (r *nodeDataResolver) TotalNodes(ctx context.Context, obj *proto.NodeData) (int, error) {
	return int(obj.TotalNodes), nil
}

func (r *nodeDataResolver) InactiveNodes(ctx context.Context, obj *proto.NodeData) (int, error) {
	return int(obj.InactiveNodes), nil
}

func (r *nodeDataResolver) Uptime(ctx context.Context, obj *proto.NodeData) (float64, error) {
	return float64(obj.Uptime), nil
}

func (r *nodeDataResolver) TendermintNodes(ctx context.Context, obj *proto.NodeData) (*NodeSet, error) {
	return toNodeSet(obj.TendermintNodes), nil
}

func (r *nodeDataResolver) ErsatzNodes(ctx context.Context, obj *proto.NodeData) (*NodeSet, error) {
	if obj.ErsatzNodes == nil || obj.ErsatzNodes.Total == 0 {
		return nil, nil
	}
	return toNodeSet(obj.ErsatzNodes), nil
}

func (r *nodeDataResolver) PendingNodes(ctx context.Context, obj *proto.NodeData) (*NodeSet, error) {
	if obj.PendingNodes == nil || obj.PendingNodes.Total == 0 {
		return nil, nil
	}
	return toNodeSet(obj.PendingNodes), nil
}
