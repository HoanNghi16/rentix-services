package propertyclient

import "fmt"

// These protobuf v1-compatible messages follow proto/property.proto.
type Complex struct {
	ID            uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Avatar        string `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar"`
	Address       string `protobuf:"bytes,4,opt,name=address,proto3" json:"address"`
	Ward          string `protobuf:"bytes,5,opt,name=ward,proto3" json:"ward"`
	City          string `protobuf:"bytes,6,opt,name=city,proto3" json:"city"`
	TotalRooms    int32  `protobuf:"varint,7,opt,name=total_rooms,json=totalRooms,proto3" json:"total_rooms"`
	InactiveRooms int32  `protobuf:"varint,8,opt,name=inactive_rooms,json=inactiveRooms,proto3" json:"inactive_rooms"`
}

func (m *Complex) Reset()         { *m = Complex{} }
func (m *Complex) String() string { return fmt.Sprintf("Complex{id:%d,name:%q}", m.ID, m.Name) }
func (*Complex) ProtoMessage()    {}

type GetComplexesRequest struct {
	UserID uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Role   string `protobuf:"bytes,2,opt,name=role,proto3" json:"role"`
}

func (m *GetComplexesRequest) Reset() { *m = GetComplexesRequest{} }
func (m *GetComplexesRequest) String() string {
	return fmt.Sprintf("GetComplexesRequest{user_id:%d}", m.UserID)
}
func (*GetComplexesRequest) ProtoMessage() {}

type GetComplexesResponse struct {
	Complexes []*Complex `protobuf:"bytes,1,rep,name=complexes,proto3" json:"complexes"`
}

func (m *GetComplexesResponse) Reset() { *m = GetComplexesResponse{} }
func (m *GetComplexesResponse) String() string {
	return fmt.Sprintf("GetComplexesResponse{complexes:%d}", len(m.Complexes))
}
func (*GetComplexesResponse) ProtoMessage() {}

type Room struct {
	ComplexID      int64   `protobuf:"varint,1,opt,name=complex_id,json=complexId,proto3" json:"complex_id"`
	ID             int64   `protobuf:"varint,2,opt,name=id,proto3" json:"id"`
	CurrentEndDate string  `protobuf:"bytes,3,opt,name=current_end_date,json=currentEndDate,proto3" json:"current_end_date"`
	RoomNumber     string  `protobuf:"bytes,4,opt,name=room_number,json=roomNumber,proto3" json:"room_number"`
	Status         string  `protobuf:"bytes,5,opt,name=status,proto3" json:"status"`
	Price          float32 `protobuf:"fixed32,6,opt,name=price,proto3" json:"price"`
}

func (m *Room) Reset()         { *m = Room{} }
func (m *Room) String() string { return fmt.Sprintf("Room{id:%d,room_number:%q}", m.ID, m.RoomNumber) }
func (*Room) ProtoMessage()    {}

type GetRoomsRequest struct {
	ComplexID int64  `protobuf:"varint,1,opt,name=complex_id,json=complexId,proto3" json:"complex_id"`
	UserID    int64  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Role      string `protobuf:"bytes,3,opt,name=role,proto3" json:"role"`
}

func (m *GetRoomsRequest) Reset() { *m = GetRoomsRequest{} }
func (m *GetRoomsRequest) String() string {
	return fmt.Sprintf("GetRoomsRequest{complex_id:%d,user_id:%d}", m.ComplexID, m.UserID)
}
func (*GetRoomsRequest) ProtoMessage() {}

type GetRoomsResponse struct {
	Rooms []*Room `protobuf:"bytes,1,rep,name=rooms,proto3" json:"rooms"`
}

func (m *GetRoomsResponse) Reset() { *m = GetRoomsResponse{} }
func (m *GetRoomsResponse) String() string {
	return fmt.Sprintf("GetRoomsResponse{rooms:%d}", len(m.Rooms))
}
func (*GetRoomsResponse) ProtoMessage() {}

type CreateComplexResponse struct {
	Complex *Complex `protobuf:"bytes,1,opt,name=complex,proto3" json:"complex"`
}

func (m *CreateComplexResponse) Reset()         { *m = CreateComplexResponse{} }
func (m *CreateComplexResponse) String() string { return "CreateComplexResponse{}" }
func (*CreateComplexResponse) ProtoMessage()    {}

type CreateComplexRequest struct {
	UserID          uint64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Role            string    `protobuf:"bytes,2,opt,name=role,proto3" json:"role"`
	Name            string    `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
	Avatar          string    `protobuf:"bytes,4,opt,name=avatar,proto3" json:"avatar"`
	Address         string    `protobuf:"bytes,5,opt,name=address,proto3" json:"address"`
	EmployeeUserIds []uint64  `protobuf:"varint,6,rep,packed,name=employee_user_ids,json=employeeUserIds,proto3" json:"employee_user_ids"`
	Services        *Services `protobuf:"bytes,7,opt,name=services,proto3" json:"services"`
}

func (m *CreateComplexRequest) Reset() { *m = CreateComplexRequest{} }
func (m *CreateComplexRequest) String() string {
	return fmt.Sprintf("CreateComplexRequest{user_id:%d,name:%q}", m.UserID, m.Name)
}
func (*CreateComplexRequest) ProtoMessage() {}

type Services struct {
	Electricity *ServiceSetting `protobuf:"bytes,1,opt,name=electricity,proto3" json:"electricity"`
	Water       *ServiceSetting `protobuf:"bytes,2,opt,name=water,proto3" json:"water"`
	Parking     *ServiceSetting `protobuf:"bytes,3,opt,name=parking,proto3" json:"parking"`
	Garbage     *ServiceSetting `protobuf:"bytes,4,opt,name=garbage,proto3" json:"garbage"`
	Elevator    *ServiceSetting `protobuf:"bytes,5,opt,name=elevator,proto3" json:"elevator"`
	Internet    *ServiceSetting `protobuf:"bytes,6,opt,name=internet,proto3" json:"internet"`
	Other       *ServiceSetting `protobuf:"bytes,7,opt,name=other,proto3" json:"other"`
}

func (m *Services) Reset()         { *m = Services{} }
func (m *Services) String() string { return "Services{}" }
func (*Services) ProtoMessage()    {}

type ServiceSetting struct {
	Price float32 `protobuf:"fixed32,1,opt,name=price,proto3" json:"price"`
	Unit  string  `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit"`
}

func (m *ServiceSetting) Reset()         { *m = ServiceSetting{} }
func (m *ServiceSetting) String() string { return fmt.Sprintf("ServiceSetting{unit:%q}", m.Unit) }
func (*ServiceSetting) ProtoMessage()    {}
