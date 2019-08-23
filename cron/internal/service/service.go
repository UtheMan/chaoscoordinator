package service

import (
	"chaosmanager/internal"
	"github.com/go-chi/render"
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
	name := r.FormValue("name")
	if name == "" {
		s.ParamsNotPresent(w, r)
	}
	err := s.ClientSet.BatchV1beta1().CronJobs("default").Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		render.Render(w, r, ContentNotFoundRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
}

func (s *CronJobService) GetCronJob(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	namespace := r.FormValue("namespace")
	if name != "" {
		s.getSingleCronJob(&v1beta1.CronJob{}, name, w, r)
	} else if namespace != "" {
		s.getCronJobList(&v1beta1.CronJobList{}, namespace, w, r)
	} else {
		s.ParamsNotPresent(w, r)
	}
}

func (s *CronJobService) getSingleCronJob(cronJob *v1beta1.CronJob, name string, w http.ResponseWriter, r *http.Request) {
	cronJob, err := s.ClientSet.BatchV1beta1().CronJobs("default").Get(name, metav1.GetOptions{})
	if err != nil {
		render.Render(w, r, ContentNotFoundRequest(err))
		return
	}
	if err := render.Render(w, r, NewCronJobResponse(cronJob)); err != nil {
		render.Render(w, r, InvalidRender(err))
		return
	}
}

func (s *CronJobService) getCronJobList(list *v1beta1.CronJobList, namespace string, w http.ResponseWriter, r *http.Request) {
	list, err := s.ClientSet.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		render.Render(w, r, ContentNotFoundRequest(err))
		return
	}
	if err := render.Render(w, r, NewCronJobListResponse(list)); err != nil {
		render.Render(w, r, InvalidRender(err))
		return
	}
}

func deployCronJob(job *internal.ChaosCronJob, clientset *kubernetes.Clientset) error {
	cronJob := setupCronJob(job)
	_, err := clientset.BatchV1beta1().CronJobs("default").Create(cronJob)
	if err != nil {
		return err
	}
	return nil
}

func setupCronJob(job *internal.ChaosCronJob) *v1beta1.CronJob {
	cronJob := initCronJob(job)
	chaosContainer := setupContainer(job)
	cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers = append(cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers, chaosContainer)
	return cronJob
}

func setupContainer(job *internal.ChaosCronJob) v1.Container {
	container := v1.Container{
		Name:    job.Name,
		Image:   "utheman/utheman_chaoscoordinator:468c33c-dirty",
		Command: job.Cmd,
		Args:    job.Args,
	}
	container.Env = append(container.Env, v1.EnvVar{
		Name:  "AZURE_AUTH_LOCATION",
		Value: "/etc/azure-auth-volume/creds",
	})
	container.Env = append(container.Env, v1.EnvVar{
		Name: "SUBSCRIPTION_ID",
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "azure-subscription-id",
				},
				Key: "subscriptionId",
			},
		},
	})
	container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
		Name:      "azure-auth-volume",
		ReadOnly:  true,
		MountPath: "/etc/azure-auth-volume",
	})
	return container
}

func initCronJob(job *internal.ChaosCronJob) *v1beta1.CronJob {
	cronJob := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.Name,
		},
		Spec: v1beta1.CronJobSpec{
			Schedule:                   job.Schedule,
			FailedJobsHistoryLimit:     new(int32),
			SuccessfulJobsHistoryLimit: new(int32),
		},
		Status: v1beta1.CronJobStatus{},
	}
	cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
	mountCronJobAuthVolumes(cronJob)
	return cronJob
}

func mountCronJobAuthVolumes(job *v1beta1.CronJob) {
	job.Spec.JobTemplate.Spec.Template.Spec.Volumes = append(job.Spec.JobTemplate.Spec.Template.Spec.Volumes, v1.Volume{
		Name: "azure-auth-volume",
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: "azure-auth",
				Items: []v1.KeyToPath{{
					Key:  "creds",
					Path: "creds",
				}},
			},
		},
	})
}

func NewCronJobResponse(job *v1beta1.CronJob) *internal.ChaosCronJobResponse {
	response := &internal.ChaosCronJobResponse{ChaosCronJob: job}
	return response
}

func NewCronJobListResponse(list *v1beta1.CronJobList) *internal.ChaosCronJobListResponse {
	response := &internal.ChaosCronJobListResponse{CronJobList: list}
	return response
}
