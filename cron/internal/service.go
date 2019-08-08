package internal

import (
	"errors"
	"k8s.io/api/batch/v1beta1"
	"net/http"
	"strings"
)

type Service interface {
	CreateCronJob(w http.ResponseWriter, r *http.Request)
	DeleteCronJob(w http.ResponseWriter, r *http.Request)
	GetCronJob(w http.ResponseWriter, r *http.Request)
}

type ChaosCronJob struct {
	Schedule string   `json:"schedule"`
	Name     string   `json:"name"`
	Cmd      []string `json:"command"`
}

type ChaosCronJobRequest struct {
	*ChaosCronJob
}

type ChaosCronJobResponse struct {
	ChaosCronJob *v1beta1.CronJob
}

func (c ChaosCronJobResponse) Render(w http.ResponseWriter, r *http.Request) error {
	//Preprocess the response before marshalling?
	return nil
}

func (c *ChaosCronJobRequest) Bind(r *http.Request) error {
	if c.ChaosCronJob == nil {
		return errors.New("Missing required chaos cron job fields.")
	}
	c.Name = strings.ToLower(c.Name)
	return nil
}
