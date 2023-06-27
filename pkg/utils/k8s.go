package utils

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetRedisEndpoints(ctx context.Context, namespace, redisClusterName string) (ips []string) {
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Fatalln(err.Error())
	}
	client, err := k8s.NewForConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	endpoints, err := client.CoreV1().Endpoints(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=" + redisClusterName,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, address := range endpoints.Items[0].Subsets[0].Addresses {
		ips = append(ips, address.IP)
	}
	return ips
}
