package service

//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks zuluprotocol/zeta/datanode/service OrderStore,ChainStore,MarketStore,MarketDataStore,PositionStore
