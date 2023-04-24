package entities

import "code.zetaprotocol.io/vega/protos/vega"

type OrderFilter struct {
	Statuses         []zeta.Order_Status
	Types            []zeta.Order_Type
	TimeInForces     []zeta.Order_TimeInForce
	ExcludeLiquidity bool
}
