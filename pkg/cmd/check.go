/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/bhiravabhatla/redis-cluster-health/pkg/exporter"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the cluster integrity of redis",
	Run: func(cmd *cobra.Command, args []string) {
		password := os.Getenv("REDIS_PASSWORD")
		cluster, _ := cmd.Flags().GetString("cluster")
		namespace, _ := cmd.Flags().GetString("namespace")
		redisExporter := exporter.NewredisCustomExporter(cluster, namespace, password)
		prometheus.MustRegister(redisExporter)

		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringP("cluster", "c", "", "redis cluster service name")
	checkCmd.Flags().StringP("namespace", "n", "master", "namespace of redis cluster")
	err := checkCmd.MarkFlagRequired("cluster")
	if err != nil {
		fmt.Println(err)
	}
	err = checkCmd.MarkFlagRequired("namespace")
	if err != nil {
		fmt.Println(err)
	}
}
