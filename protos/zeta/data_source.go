package zeta

import (
	datapb "zuluprotocol/zeta/protos/zeta/data/v1"
)

func (o *DataSourceSpecConfiguration) ToOracleSpec(d *DataSourceSpec) *OracleSpec {
	return NewOracleSpec(d)
}

func (DataSourceSpec) IsEvent() {}

func NewDataSourceSpec(sc *DataSourceDefinition) *DataSourceSpec {
	ds := &DataSourceSpec{}
	tp := sc.GetSourceType()
	if tp != nil {
		switch sc.SourceType.(type) {
		case *DataSourceDefinition_External:
			ext := sc.GetExternal()
			if ext != nil {
				o := ext.GetOracle()
				if o != nil {
					ds.Id = datapb.NewID(o.Signers, o.Filters)
				}
			}
		case *DataSourceDefinition_Internal:
			in := sc.GetInternal()
			if in != nil {
				t := in.GetTime()
				if t != nil {
					//
				}
			}
		}
	}

	ds.Data = sc
	return ds
}
