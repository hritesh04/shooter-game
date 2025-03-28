// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: movementEmitter.proto

package proto

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Action int32

const (
	Action_Movement Action = 0
	Action_Fire     Action = 1
	Action_Join     Action = 2
	Action_Info     Action = 3
)

// Enum value maps for Action.
var (
	Action_name = map[int32]string{
		0: "Movement",
		1: "Fire",
		2: "Join",
		3: "Info",
	}
	Action_value = map[string]int32{
		"Movement": 0,
		"Fire":     1,
		"Join":     2,
		"Info":     3,
	}
)

func (x Action) Enum() *Action {
	p := new(Action)
	*p = x
	return p
}

func (x Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Action) Descriptor() protoreflect.EnumDescriptor {
	return file_movementEmitter_proto_enumTypes[0].Descriptor()
}

func (Action) Type() protoreflect.EnumType {
	return &file_movementEmitter_proto_enumTypes[0]
}

func (x Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Action.Descriptor instead.
func (Action) EnumDescriptor() ([]byte, []int) {
	return file_movementEmitter_proto_rawDescGZIP(), []int{0}
}

type Direction int32

const (
	Direction_UP    Direction = 0
	Direction_DOWN  Direction = 1
	Direction_LEFT  Direction = 2
	Direction_RIGHT Direction = 3
	Direction_SHOOT Direction = 4
)

// Enum value maps for Direction.
var (
	Direction_name = map[int32]string{
		0: "UP",
		1: "DOWN",
		2: "LEFT",
		3: "RIGHT",
		4: "SHOOT",
	}
	Direction_value = map[string]int32{
		"UP":    0,
		"DOWN":  1,
		"LEFT":  2,
		"RIGHT": 3,
		"SHOOT": 4,
	}
)

func (x Direction) Enum() *Direction {
	p := new(Direction)
	*p = x
	return p
}

func (x Direction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Direction) Descriptor() protoreflect.EnumDescriptor {
	return file_movementEmitter_proto_enumTypes[1].Descriptor()
}

func (Direction) Type() protoreflect.EnumType {
	return &file_movementEmitter_proto_enumTypes[1]
}

func (x Direction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Direction.Descriptor instead.
func (Direction) EnumDescriptor() ([]byte, []int) {
	return file_movementEmitter_proto_rawDescGZIP(), []int{1}
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   Action    `protobuf:"varint,1,opt,name=type,proto3,enum=serverrpc.Action" json:"type,omitempty"`
	Data   Direction `protobuf:"varint,2,opt,name=data,proto3,enum=serverrpc.Direction" json:"data,omitempty"`
	RoomID string    `protobuf:"bytes,3,opt,name=roomID,proto3" json:"roomID,omitempty"`
	Name   string    `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Player []*Player `protobuf:"bytes,5,rep,name=player,proto3" json:"player,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movementEmitter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_movementEmitter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_movementEmitter_proto_rawDescGZIP(), []int{0}
}

func (x *Data) GetType() Action {
	if x != nil {
		return x.Type
	}
	return Action_Movement
}

func (x *Data) GetData() Direction {
	if x != nil {
		return x.Data
	}
	return Direction_UP
}

func (x *Data) GetRoomID() string {
	if x != nil {
		return x.RoomID
	}
	return ""
}

func (x *Data) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Data) GetPlayer() []*Player {
	if x != nil {
		return x.Player
	}
	return nil
}

type Room struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Player []*Player `protobuf:"bytes,3,rep,name=player,proto3" json:"player,omitempty"`
}

func (x *Room) Reset() {
	*x = Room{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movementEmitter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Room) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Room) ProtoMessage() {}

func (x *Room) ProtoReflect() protoreflect.Message {
	mi := &file_movementEmitter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Room.ProtoReflect.Descriptor instead.
func (*Room) Descriptor() ([]byte, []int) {
	return file_movementEmitter_proto_rawDescGZIP(), []int{1}
}

func (x *Room) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Room) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Room) GetPlayer() []*Player {
	if x != nil {
		return x.Player
	}
	return nil
}

type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	X    float32 `protobuf:"fixed32,2,opt,name=x,proto3" json:"x,omitempty"`
	Y    float32 `protobuf:"fixed32,3,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movementEmitter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_movementEmitter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_movementEmitter_proto_rawDescGZIP(), []int{2}
}

func (x *Player) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Player) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Player) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

var File_movementEmitter_proto protoreflect.FileDescriptor

var file_movementEmitter_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x6d, 0x69, 0x74, 0x74, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72,
	0x70, 0x63, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xae, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x72, 0x70, 0x63, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f,
	0x6f, 0x6d, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d,
	0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72,
	0x70, 0x63, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x22, 0x55, 0x0a, 0x04, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a,
	0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x22, 0x38, 0x0a, 0x06, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x01, 0x79, 0x2a, 0x34, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0c, 0x0a, 0x08,
	0x4d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x69,
	0x72, 0x65, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x10, 0x02, 0x12, 0x08,
	0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x10, 0x03, 0x2a, 0x3d, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x06, 0x0a, 0x02, 0x55, 0x50, 0x10, 0x00, 0x12, 0x08, 0x0a,
	0x04, 0x44, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x45, 0x46, 0x54, 0x10,
	0x02, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05,
	0x53, 0x48, 0x4f, 0x4f, 0x54, 0x10, 0x04, 0x32, 0xed, 0x01, 0x0a, 0x0f, 0x4d, 0x6f, 0x76, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x45, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72, 0x12, 0x46, 0x0a, 0x08, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x6f, 0x76, 0x65, 0x12, 0x0f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x72, 0x70, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0e, 0x12, 0x0c, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x4d, 0x6f, 0x76, 0x65, 0x28,
	0x01, 0x30, 0x01, 0x12, 0x4b, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f,
	0x6d, 0x12, 0x0f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x6f,
	0x6f, 0x6d, 0x1a, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x3a, 0x01, 0x2a,
	0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x6d,
	0x12, 0x45, 0x0a, 0x08, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x0f, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x1a, 0x0f, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x22, 0x17,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01, 0x2a, 0x22, 0x0c, 0x2f, 0x76, 0x31, 0x2f, 0x6a,
	0x6f, 0x69, 0x6e, 0x52, 0x6f, 0x6f, 0x6d, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x72, 0x69, 0x74, 0x65, 0x73, 0x68, 0x30, 0x34, 0x2f,
	0x73, 0x68, 0x6f, 0x6f, 0x74, 0x65, 0x72, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_movementEmitter_proto_rawDescOnce sync.Once
	file_movementEmitter_proto_rawDescData = file_movementEmitter_proto_rawDesc
)

func file_movementEmitter_proto_rawDescGZIP() []byte {
	file_movementEmitter_proto_rawDescOnce.Do(func() {
		file_movementEmitter_proto_rawDescData = protoimpl.X.CompressGZIP(file_movementEmitter_proto_rawDescData)
	})
	return file_movementEmitter_proto_rawDescData
}

var file_movementEmitter_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_movementEmitter_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_movementEmitter_proto_goTypes = []any{
	(Action)(0),    // 0: serverrpc.Action
	(Direction)(0), // 1: serverrpc.Direction
	(*Data)(nil),   // 2: serverrpc.Data
	(*Room)(nil),   // 3: serverrpc.Room
	(*Player)(nil), // 4: serverrpc.Player
}
var file_movementEmitter_proto_depIdxs = []int32{
	0, // 0: serverrpc.Data.type:type_name -> serverrpc.Action
	1, // 1: serverrpc.Data.data:type_name -> serverrpc.Direction
	4, // 2: serverrpc.Data.player:type_name -> serverrpc.Player
	4, // 3: serverrpc.Room.player:type_name -> serverrpc.Player
	2, // 4: serverrpc.MovementEmitter.SendMove:input_type -> serverrpc.Data
	3, // 5: serverrpc.MovementEmitter.CreateRoom:input_type -> serverrpc.Room
	3, // 6: serverrpc.MovementEmitter.JoinRoom:input_type -> serverrpc.Room
	2, // 7: serverrpc.MovementEmitter.SendMove:output_type -> serverrpc.Data
	4, // 8: serverrpc.MovementEmitter.CreateRoom:output_type -> serverrpc.Player
	3, // 9: serverrpc.MovementEmitter.JoinRoom:output_type -> serverrpc.Room
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_movementEmitter_proto_init() }
func file_movementEmitter_proto_init() {
	if File_movementEmitter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_movementEmitter_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Data); i {
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
		file_movementEmitter_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Room); i {
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
		file_movementEmitter_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Player); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_movementEmitter_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_movementEmitter_proto_goTypes,
		DependencyIndexes: file_movementEmitter_proto_depIdxs,
		EnumInfos:         file_movementEmitter_proto_enumTypes,
		MessageInfos:      file_movementEmitter_proto_msgTypes,
	}.Build()
	File_movementEmitter_proto = out.File
	file_movementEmitter_proto_rawDesc = nil
	file_movementEmitter_proto_goTypes = nil
	file_movementEmitter_proto_depIdxs = nil
}
