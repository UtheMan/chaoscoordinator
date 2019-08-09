package service

import (
	"github.com/go-chi/render"
	"github.com/utheman/chaoscoordinator/cron/internal"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

type CronJobService struct {
	ClientSet *kubernetes.Clientset
}

func (s *CronJobService) CreateCronJob(w http.ResponseWriter, r *http.Request) {
	data := &internal.ChaosCronJobRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequest(err))
		return
	}
	if err := deployCronJob(data.ChaosCronJob, s.ClientSet); err != nil {
		render.Render(w, r, InvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
}

func (s *CronJobService) DeleteCronJob(w http.ResponseWriter, r *http.Request) {
	err := s.ClientSet.BatchV1beta1().CronJobs("default").Delete("test", &metav1.DeleteOptions{})
	if err != nil {
		render.Render(w, r, InvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
}

func (s *CronJobService) GetCronJob(w http.ResponseWriter, r *http.Request) {
	var cronJob = &v1beta1.CronJob{}
	cronJob, err := s.ClientSet.BatchV1beta1().CronJobs("default").Get("test", metav1.GetOptions{})
	if err != nil {
		render.Render(w, r, InvalidRequest(err))
		return
	}
	if err := render.Render(w, r, NewCronJobResponse(cronJob)); err != nil {
		render.Render(w, r, InvalidRender(err))
		return
	}
}

func deployCronJob(job *internal.ChaosCronJob, clientset *kubernetes.Clientset) error {
	testCronJob := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1beta1.CronJobSpec{},
		Status:     v1beta1.CronJobStatus{},
	}
	testCronJob.ObjectMeta.Name = job.Name
	testCronJob.Spec.Schedule = job.Schedule
	testChaosContainer := v1.Container{
		Name:  job.Name,
		Image: "utheman/utheman_chaoscoordinator:175adfc-dirty",
		Args:  job.Cmd,
	}
	testCronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
	testCronJob.Spec.JobTemplate.Spec.Template.Spec.Containers = append(testCronJob.Spec.JobTemplate.Spec.Template.Spec.Containers, testChaosContainer)
	_, err := clientset.BatchV1beta1().CronJobs("default").Create(testCronJob)
	if err != nil {
		return err
	}
	return nil
}

func NewCronJobResponse(job *v1beta1.CronJob) *internal.ChaosCronJobResponse {
	response := &internal.ChaosCronJobResponse{ChaosCronJob: job}
	return response
}
