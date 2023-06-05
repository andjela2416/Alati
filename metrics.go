package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// Initial count.
	currentCount = 0

	// The Prometheus metric that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_app_http_hit_total",
			Help: "Total number of http hits.",
		},
	)
	createConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_config_http_hit_total",
			Help: "Total number of create config hits.",
		},
	)

	getAllConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_config_http_hit_total",
			Help: "Total number of get all config hits.",
		},
	)

	getConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_http_hit_total",
			Help: "Total number of get config hits.",
		},
	)

	delConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_config_http_hit_total",
			Help: "Total number of del config hits.",
		},
	)

	createGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_group_http_hit_total",
			Help: "Total number of create group hits.",
		},
	)

	getAllGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_group_http_hit_total",
			Help: "Total number of get all group hits.",
		},
	)

	getGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_group_http_hit_total",
			Help: "Total number of get group hits.",
		},
	)

	delGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_group_http_hit_total",
			Help: "Total number of del group hits.",
		},
	)

	appendGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "append_group_http_hit_total",
			Help: "Total number of append group hits.",
		},
	)

	getConfigByLabelsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_by_labels_http_hit_total",
			Help: "Total number of get config by labels hits.",
		},
	)

	delConfigByLabelsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_config_by_labels_http_hit_total",
			Help: "Total number of del config by labels hits.",
		},
	)

	getGroupByLabelsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_group_by_labels_http_hit_total",
			Help: "Total number of get group by labels hits.",
		},
	)

	delGroupByLabelsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "del_group_by_labels_http_hit_total",
			Help: "Total number of del group by labels hits.",
		},
	)
	addConfigToGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "add_config_to_group_http_hit_total",
			Help: "Total number of add config to group hits.",
		},
	)
	swaggerHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "swagger_http_hit_total",
			Help: "Total number of swagger hits.",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		httpHits,
		createConfigHits,
		getAllConfigHits,
		getConfigHits,
		delConfigHits,
		createGroupHits,
		getAllGroupHits,
		getGroupHits,
		delGroupHits,
		appendGroupHits,
		getConfigByLabelsHits,
		delConfigByLabelsHits,
		getGroupByLabelsHits,
		delGroupByLabelsHits,
		addConfigToGroupHits,
		swaggerHits,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		f(w, r) // original function call
	}
}

func CountCreateConfig(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		createConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetAllConfig(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getAllConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetConfig(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelConfig(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigHits.Inc()
		f(w, r) // original function call
	}
}

func CountCreateGroup(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		createGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetAllGroup(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getAllGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetGroup(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelGroup(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountAppendGroup(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		appendGroupHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetConfigByLabels(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigByLabelsHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelConfigByLabels(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigByLabelsHits.Inc()
		f(w, r) // original function call
	}
}

func CountGetGroupByLabels(f func(w http.ResponseWriter, req *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigByLabelsHits.Inc()
		f(w, r) // original function call
	}
}

func CountDelGroupByLabels(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		delConfigByLabelsHits.Inc()
		f(w, r) // original function call
	}
}

func CountAddConfigToGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		addConfigToGroupHits.Inc()
		f(w, r) // original function call
	}
}
func SwaggerHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		swaggerHits.Inc()
		f(w, r) // original function call
	}
}
