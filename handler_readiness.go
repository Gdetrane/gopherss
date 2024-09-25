package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request)  {
  type healthRespPayload struct {
    ApiVersion string `json:"ApiVersion"`
    HealthStatus string `json:"HealthStatus"`
  }
  respondWithJSON(w, 200, healthRespPayload{
    ApiVersion: "v1",
    HealthStatus: "OK",
  })
}
