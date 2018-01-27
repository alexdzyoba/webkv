package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 200

	key := strings.Trim(r.URL.Path, "/")

	timer := prometheus.NewTimer(s.Metrics.RedisDurations)
	defer timer.ObserveDuration()

	val, err := s.RedisClient.Get(key).Result()
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		status = 404
		s.Metrics.RedisRequests.WithLabelValues("fail").Inc()
	}
	s.Metrics.RedisRequests.WithLabelValues("success").Inc()

	fmt.Fprint(w, val)
	log.Printf("url=\"%s\" remote=\"%s\" key=\"%s\" status=%d\n",
		r.URL, r.RemoteAddr, key, status)
}

func (s *Service) Check() (bool, error) {
	_, err := s.RedisClient.Ping().Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
