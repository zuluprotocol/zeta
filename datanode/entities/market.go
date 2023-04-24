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
	"errors"
	"fmt"
	"math"
	"time"

	"code.zetaprotocol.io/vega/libs/num"
	v2 "code.zetaprotocol.io/vega/protos/data-node/api/v2"
	"code.zetaprotocol.io/vega/protos/vega"
	"github.com/shopspring/decimal"
)

type _Market struct{}

type MarketID = ID[_Market]

type Market struct {
	ID                            MarketID
	TxHash                        TxHash
	ZetaTime                      time.Time
	InstrumentID                  string
	TradableInstrument            TradableInstrument
	DecimalPlaces                 int
	Fees                          Fees
	OpeningAuction                AuctionDuration
	PriceMonitoringSettings       PriceMonitoringSettings
	LiquidityMonitoringParameters LiquidityMonitoringParameters
	TradingMode                   MarketTradingMode
	State                         MarketState
	MarketTimestamps              MarketTimestamps
	PositionDecimalPlaces         int
	LpPriceRange                  string
}

type MarketCursor struct {
	ZetaTime time.Time `json:"zetaTime"`
	ID       MarketID  `json:"id"`
}

func (mc MarketCursor) String() string {
	bs, err := json.Marshal(mc)
	if err != nil {
		panic(fmt.Errorf("could not marshal market cursor: %w", err))
	}
	return string(bs)
}

func (mc *MarketCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}

	return json.Unmarshal([]byte(cursorString), mc)
}

func NewMarketFromProto(market *zeta.Market, txHash TxHash, vegaTime time.Time) (*Market, error) {
	var err error
	var liquidityMonitoringParameters LiquidityMonitoringParameters
	var marketTimestamps MarketTimestamps
	var priceMonitoringSettings PriceMonitoringSettings
	var openingAuction AuctionDuration
	var fees Fees

	if fees, err = feesFromProto(market.Fees); err != nil {
		return nil, err
	}

	if market.OpeningAuction != nil {
		openingAuction.Duration = market.OpeningAuction.Duration
		openingAuction.Volume = market.OpeningAuction.Volume
	}

	if priceMonitoringSettings, err = priceMonitoringSettingsFromProto(market.PriceMonitoringSettings); err != nil {
		return nil, err
	}

	if liquidityMonitoringParameters, err = liquidityMonitoringParametersFromProto(market.LiquidityMonitoringParameters); err != nil {
		return nil, err
	}

	if marketTimestamps, err = marketTimestampsFromProto(market.MarketTimestamps); err != nil {
		return nil, err
	}

	if market.DecimalPlaces > math.MaxInt {
		return nil, fmt.Errorf("%d is not a valid number for decimal places", market.DecimalPlaces)
	}

	if market.PositionDecimalPlaces > math.MaxInt {
		return nil, fmt.Errorf("%d is not a valid number for position decimal places", market.PositionDecimalPlaces)
	}

	lppr, err := num.DecimalFromString(market.LpPriceRange)
	if err != nil || lppr.IsNegative() || lppr.IsZero() || lppr.GreaterThan(num.DecimalFromInt64(100)) {
		return nil, fmt.Errorf("%v is not a valid number for LP price range", market.LpPriceRange)
	}

	dps := int(market.DecimalPlaces)
	positionDps := int(market.PositionDecimalPlaces)

	return &Market{
		ID:                            MarketID(market.Id),
		TxHash:                        txHash,
		ZetaTime:                      zetaTime,
		InstrumentID:                  market.TradableInstrument.Instrument.Id,
		TradableInstrument:            TradableInstrument{market.TradableInstrument},
		DecimalPlaces:                 dps,
		Fees:                          fees,
		OpeningAuction:                openingAuction,
		PriceMonitoringSettings:       priceMonitoringSettings,
		LiquidityMonitoringParameters: liquidityMonitoringParameters,
		TradingMode:                   MarketTradingMode(market.TradingMode),
		State:                         MarketState(market.State),
		MarketTimestamps:              marketTimestamps,
		PositionDecimalPlaces:         positionDps,
		LpPriceRange:                  market.LpPriceRange,
	}, nil
}

func (m Market) ToProto() *zeta.Market {
	return &zeta.Market{
		Id:                 m.ID.String(),
		TradableInstrument: m.TradableInstrument.ToProto(),
		DecimalPlaces:      uint64(m.DecimalPlaces),
		Fees:               m.Fees.ToProto(),
		OpeningAuction: &zeta.AuctionDuration{
			Duration: m.OpeningAuction.Duration,
			Volume:   m.OpeningAuction.Volume,
		},
		PriceMonitoringSettings:       m.PriceMonitoringSettings.ToProto(),
		LiquidityMonitoringParameters: m.LiquidityMonitoringParameters.ToProto(),
		TradingMode:                   zeta.Market_TradingMode(m.TradingMode),
		State:                         zeta.Market_State(m.State),
		MarketTimestamps:              m.MarketTimestamps.ToProto(),
		PositionDecimalPlaces:         int64(m.PositionDecimalPlaces),
		LpPriceRange:                  m.LpPriceRange,
	}
}

func (m Market) Cursor() *Cursor {
	mc := MarketCursor{
		ZetaTime: m.VegaTime,
		ID:       m.ID,
	}
	return NewCursor(mc.String())
}

func (m Market) ToProtoEdge(_ ...any) (*v2.MarketEdge, error) {
	return &v2.MarketEdge{
		Node:   m.ToProto(),
		Cursor: m.Cursor().Encode(),
	}, nil
}

type MarketTimestamps struct {
	Proposed int64 `json:"proposed,omitempty"`
	Pending  int64 `json:"pending,omitempty"`
	Open     int64 `json:"open,omitempty"`
	Close    int64 `json:"close,omitempty"`
}

func (mt MarketTimestamps) ToProto() *zeta.MarketTimestamps {
	return &zeta.MarketTimestamps{
		Proposed: mt.Proposed,
		Pending:  mt.Pending,
		Open:     mt.Open,
		Close:    mt.Close,
	}
}

func marketTimestampsFromProto(ts *zeta.MarketTimestamps) (MarketTimestamps, error) {
	if ts == nil {
		return MarketTimestamps{}, errors.New("market timestamps cannot be nil")
	}

	return MarketTimestamps{
		Proposed: ts.Proposed,
		Pending:  ts.Pending,
		Open:     ts.Open,
		Close:    ts.Close,
	}, nil
}

type TargetStakeParameters struct {
	TimeWindow     int64   `json:"timeWindow,omitempty"`
	ScalingFactors float64 `json:"scalingFactor,omitempty"`
}

func (tsp TargetStakeParameters) ToProto() *zeta.TargetStakeParameters {
	return &zeta.TargetStakeParameters{
		TimeWindow:    tsp.TimeWindow,
		ScalingFactor: tsp.ScalingFactors,
	}
}

type LiquidityMonitoringParameters struct {
	TargetStakeParameters *TargetStakeParameters `json:"targetStakeParameters,omitempty"`
	TriggeringRatio       string                 `json:"triggeringRatio,omitempty"`
	AuctionExtension      int64                  `json:"auctionExtension,omitempty"`
}

func (lmp LiquidityMonitoringParameters) ToProto() *zeta.LiquidityMonitoringParameters {
	if lmp.TargetStakeParameters == nil {
		return nil
	}
	return &zeta.LiquidityMonitoringParameters{
		TargetStakeParameters: lmp.TargetStakeParameters.ToProto(),
		TriggeringRatio:       lmp.TriggeringRatio,
		AuctionExtension:      lmp.AuctionExtension,
	}
}

func liquidityMonitoringParametersFromProto(lmp *zeta.LiquidityMonitoringParameters) (LiquidityMonitoringParameters, error) {
	if lmp == nil {
		return LiquidityMonitoringParameters{}, errors.New("liquidity monitoring parameters cannot be Nil")
	}

	var tsp *TargetStakeParameters

	if lmp.TargetStakeParameters != nil {
		tsp = &TargetStakeParameters{
			TimeWindow:     lmp.TargetStakeParameters.TimeWindow,
			ScalingFactors: lmp.TargetStakeParameters.ScalingFactor,
		}
	}

	return LiquidityMonitoringParameters{
		TargetStakeParameters: tsp,
		TriggeringRatio:       lmp.TriggeringRatio,
		AuctionExtension:      lmp.AuctionExtension,
	}, nil
}

type PriceMonitoringParameters struct {
	Triggers []*PriceMonitoringTrigger `json:"triggers,omitempty"`
}

func priceMonitoringParametersFromProto(pmp *zeta.PriceMonitoringParameters) PriceMonitoringParameters {
	if len(pmp.Triggers) == 0 {
		return PriceMonitoringParameters{}
	}

	triggers := make([]*PriceMonitoringTrigger, 0, len(pmp.Triggers))

	for _, trigger := range pmp.Triggers {
		probability, _ := decimal.NewFromString(trigger.Probability)
		triggers = append(triggers, &PriceMonitoringTrigger{
			Horizon:          uint64(trigger.Horizon),
			Probability:      probability,
			AuctionExtension: uint64(trigger.AuctionExtension),
		})
	}

	return PriceMonitoringParameters{
		Triggers: triggers,
	}
}

type PriceMonitoringSettings struct {
	Parameters *PriceMonitoringParameters `json:"priceMonitoringParameters,omitempty"`
}

func (s PriceMonitoringSettings) ToProto() *zeta.PriceMonitoringSettings {
	if s.Parameters == nil {
		return nil
	}
	triggers := make([]*zeta.PriceMonitoringTrigger, 0, len(s.Parameters.Triggers))

	if len(s.Parameters.Triggers) > 0 {
		for _, trigger := range s.Parameters.Triggers {
			triggers = append(triggers, trigger.ToProto())
		}
	}

	return &zeta.PriceMonitoringSettings{
		Parameters: &zeta.PriceMonitoringParameters{
			Triggers: triggers,
		},
	}
}

func priceMonitoringSettingsFromProto(pms *zeta.PriceMonitoringSettings) (PriceMonitoringSettings, error) {
	if pms == nil {
		return PriceMonitoringSettings{}, errors.New("price monitoring settings cannot be nil")
	}

	parameters := priceMonitoringParametersFromProto(pms.Parameters)
	return PriceMonitoringSettings{
		Parameters: &parameters,
	}, nil
}

type AuctionDuration struct {
	Duration int64  `json:"duration,omitempty"`
	Volume   uint64 `json:"volume,omitempty"`
}

type FeeFactors struct {
	MakerFee          string `json:"makerFee,omitempty"`
	InfrastructureFee string `json:"infrastructureFee,omitempty"`
	LiquidityFee      string `json:"liquidityFee,omitempty"`
}

type Fees struct {
	Factors *FeeFactors `json:"factors,omitempty"`
}

func (f Fees) ToProto() *zeta.Fees {
	if f.Factors == nil {
		return nil
	}

	return &zeta.Fees{
		Factors: &zeta.FeeFactors{
			MakerFee:          f.Factors.MakerFee,
			InfrastructureFee: f.Factors.InfrastructureFee,
			LiquidityFee:      f.Factors.LiquidityFee,
		},
	}
}

func feesFromProto(fees *zeta.Fees) (Fees, error) {
	if fees == nil {
		return Fees{}, errors.New("fees cannot be Nil")
	}

	return Fees{
		Factors: &FeeFactors{
			MakerFee:          fees.Factors.MakerFee,
			InfrastructureFee: fees.Factors.InfrastructureFee,
			LiquidityFee:      fees.Factors.LiquidityFee,
		},
	}, nil
}