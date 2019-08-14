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
	cronJob := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1beta1.CronJobSpec{},
		Status:     v1beta1.CronJobStatus{},
	}
	cronJob.ObjectMeta.Name = job.Name
	cronJob.Spec.Schedule = job.Schedule
	cronJob.Spec.FailedJobsHistoryLimit = new(int32)
	cronJob.Spec.SuccessfulJobsHistoryLimit = new(int32)
	testChaosContainer := v1.Container{
		Name:    job.Name,
		Image:   "utheman/utheman_chaoscoordinator:4124186-dirty",
		Command: job.Cmd,
		Args:    job.Args,
	}
	cronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
	azureCreds := v1.KeyToPath{
		Key:  "creds",
		Path: "creds",
	}
	secretVolume := v1.Volume{
		Name: "azure-auth-volume",
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: "azure-auth",
				Items:      []v1.KeyToPath{azureCreds},
			},
		},
	}
	volumeMount := v1.VolumeMount{
		Name:      "azure-auth-volume",
		ReadOnly:  true,
		MountPath: "/etc/azure-auth-volume",
	}
	azureAuthLocation := v1.EnvVar{
		Name:  "AZURE_AUTH_LOCATION",
		Value: "/etc/azure-auth-volume/creds",
	}
	azureSubscriptionId := v1.EnvVar{
		Name: "SUBSCRIPTION_ID",
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "azure-subscription-id",
				},
				Key: "subscriptionId",
			},
		},
	}
	testChaosContainer.Env = append(testChaosContainer.Env, azureAuthLocation)
	testChaosContainer.Env = append(testChaosContainer.Env, azureSubscriptionId)
	testChaosContainer.VolumeMounts = append(testChaosContainer.VolumeMounts, volumeMount)
	cronJob.Spec.JobTemplate.Spec.Template.Spec.Volumes = append(cronJob.Spec.JobTemplate.Spec.Template.Spec.Volumes, secretVolume)
	cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers = append(cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers, testChaosContainer)
	_, err := clientset.BatchV1beta1().CronJobs("default").Create(cronJob)
	if err != nil {
		return err
	}
	return nil
}

func NewCronJobResponse(job *v1beta1.CronJob) *internal.ChaosCronJobResponse {
	response := &internal.ChaosCronJobResponse{ChaosCronJob: job}
	return response
}

func NewCronJobListResponse(list *v1beta1.CronJobList) *internal.ChaosCronJobListResponse {
	response := &internal.ChaosCronJobListResponse{CronJobList: list}
	return response
}
