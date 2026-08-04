package pakasir_service

import (
	"raznar.id/invoice-broker/configs"
	base_service "raznar.id/invoice-broker/internal/services/base"
)

type PakasirService struct {
	base_service.Service
}

func New(c *configs.Config) *PakasirService {
	v := &PakasirService{}
	v.Set(c)

	return v
}
