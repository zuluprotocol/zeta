// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: zeta/assets.proto

package zeta

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Asset_Status int32

const (
	// Default value, always invalid
	Asset_STATUS_UNSPECIFIED Asset_Status = 0
	// Asset is proposed and under vote
	Asset_STATUS_PROPOSED Asset_Status = 1
	// Asset has been rejected from governance
	Asset_STATUS_REJECTED Asset_Status = 2
	// Asset is pending listing from the bridge
	Asset_STATUS_PENDING_LISTING Asset_Status = 3
	// Asset is fully usable in the network
	Asset_STATUS_ENABLED Asset_Status = 4
)

// Enum value maps for Asset_Status.
var (
	Asset_Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "STATUS_PROPOSED",
		2: "STATUS_REJECTED",
		3: "STATUS_PENDING_LISTING",
		4: "STATUS_ENABLED",
	}
	Asset_Status_value = map[string]int32{
		"STATUS_UNSPECIFIED":     0,
		"STATUS_PROPOSED":        1,
		"STATUS_REJECTED":        2,
		"STATUS_PENDING_LISTING": 3,
		"STATUS_ENABLED":         4,
	}
)

func (x Asset_Status) Enum() *Asset_Status {
	p := new(Asset_Status)
	*p = x
	return p
}

func (x Asset_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Asset_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_zeta_assets_proto_enumTypes[0].Descriptor()
}

func (Asset_Status) Type() protoreflect.EnumType {
	return &file_zeta_assets_proto_enumTypes[0]
}

func (x Asset_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Asset_Status.Descriptor instead.
func (Asset_Status) EnumDescriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{0, 0}
}

// The Zeta representation of an external asset
type Asset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Internal identifier of the asset
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The definition of the external source for this asset
	Details *AssetDetails `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
	// Status of the asset
	Status Asset_Status `protobuf:"varint,3,opt,name=status,proto3,enum=zeta.Asset_Status" json:"status,omitempty"`
}

func (x *Asset) Reset() {
	*x = Asset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zeta_assets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Asset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Asset) ProtoMessage() {}

func (x *Asset) ProtoReflect() protoreflect.Message {
	mi := &file_zeta_assets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Asset.ProtoReflect.Descriptor instead.
func (*Asset) Descriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{0}
}

func (x *Asset) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Asset) GetDetails() *AssetDetails {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *Asset) GetStatus() Asset_Status {
	if x != nil {
		return x.Status
	}
	return Asset_STATUS_UNSPECIFIED
}

// The Zeta representation of an external asset
type AssetDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the asset (e.g: Great British Pound)
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Symbol of the asset (e.g: GBP)
	Symbol string `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Number of decimal / precision handled by this asset
	Decimals uint64 `protobuf:"varint,4,opt,name=decimals,proto3" json:"decimals,omitempty"`
	// The minimum economically meaningful amount in the asset
	Quantum string `protobuf:"bytes,5,opt,name=quantum,proto3" json:"quantum,omitempty"`
	// The source
	//
	// Types that are assignable to Source:
	//
	//	*AssetDetails_BuiltinAsset
	//	*AssetDetails_Erc20
	Source isAssetDetails_Source `protobuf_oneof:"source"`
}

func (x *AssetDetails) Reset() {
	*x = AssetDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zeta_assets_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetDetails) ProtoMessage() {}

func (x *AssetDetails) ProtoReflect() protoreflect.Message {
	mi := &file_zeta_assets_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetDetails.ProtoReflect.Descriptor instead.
func (*AssetDetails) Descriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{1}
}

func (x *AssetDetails) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AssetDetails) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *AssetDetails) GetDecimals() uint64 {
	if x != nil {
		return x.Decimals
	}
	return 0
}

func (x *AssetDetails) GetQuantum() string {
	if x != nil {
		return x.Quantum
	}
	return ""
}

func (m *AssetDetails) GetSource() isAssetDetails_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (x *AssetDetails) GetBuiltinAsset() *BuiltinAsset {
	if x, ok := x.GetSource().(*AssetDetails_BuiltinAsset); ok {
		return x.BuiltinAsset
	}
	return nil
}

func (x *AssetDetails) GetErc20() *ERC20 {
	if x, ok := x.GetSource().(*AssetDetails_Erc20); ok {
		return x.Erc20
	}
	return nil
}

type isAssetDetails_Source interface {
	isAssetDetails_Source()
}

type AssetDetails_BuiltinAsset struct {
	// A built-in asset
	BuiltinAsset *BuiltinAsset `protobuf:"bytes,101,opt,name=builtin_asset,json=builtinAsset,proto3,oneof"`
}

type AssetDetails_Erc20 struct {
	// An Ethereum ERC20 asset
	Erc20 *ERC20 `protobuf:"bytes,102,opt,name=erc20,proto3,oneof"`
}

func (*AssetDetails_BuiltinAsset) isAssetDetails_Source() {}

func (*AssetDetails_Erc20) isAssetDetails_Source() {}

// A Zeta internal asset
type BuiltinAsset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Maximum amount that can be requested by a party through the built-in asset faucet at a time
	MaxFaucetAmountMint string `protobuf:"bytes,1,opt,name=max_faucet_amount_mint,json=maxFaucetAmountMint,proto3" json:"max_faucet_amount_mint,omitempty"`
}

func (x *BuiltinAsset) Reset() {
	*x = BuiltinAsset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zeta_assets_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuiltinAsset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuiltinAsset) ProtoMessage() {}

func (x *BuiltinAsset) ProtoReflect() protoreflect.Message {
	mi := &file_zeta_assets_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuiltinAsset.ProtoReflect.Descriptor instead.
func (*BuiltinAsset) Descriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{2}
}

func (x *BuiltinAsset) GetMaxFaucetAmountMint() string {
	if x != nil {
		return x.MaxFaucetAmountMint
	}
	return ""
}

// An ERC20 token based asset, living on the ethereum network
type ERC20 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The address of the contract for the token, on the ethereum network
	ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	// The lifetime limits deposit per address
	// note: this is a temporary measure that can be changed by governance
	LifetimeLimit string `protobuf:"bytes,2,opt,name=lifetime_limit,json=lifetimeLimit,proto3" json:"lifetime_limit,omitempty"`
	// The maximum you can withdraw instantly. All withdrawals over the threshold will be delayed by the withdrawal delay.
	// There’s no limit on the size of a withdrawal
	// note: this is a temporary measure that can be changed by governance
	WithdrawThreshold string `protobuf:"bytes,3,opt,name=withdraw_threshold,json=withdrawThreshold,proto3" json:"withdraw_threshold,omitempty"`
}

func (x *ERC20) Reset() {
	*x = ERC20{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zeta_assets_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ERC20) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ERC20) ProtoMessage() {}

func (x *ERC20) ProtoReflect() protoreflect.Message {
	mi := &file_zeta_assets_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ERC20.ProtoReflect.Descriptor instead.
func (*ERC20) Descriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{3}
}

func (x *ERC20) GetContractAddress() string {
	if x != nil {
		return x.ContractAddress
	}
	return ""
}

func (x *ERC20) GetLifetimeLimit() string {
	if x != nil {
		return x.LifetimeLimit
	}
	return ""
}

func (x *ERC20) GetWithdrawThreshold() string {
	if x != nil {
		return x.WithdrawThreshold
	}
	return ""
}

// The changes to apply on an existing asset.
type AssetDetailsUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The minimum economically meaningful amount in the asset
	Quantum string `protobuf:"bytes,5,opt,name=quantum,proto3" json:"quantum,omitempty"`
	// The source
	//
	// Types that are assignable to Source:
	//
	//	*AssetDetailsUpdate_Erc20
	Source isAssetDetailsUpdate_Source `protobuf_oneof:"source"`
}

func (x *AssetDetailsUpdate) Reset() {
	*x = AssetDetailsUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zeta_assets_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetDetailsUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetDetailsUpdate) ProtoMessage() {}

func (x *AssetDetailsUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_zeta_assets_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetDetailsUpdate.ProtoReflect.Descriptor instead.
func (*AssetDetailsUpdate) Descriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{4}
}

func (x *AssetDetailsUpdate) GetQuantum() string {
	if x != nil {
		return x.Quantum
	}
	return ""
}

func (m *AssetDetailsUpdate) GetSource() isAssetDetailsUpdate_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (x *AssetDetailsUpdate) GetErc20() *ERC20Update {
	if x, ok := x.GetSource().(*AssetDetailsUpdate_Erc20); ok {
		return x.Erc20
	}
	return nil
}

type isAssetDetailsUpdate_Source interface {
	isAssetDetailsUpdate_Source()
}

type AssetDetailsUpdate_Erc20 struct {
	// An Ethereum ERC20 asset
	Erc20 *ERC20Update `protobuf:"bytes,101,opt,name=erc20,proto3,oneof"`
}

func (*AssetDetailsUpdate_Erc20) isAssetDetailsUpdate_Source() {}

type ERC20Update struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The lifetime limits deposit per address.
	// This will be interpreted against the asset decimals.
	// note: this is a temporary measure that can be changed by governance
	LifetimeLimit string `protobuf:"bytes,1,opt,name=lifetime_limit,json=lifetimeLimit,proto3" json:"lifetime_limit,omitempty"`
	// The maximum you can withdraw instantly. All withdrawals over the threshold will be delayed by the withdrawal delay.
	// There’s no limit on the size of a withdrawal
	// note: this is a temporary measure that can be changed by governance
	WithdrawThreshold string `protobuf:"bytes,2,opt,name=withdraw_threshold,json=withdrawThreshold,proto3" json:"withdraw_threshold,omitempty"`
}

func (x *ERC20Update) Reset() {
	*x = ERC20Update{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zeta_assets_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ERC20Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ERC20Update) ProtoMessage() {}

func (x *ERC20Update) ProtoReflect() protoreflect.Message {
	mi := &file_zeta_assets_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ERC20Update.ProtoReflect.Descriptor instead.
func (*ERC20Update) Descriptor() ([]byte, []int) {
	return file_zeta_assets_proto_rawDescGZIP(), []int{5}
}

func (x *ERC20Update) GetLifetimeLimit() string {
	if x != nil {
		return x.LifetimeLimit
	}
	return ""
}

func (x *ERC20Update) GetWithdrawThreshold() string {
	if x != nil {
		return x.WithdrawThreshold
	}
	return ""
}

var File_zeta_assets_proto protoreflect.FileDescriptor

var file_zeta_assets_proto_rawDesc = []byte{
	0x0a, 0x11, 0x76, 0x65, 0x67, 0x61, 0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x76, 0x65, 0x67, 0x61, 0x22, 0xed, 0x01, 0x0a, 0x05, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x76, 0x65, 0x67, 0x61, 0x2e, 0x41, 0x73, 0x73, 0x65,
	0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x12, 0x2a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x76, 0x65, 0x67, 0x61, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x7a, 0x0a,
	0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x13, 0x0a, 0x0f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x50, 0x4f, 0x53,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x52,
	0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x5f, 0x4c, 0x49, 0x53, 0x54,
	0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f,
	0x45, 0x4e, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x22, 0xe0, 0x01, 0x0a, 0x0c, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61,
	0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61,
	0x6c, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x12, 0x39, 0x0a, 0x0d,
	0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x65, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x76, 0x65, 0x67, 0x61, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x74,
	0x69, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x48, 0x00, 0x52, 0x0c, 0x62, 0x75, 0x69, 0x6c, 0x74,
	0x69, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x23, 0x0a, 0x05, 0x65, 0x72, 0x63, 0x32, 0x30,
	0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x76, 0x65, 0x67, 0x61, 0x2e, 0x45, 0x52,
	0x43, 0x32, 0x30, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x63, 0x32, 0x30, 0x42, 0x08, 0x0a, 0x06,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0x43, 0x0a, 0x0c,
	0x42, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x33, 0x0a, 0x16,
	0x6d, 0x61, 0x78, 0x5f, 0x66, 0x61, 0x75, 0x63, 0x65, 0x74, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x6d, 0x61,
	0x78, 0x46, 0x61, 0x75, 0x63, 0x65, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x4d, 0x69, 0x6e,
	0x74, 0x22, 0x88, 0x01, 0x0a, 0x05, 0x45, 0x52, 0x43, 0x32, 0x30, 0x12, 0x29, 0x0a, 0x10, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69,
	0x6d, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x2d, 0x0a,
	0x12, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68,
	0x6f, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x77, 0x69, 0x74, 0x68, 0x64,
	0x72, 0x61, 0x77, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x22, 0x7b, 0x0a, 0x12,
	0x41, 0x73, 0x73, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x12, 0x29, 0x0a, 0x05,
	0x65, 0x72, 0x63, 0x32, 0x30, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x65,
	0x67, 0x61, 0x2e, 0x45, 0x52, 0x43, 0x32, 0x30, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x00,
	0x52, 0x05, 0x65, 0x72, 0x63, 0x32, 0x30, 0x42, 0x08, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08,
	0x03, 0x10, 0x04, 0x4a, 0x04, 0x08, 0x04, 0x10, 0x05, 0x22, 0x63, 0x0a, 0x0b, 0x45, 0x52, 0x43,
	0x32, 0x30, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x69, 0x66, 0x65,
	0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12,
	0x2d, 0x0a, 0x12, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x5f, 0x74, 0x68, 0x72, 0x65,
	0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x77, 0x69, 0x74,
	0x68, 0x64, 0x72, 0x61, 0x77, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x42, 0x27,
	0x5a, 0x25, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x76, 0x65, 0x67, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x65, 0x67, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2f, 0x76, 0x65, 0x67, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_zeta_assets_proto_rawDescOnce sync.Once
	file_zeta_assets_proto_rawDescData = file_vega_assets_proto_rawDesc
)

func file_zeta_assets_proto_rawDescGZIP() []byte {
	file_zeta_assets_proto_rawDescOnce.Do(func() {
		file_zeta_assets_proto_rawDescData = protoimpl.X.CompressGZIP(file_vega_assets_proto_rawDescData)
	})
	return file_zeta_assets_proto_rawDescData
}

var file_zeta_assets_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_zeta_assets_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_zeta_assets_proto_goTypes = []interface{}{
	(Asset_Status)(0),          // 0: zeta.Asset.Status
	(*Asset)(nil),              // 1: zeta.Asset
	(*AssetDetails)(nil),       // 2: zeta.AssetDetails
	(*BuiltinAsset)(nil),       // 3: zeta.BuiltinAsset
	(*ERC20)(nil),              // 4: zeta.ERC20
	(*AssetDetailsUpdate)(nil), // 5: zeta.AssetDetailsUpdate
	(*ERC20Update)(nil),        // 6: zeta.ERC20Update
}
var file_zeta_assets_proto_depIdxs = []int32{
	2, // 0: zeta.Asset.details:type_name -> vega.AssetDetails
	0, // 1: zeta.Asset.status:type_name -> vega.Asset.Status
	3, // 2: zeta.AssetDetails.builtin_asset:type_name -> vega.BuiltinAsset
	4, // 3: zeta.AssetDetails.erc20:type_name -> vega.ERC20
	6, // 4: zeta.AssetDetailsUpdate.erc20:type_name -> vega.ERC20Update
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_zeta_assets_proto_init() }
func file_zeta_assets_proto_init() {
	if File_zeta_assets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_zeta_assets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Asset); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zeta_assets_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetDetails); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zeta_assets_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuiltinAsset); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zeta_assets_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ERC20); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zeta_assets_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetDetailsUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zeta_assets_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ERC20Update); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_zeta_assets_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*AssetDetails_BuiltinAsset)(nil),
		(*AssetDetails_Erc20)(nil),
	}
	file_zeta_assets_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*AssetDetailsUpdate_Erc20)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_zeta_assets_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_zeta_assets_proto_goTypes,
		DependencyIndexes: file_zeta_assets_proto_depIdxs,
		EnumInfos:         file_zeta_assets_proto_enumTypes,
		MessageInfos:      file_zeta_assets_proto_msgTypes,
	}.Build()
	File_zeta_assets_proto = out.File
	file_zeta_assets_proto_rawDesc = nil
	file_zeta_assets_proto_goTypes = nil
	file_zeta_assets_proto_depIdxs = nil
}