package entities

import "zuluprotocol/zeta/zeta/protos/zeta"

type OrderFilter struct {
	Statuses         []zeta.Order_Status
	Types            []zeta.Order_Type
	TimeInForces     []zeta.Order_TimeInForce
	ExcludeLiquidity bool
}
