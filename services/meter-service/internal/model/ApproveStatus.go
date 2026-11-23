package model

type ApproveStatus string

const (
	PENDING  ApproveStatus = "PENDING"
	APPROVED ApproveStatus = "APPROVED"
	REFUSED  ApproveStatus = "REFUSED"
)