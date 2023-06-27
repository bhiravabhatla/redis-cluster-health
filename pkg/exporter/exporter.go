package exporter

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thoas/go-funk"

	"github.com/bhiravabhatla/redis-cluster-health/pkg/utils"
)

type RedisCustomExporter struct {
	clusterStatus *prometheus.Desc
	clusterAddr   string
	namespace     string
	redisPassword string
}

func NewredisCustomExporter(clusterAddr, namespace, password string) *RedisCustomExporter {
	return &RedisCustomExporter{
		clusterStatus: prometheus.NewDesc("cluster_corrupted",
			"Shows whether redis nodes have corrupted cluster state. 0 implies cluster state is not corrupted.",
			nil, map[string]string{
				"cluster":   clusterAddr,
				"namespace": namespace,
			},
		),
		clusterAddr:   clusterAddr,
		namespace:     namespace,
		redisPassword: password,
	}
}

// Describe Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *RedisCustomExporter) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the metric you create for a given collector
	ch <- collector.clusterStatus
}

// Collect implements required collect function for all promehteus collectors
func (collector *RedisCustomExporter) Collect(ch chan<- prometheus.Metric) {

	var clusterCorrupted int

	redisEndpoints := utils.GetRedisEndpoints(context.Background(), collector.namespace, collector.clusterAddr)
	log.Println("current redis cluster endpoints are : ")
	log.Println(redisEndpoints)

	expectedClusterInfo := utils.GetRedisClusterDetails(redisEndpoints[0], collector.redisPassword)

	for _, endPoint := range redisEndpoints {
		isEqual := funk.Equal(expectedClusterInfo, utils.GetRedisClusterDetails(endPoint, collector.redisPassword))
		if !isEqual {
			clusterCorrupted++
		}
	}

	m1 := prometheus.MustNewConstMetric(collector.clusterStatus, prometheus.GaugeValue, float64(clusterCorrupted))
	m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
	ch <- m1
}
