package service

import (
	"time"

	"github.com/go-redis/redis"
	consul "github.com/hashicorp/consul/api"
)

type Service struct {
	Name        string
	TTL         time.Duration
	RedisClient redis.UniversalClient
	ConsulAgent *consul.Agent
}

func New(addrs []string, ttl time.Duration) (*Service, error) {
	s := new(Service)
	s.Name = "webkv"
	s.TTL = ttl
	s.RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: addrs,
	})

	ok, err := s.Check()
	if !ok {
		return nil, err
	}

	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	s.ConsulAgent = c.Agent()

	serviceDef := &consul.AgentServiceRegistration{
		Name: s.Name,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}

	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}
	go s.UpdateTTL(s.Check)

	return s, nil
}
