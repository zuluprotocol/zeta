package gql

import (
	"context"

	v2 "code.zetaprotocol.io/vega/protos/data-node/api/v2"
	"code.zetaprotocol.io/vega/protos/vega"
)

type orderFilterResolver ZetaResolverRoot

func (o orderFilterResolver) Status(ctx context.Context, obj *v2.OrderFilter, data []zeta.Order_Status) error {
	obj.Statuses = data
	return nil
}

func (o orderFilterResolver) TimeInForce(ctx context.Context, obj *v2.OrderFilter, data []zeta.Order_TimeInForce) error {
	obj.TimeInForces = data
	return nil
}