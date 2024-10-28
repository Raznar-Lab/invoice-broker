package services

import services_invoice "raznar.id/invoice-broker/internal/app/services/invoice"

func Invoice() *services_invoice.InvoiceService {
	return &services_invoice.InvoiceService{}
}
