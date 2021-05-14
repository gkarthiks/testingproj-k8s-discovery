package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	discovery "github.com/gkarthiks/k8s-discovery"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsTypes "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

var (
	k8s *discovery.K8s
)

func main() {
	k8s, _ = discovery.NewK8s()
	namespace, _ := k8s.GetNamespace()
	version, _ := k8s.GetVersion()
	fmt.Printf("Specified Namespace: %s\n", namespace)
	fmt.Printf("Version of running Kubernetes: %s\n", version)

	pods, err := k8s.Clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	//cronJobs, err := k8s.Clientset.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	}
	for idx, crons := range pods.Items {
		fmt.Printf("%d -> %s\n", idx, crons.Name)
	}

	fmt.Println("-----------------------Done with the pod listing; moving towards the metrics-----------------------")

	podMetrics, err := k8s.MetricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	var podMetric metricsTypes.PodMetrics
	getPodUsageMetrics := func(pod metricsTypes.PodMetrics) {
		for _, container := range pod.Containers {
			cpuQuantityDec := container.Usage.Cpu().AsDec().String()
			cpuUsageFloat, _ := strconv.ParseFloat(cpuQuantityDec, 64)

			fmt.Printf("CPU Usage Float: %v\n", cpuUsageFloat)

			memoryQuantityDec := container.Usage.Memory().AsDec().String()
			memoryUsageFloat, _ := strconv.ParseFloat(memoryQuantityDec, 64)
			fmt.Printf("Memory Usage Float: %v\n\n", memoryUsageFloat)
		}
	}

	for _, podMetric = range podMetrics.Items {
		getPodUsageMetrics(podMetric)
	}

}
