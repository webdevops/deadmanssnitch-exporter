package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	prometheusCommon "github.com/webdevops/go-prometheus-common"
)

type MetricsCollectorSnitch struct {
	CollectorProcessorGeneral

	prometheus struct {
		snitchInfo      *prometheus.GaugeVec
		snitchStatus    *prometheus.GaugeVec
		snitchHeartbeat *prometheus.GaugeVec
	}
}

func (m *MetricsCollectorSnitch) Setup(collector *CollectorGeneral) {
	m.CollectorReference = collector

	m.prometheus.snitchInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deadmanssnitch_snitch_info",
			Help: "DeadMansSnitch snitch information",
		},
		[]string{
			"snitchName",
			"interval",
			"alertType",
		},
	)

	m.prometheus.snitchStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deadmanssnitch_snitch_status",
			Help: "DeadMansSnitch snitch status",
		},
		[]string{
			"snitchName",
		},
	)
	m.prometheus.snitchHeartbeat = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deadmanssnitch_snitch_heartbeat",
			Help: "DeadMansSnitch snitch last heartbeat update",
		},
		[]string{
			"snitchName",
		},
	)

	prometheus.MustRegister(m.prometheus.snitchInfo)
	prometheus.MustRegister(m.prometheus.snitchStatus)
	prometheus.MustRegister(m.prometheus.snitchHeartbeat)
}

func (m *MetricsCollectorSnitch) Reset() {
	m.prometheus.snitchInfo.Reset()
	m.prometheus.snitchStatus.Reset()
	m.prometheus.snitchHeartbeat.Reset()
}

func (m *MetricsCollectorSnitch) Collect(ctx context.Context, callback chan<- func()) {
	list, err := DmsClient.ListSnitches()
	m.CollectorReference.PrometheusAPICounter().WithLabelValues("ListSnitches").Inc()

	if err != nil {
		m.logger().Panic(err)
	}

	snitchInfoMetricList := prometheusCommon.NewMetricsList()
	snitchStatusMetricList := prometheusCommon.NewMetricsList()
	snitchHeartbeatMetricList := prometheusCommon.NewMetricsList()

	for _, snitch := range list {
		snitchInfoMetricList.AddInfo(prometheus.Labels{
			"snitchName": snitch.Name,
			"interval":   snitch.Interval,
			"alertType":  snitch.AlertType,
		})

		snitchStatusMetricList.AddBool(prometheus.Labels{
			"snitchName": snitch.Name,
		}, snitch.IsHealthy())

		if snitch.CheckedInAt != nil && snitch.CheckedInAt.Unix() > 0 {
			snitchHeartbeatMetricList.AddTime(prometheus.Labels{
				"snitchName": snitch.Name,
			}, *snitch.CheckedInAt)
		}
	}

	// set metrics
	callback <- func() {
		snitchInfoMetricList.GaugeSet(m.prometheus.snitchInfo)
		snitchStatusMetricList.GaugeSet(m.prometheus.snitchStatus)
		snitchHeartbeatMetricList.GaugeSet(m.prometheus.snitchHeartbeat)
	}
}
