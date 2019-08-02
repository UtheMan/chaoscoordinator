package main

import (
	"flag"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	//"k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "genesys-config-k8s-genesys"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}


	//create test cron job
	testCronJob := &v1beta1.CronJob{
		TypeMeta:   metav1.TypeMeta{"CronJob", "batch/v1beta1"},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1beta1.CronJobSpec{},
		Status:     v1beta1.CronJobStatus{},
	}
	testCronJob.TypeMeta.APIVersion = "batch/v1beta1"
	testCronJob.TypeMeta.Kind = "CronJob"
	testCronJob.ObjectMeta.Name = "test"
	testCronJob.Spec.Schedule = "*/1 * * * *"
	
	testChaosContainer :=  v1.Container{}
	testChaosContainer.Name = "test"
	testChaosContainer.Image = "utheman/utheman_chaoscoordinator:599f0ea-dirty"
	testChaosContainer.Args = [] string {"./chaos", "vm", "kill", "-m random", "-r controlplane", "-"}
	testCronJob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyOnFailure
	testCronJob.Spec.JobTemplate.Spec.Template.Spec.Containers = append(testCronJob.Spec.JobTemplate.Spec.Template.Spec.Containers, testChaosContainer)

	cronJob, _ := clientset.BatchV1beta1().CronJobs("default").Create(testCronJob)
	println(cronJob)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
