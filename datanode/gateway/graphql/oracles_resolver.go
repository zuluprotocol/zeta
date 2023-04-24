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

	"zuluprotocol/zeta/zeta/libs/ptr"
	v2 "zuluprotocol/zeta/zeta/protos/data-node/api/v2"
	zetapb "code.zetaprotocol.io/zeta/protos/zeta"
	v1 "zuluprotocol/zeta/zeta/protos/zeta/data/v1"
)

type oracleSpecResolver ZetaResolverRoot

func (o *oracleSpecResolver) DataSourceSpec(_ context.Context, obj *zetapb.OracleSpec) (extDss *ExternalDataSourceSpec, _ error) {
	extDss = &ExternalDataSourceSpec{Spec: &DataSourceSpec{Data: &DataSourceDefinition{}}}
	if obj.ExternalDataSourceSpec != nil {
		extDss.Spec = resolveDataSourceSpec(obj.ExternalDataSourceSpec.Spec)
	}
	return
}

func (o *oracleSpecResolver) DataConnection(ctx context.Context, obj *zetapb.OracleSpec, pagination *v2.Pagination) (*v2.OracleDataConnection, error) {
	var specID *string
	if ed := obj.ExternalDataSourceSpec; ed != nil && ed.Spec != nil && ed.Spec.Id != "" {
		specID = &ed.Spec.Id
	}

	req := v2.ListOracleDataRequest{
		OracleSpecId: specID,
		Pagination:   pagination,
	}

	resp, err := o.tradingDataClientV2.ListOracleData(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.OracleData, nil
}

type oracleDataResolver ZetaResolverRoot

func (o *oracleDataResolver) ExternalData(_ context.Context, obj *zetapb.OracleData) (ed *ExternalData, _ error) {
	ed = &ExternalData{
		Data: &Data{},
	}

	oed := obj.ExternalData
	if oed == nil || oed.Data == nil {
		return
	}

	ed.Data.Signers = resolveSigners(oed.Data.Signers)
	ed.Data.Data = oed.Data.Data
	ed.Data.MatchedSpecIds = oed.Data.MatchedSpecIds
	ed.Data.BroadcastAt = oed.Data.BroadcastAt

	return
}

func resolveSigners(obj []*v1.Signer) (signers []*Signer) {
	for i := range obj {
		signers = append(signers, &Signer{Signer: resolveSigner(obj[i].Signer)})
	}
	return
}

func resolveSigner(obj any) (signer SignerKind) {
	switch sig := obj.(type) {
	case *v1.Signer_PubKey:
		signer = &PubKey{Key: &sig.PubKey.Key}
	case *v1.Signer_EthAddress:
		signer = &ETHAddress{Address: &sig.EthAddress.Address}
	}
	return
}

func resolveDataSourceDefinition(d *zetapb.DataSourceDefinition) (ds *DataSourceDefinition) {
	ds = &DataSourceDefinition{}
	if d == nil {
		return
	}
	switch dst := d.SourceType.(type) {
	case *zetapb.DataSourceDefinition_External:
		ds.SourceType = DataSourceDefinitionExternal{
			SourceType: dst.External.GetOracle(),
		}
	case *zetapb.DataSourceDefinition_Internal:
		ds.SourceType = DataSourceDefinitionInternal{
			SourceType: dst.Internal.GetTime(),
		}
	}
	return
}

func resolveDataSourceSpec(d *zetapb.DataSourceSpec) (ds *DataSourceSpec) {
	ds = &DataSourceSpec{
		Data: &DataSourceDefinition{},
	}
	if d == nil {
		return
	}

	ds.ID = d.GetId()
	ds.CreatedAt = d.CreatedAt
	if d.UpdatedAt != 0 {
		ds.UpdatedAt = ptr.From(d.UpdatedAt)
	}

	switch d.Status {
	case zetapb.DataSourceSpec_STATUS_ACTIVE:
		ds.Status = DataSourceSpecStatusStatusActive
	case zetapb.DataSourceSpec_STATUS_DEACTIVATED:
		ds.Status = DataSourceSpecStatusStatusDeactivated
	}

	if d.Data != nil {
		ds.Data = resolveDataSourceDefinition(d.Data)
	}

	return
}
