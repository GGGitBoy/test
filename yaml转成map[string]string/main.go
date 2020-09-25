package main

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/util/yaml"
)

var data = `
---
  test:
  - aa
  - bb
  enabled: "false"
  exporter-coredns: 
    apiGroup: "monitoring.coreos.com"
  exporter-fluentd: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
  exporter-gpu-node: 
    enabled: "false"
  exporter-kube-controller-manager: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
    endpoints: 
      - "172.31.11.115"
      - "172.31.11.218"
      - "172.31.6.109"
  exporter-kube-dns: 
    apiGroup: "monitoring.coreos.com"
  exporter-kube-etcd: 
    apiGroup: "monitoring.coreos.com"
    certFile: "/etc/prometheus/secrets/exporter-etcd-cert/kube-etcd-172-31-11-115.pem"
    enabled: "true"
    endpoints: 
      - "172.31.11.115"
      - "172.31.11.218"
      - "172.31.6.109"
    keyFile: "/etc/prometheus/secrets/exporter-etcd-cert/kube-etcd-172-31-11-115-key.pem"
    ports: 
      metrics: 
        port: "2379"
  exporter-kube-scheduler: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
    endpoints: 
      - "172.31.11.115"
      - "172.31.11.218"
      - "172.31.6.109"
  exporter-kube-state: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
  exporter-kubelets: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
    https: "true"
  exporter-kubernetes: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
  exporter-node: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
    ports: 
      metrics: 
        port: "9796"
    resources: 
      limits: 
        cpu: "200m"
        memory: "50Mi"
  grafana: 
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
    persistence: 
      enabled: "false"
      size: "10Gi"
      storageClass: "default"
    serviceAccountName: "cluster-monitoring"
  operator-init: 
    enabled: "true"
  operator: 
    resources: 
      limits: 
        memory: "100Mi"
  prometheus: 
    additionalAlertManagerConfigs: 
      - 
        kubernetes_sd_configs: 
          - 
            role: "endpoints"
        relabel_configs: 
          - 
            action: "keep"
            regex: "cattle-prometheus;access-alertmanager"
            source_labels: 
              - "__meta_kubernetes_namespace"
              - "__meta_kubernetes_endpoints_name"
    apiGroup: "monitoring.coreos.com"
    enabled: "true"
    externalLabels: 
      prometheus_from: "test"
    persistence: 
      enabled: "false"
      size: "50Gi"
      storageClass: "default"
    persistent: 
      useReleaseName: "true"
    resources: 
      core: 
        limits: 
          cpu: "1000m"
          memory: "1000Mi"
        requests: 
          cpu: "750m"
          memory: "750Mi"
    retention: "12h"
    ruleSelector: 
      matchExpressions: 
        - 
          key: "source"
          operator: "In"
          values: 
            - "rancher-alert"
            - "rancher-monitoring"
    secrets: 
      - "exporter-etcd-cert"
    serviceAccountNameOverride: "cluster-monitoring"
`

func main() {
	jsonData, err := yaml.ToJSON([]byte(data))
	if err != nil {
		panic(err)
	}

	var mapData map[string]interface{}
	err = json.Unmarshal(jsonData, &mapData)
	if err != nil {
		fmt.Printf("解码失败：%v\n", err)
		return
	}

	answer := make(map[string]string)

	for k, v := range mapData {
		getAnswer(k, v, answer)
	}

	for ka, va := range answer {
		fmt.Printf("%s  :  %s\n\n", ka, va)
	}

}

func getAnswer(k string, v interface{}, answer map[string]string) {
	if data, ok := v.(map[string]interface{}); ok {
		// 当值 v 为 map[string]interface{} 时
		for kAdd, vNew := range data {
			kNew := k + "." + kAdd
			getAnswer(kNew, vNew, answer)
		}
	} else if data, ok := v.([]interface{}); ok {
		// 当值 v 为 []interface{} 时
		for index, vNew := range data {
			kNew := fmt.Sprintf("%s[%v]", k, index)
			getAnswer(kNew, vNew, answer)
		}
	} else {
		// 当值 v 为其它
		str := fmt.Sprintf("%v", v)
		answer[k] = str
	}
}
