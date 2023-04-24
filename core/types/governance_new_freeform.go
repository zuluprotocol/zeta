// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.ZETA file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package types

import (
	"fmt"

	zetapb "code.zetaprotocol.io/zeta/protos/zeta"
)

type ProposalTermsNewFreeform struct {
	NewFreeform *NewFreeform
}

func (f ProposalTermsNewFreeform) String() string {
	return fmt.Sprintf(
		"newFreeForm(%s)",
		reflectPointerToString(f.NewFreeform),
	)
}

func (f ProposalTermsNewFreeform) IntoProto() *zetapb.ProposalTerms_NewFreeform {
	var newFreeform *zetapb.NewFreeform
	if f.NewFreeform != nil {
		newFreeform = f.NewFreeform.IntoProto()
	}
	return &zetapb.ProposalTerms_NewFreeform{
		NewFreeform: newFreeform,
	}
}

func (f ProposalTermsNewFreeform) isPTerm() {}

func (f ProposalTermsNewFreeform) oneOfProto() interface{} {
	return f.IntoProto()
}

func (f ProposalTermsNewFreeform) GetTermType() ProposalTermsType {
	return ProposalTermsTypeNewFreeform
}

func (f ProposalTermsNewFreeform) DeepClone() proposalTerm {
	if f.NewFreeform == nil {
		return &ProposalTermsNewFreeform{}
	}
	return &ProposalTermsNewFreeform{
		NewFreeform: f.NewFreeform.DeepClone(),
	}
}

func NewNewFreeformFromProto(_ *zetapb.ProposalTerms_NewFreeform) *ProposalTermsNewFreeform {
	return &ProposalTermsNewFreeform{
		NewFreeform: &NewFreeform{},
	}
}

type NewFreeform struct{}

func (n NewFreeform) IntoProto() *zetapb.NewFreeform {
	return &zetapb.NewFreeform{}
}

func (n NewFreeform) String() string {
	return ""
}

func (n NewFreeform) DeepClone() *NewFreeform {
	return &NewFreeform{}
}
