package model

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "PENDING"
	InvoiceStatusUnpaid  InvoiceStatus = "UNPAID"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
)