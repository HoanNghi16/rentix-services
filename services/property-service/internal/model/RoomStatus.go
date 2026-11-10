package model

type RoomStatus string

const (
	RoomEmpty       RoomStatus = "EMPTY"
	RoomOccupied    RoomStatus = "OCCUPIED"
	RoomMaintenance RoomStatus = "MAINTENANCE"
)