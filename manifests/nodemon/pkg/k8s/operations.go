package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
	"strconv"
	"strings"
)

type PodData struct {
	Index int64
	IP    string
	Name  string
}

func GetPodIps(ctx context.Context, label, ns string, client *kubernetes.Clientset) ([]PodData, error) {
	pods, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods for label=%v, ns=%v, error=%w", label, ns, err)
	}
	podData := []PodData{}
	stsName := ""
	for _, pod := range pods.Items {
		if stsName == "" {
			owners := pod.GetOwnerReferences()
			if len(owners) > 0 {
				stsName = owners[0].Name
			} else {
				return nil, fmt.Errorf("no owner ref found for assumes sts pod")
			}
		}
		index, err := strconv.ParseInt(strings.Replace(pod.Name, stsName+"-", "", -1), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sts pod index in pod=%v, error=%w", pod.Name, err)
		}
		podData = append(podData, PodData{index, pod.Status.PodIP, pod.Name})
	}
	sort.Slice(podData, func(i, j int) bool {
		return podData[i].Index < podData[j].Index
	})

	return podData, nil
}
