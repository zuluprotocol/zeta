package types

import (
	"fmt"

	zetapb "code.zetaprotocol.io/zeta/protos/zeta"
)

// DataSourceSpecConfigurationTime is used internally.
type DataSourceSpecConfigurationTime struct {
	Conditions []*DataSourceSpecCondition
}

func (s DataSourceSpecConfigurationTime) isDataSourceType() {}

func (s DataSourceSpecConfigurationTime) oneOfProto() interface{} {
	return s
}

// /
// String returns the content of DataSourceSpecConfigurationTime as a string.
func (s DataSourceSpecConfigurationTime) String() string {
	return fmt.Sprintf(
		"conditions(%s)", DataSourceSpecConditions(s.Conditions).String(),
	)
}

func (s DataSourceSpecConfigurationTime) IntoProto() *zetapb.DataSourceSpecConfigurationTime {
	return &zetapb.DataSourceSpecConfigurationTime{
		Conditions: DataSourceSpecConditions(s.Conditions).IntoProto(),
	}
}

func (s DataSourceSpecConfigurationTime) DeepClone() dataSourceType {
	conditions := []*DataSourceSpecCondition{}
	conditions = append(conditions, s.Conditions...)

	return &DataSourceSpecConfigurationTime{
		Conditions: conditions,
	}
}

func DataSourceSpecConfigurationTimeFromProto(protoConfig *zetapb.DataSourceSpecConfigurationTime) *DataSourceSpecConfigurationTime {
	return &DataSourceSpecConfigurationTime{
		Conditions: DataSourceSpecConditionsFromProto(protoConfig.Conditions),
	}
}

type DataSourceDefinitionInternalTime struct {
	Time *DataSourceSpecConfigurationTime
}

func (i *DataSourceDefinitionInternalTime) isDataSourceType() {}

func (i *DataSourceDefinitionInternalTime) oneOfProto() interface{} {
	return i.IntoProto()
}

func (i *DataSourceDefinitionInternalTime) IntoProto() *zetapb.DataSourceDefinitionInternal_Time {
	ids := &zetapb.DataSourceSpecConfigurationTime{}
	if i.Time != nil {
		ids = i.Time.IntoProto()
	}

	return &zetapb.DataSourceDefinitionInternal_Time{
		Time: ids,
	}
}

func (i *DataSourceDefinitionInternalTime) DeepClone() dataSourceType {
	if i.Time == nil {
		return &DataSourceDefinitionInternalTime{
			Time: &DataSourceSpecConfigurationTime{},
		}
	}

	return nil
}

func (i *DataSourceDefinitionInternalTime) String() string {
	if i.Time == nil {
		return ""
	}
	return i.Time.String()
}

func DataSourceDefinitionInternalTimeFromProto(protoConfig *zetapb.DataSourceDefinitionInternal_Time) *DataSourceDefinitionInternalTime {
	ids := &DataSourceDefinitionInternalTime{
		Time: &DataSourceSpecConfigurationTime{},
	}

	if protoConfig != nil {
		if protoConfig.Time != nil {
			ids.Time = DataSourceSpecConfigurationTimeFromProto(protoConfig.Time)
		}
	}

	return ids
}
