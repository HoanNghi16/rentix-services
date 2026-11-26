package grpcapi

import "fmt"

// These protobuf v1-compatible messages match proto/employee.proto.
type GetEmployeesRequest struct {
	OwnerID uint64 `protobuf:"varint,1,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	Role    string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (m *GetEmployeesRequest) Reset() { *m = GetEmployeesRequest{} }
func (m *GetEmployeesRequest) String() string {
	return fmt.Sprintf("GetEmployeesRequest{owner_id:%d}", m.OwnerID)
}
func (*GetEmployeesRequest) ProtoMessage() {}

type EmployeeAssignment struct {
	ComplexID uint64 `protobuf:"varint,1,opt,name=complex_id,json=complexId,proto3" json:"complex_id,omitempty"`
}

func (m *EmployeeAssignment) Reset() { *m = EmployeeAssignment{} }
func (m *EmployeeAssignment) String() string {
	return fmt.Sprintf("EmployeeAssignment{complex_id:%d}", m.ComplexID)
}
func (*EmployeeAssignment) ProtoMessage() {}

type Employee struct {
	UserID      uint64                `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FullName    string                `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Email       string                `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phone       string                `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	Assignments []*EmployeeAssignment `protobuf:"bytes,5,rep,name=assignments,proto3" json:"assignments,omitempty"`
}

func (m *Employee) Reset()         { *m = Employee{} }
func (m *Employee) String() string { return fmt.Sprintf("Employee{user_id:%d}", m.UserID) }
func (*Employee) ProtoMessage()    {}

type GetEmployeesResponse struct {
	Employees []*Employee `protobuf:"bytes,1,rep,name=employees,proto3" json:"employees,omitempty"`
}

func (m *GetEmployeesResponse) Reset() { *m = GetEmployeesResponse{} }
func (m *GetEmployeesResponse) String() string {
	return fmt.Sprintf("GetEmployeesResponse{employees:%d}", len(m.Employees))
}
func (*GetEmployeesResponse) ProtoMessage() {}
