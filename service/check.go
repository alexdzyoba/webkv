package service

import (
	"log"
	"time"
)

func (s *Service) UpdateTTL(check func() (bool, error)) {
	ticker := time.NewTicker(s.TTL / 2)
	for range ticker.C {
		s.update(check)
	}
}

func (s *Service) update(check func() (bool, error)) {
	ok, err := check()
	if !ok {
		log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		if agentErr := s.ConsulAgent.FailTTL("service:"+s.Name, err.Error()); agentErr != nil {
			log.Print(agentErr)
		}
	} else {
		if agentErr := s.ConsulAgent.PassTTL("service:"+s.Name, ""); agentErr != nil {
			log.Print(agentErr)
		}
	}
}
