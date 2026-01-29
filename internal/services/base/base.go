package base_service

import "raznar.id/invoice-broker/configs"

type Service struct {
	Config *configs.Config
}

func (s *Service) Set(c *configs.Config) {
	s.Config = c
}
