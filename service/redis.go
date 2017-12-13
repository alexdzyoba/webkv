package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 200

	key := strings.Trim(r.URL.Path, "/")
	val, err := s.RedisClient.Get(key).Result()
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		status = 404
	}

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
