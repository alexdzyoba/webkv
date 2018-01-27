package service

import (
	"time"

	"github.com/go-redis/redis"
	consul "github.com/hashicorp/consul/api"
)

type Service struct {
	Name        string
	TTL         time.Duration
	Port        int
	RedisClient redis.UniversalClient
	ConsulAgent *consul.Agent
}

type Metrics struct {
	RedisRequests *prometheus.CounterVec
}

func New(addrs []string, ttl time.Duration, port int) (*Service, error) {
	s := new(Service)
	s.Name = "webkv"
	s.Port = port
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
		Port: s.Port,
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
