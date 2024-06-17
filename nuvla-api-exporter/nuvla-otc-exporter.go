package nuvlaedge_otc_exporter

import (
	"context"
	"fmt"
	nuvla "github.com/nuvla/api-client-go"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"io"
	"net/http/httputil"
	"strings"
	"time"
)

var (
	typeStr = component.MustNewType("nuvlaedge_otc_exporter")
)

type NuvlaAPIExporter struct {
	cfg      *Config
	nuvlaApi *nuvla.NuvlaClient
	settings component.TelemetrySettings
}

func newNuvlaOTCExporter(
	cfg *Config,
	set *exporter.CreateSettings,
) (*NuvlaAPIExporter, error) {
	return &NuvlaAPIExporter{
		cfg:      cfg,
		nuvlaApi: nil,
		settings: set.TelemetrySettings,
	}, nil
}

func (e *NuvlaAPIExporter) Start(_ context.Context, _ component.Host) error {
	if e.cfg.Enabled {
		err := e.StartNuvlaApiClient()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *NuvlaAPIExporter) StartNuvlaApiClient() error {

	e.nuvlaApi = nuvla.NewNuvlaClientFromOpts(nil, nuvla.WithEndpoint(e.cfg.Endpoint),
		nuvla.WithInsecureSession(e.cfg.Insecure))
	err := e.nuvlaApi.LoginApiKeys(e.cfg.ApiKey, e.cfg.ApiSecret)

	if err != nil {
		e.settings.Logger.Error("Error logging in with api keys: ", zap.Error(err))
		return err
	}
	return nil
}

func (e *NuvlaAPIExporter) ConsumeMetrics(_ context.Context, pm pmetric.Metrics) error {
	rms := pm.ResourceMetrics()

	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		attrs := rm.Resource().Attributes()
		nuvlaDeploymentUUID, _ := attrs.Get("nuvla.deployment.uuid")
		serviceVal, _ := attrs.Get("service.name")
		serviceName := serviceVal.Str()

		scm := rm.ScopeMetrics()
		for j := 0; j < scm.Len(); j++ {
			sc := scm.At(j)

			ms := sc.Metrics()
			if ms.Len() == 0 {
				continue
			}
			uuid := nuvlaDeploymentUUID.Str()
			var metricMap []map[string]interface{}
			for k := 0; k < ms.Len(); k++ {
				currMetric := ms.At(k)
				e.updateMetric(&serviceName, &currMetric, &metricMap, &uuid)
			}
			if e.cfg.Enabled {
				err := e.sendMetricsToNuvla(&metricMap)
				if err != nil {
					e.settings.Logger.Error("Error sending metrics to Nuvla: ", zap.Error(err))
					return err
				}
			}
		}
	}
	return nil
}

func (e *NuvlaAPIExporter) sendMetricsToNuvla(metricMap *[]map[string]interface{}) error {
	e.settings.Logger.Info("Sending metrics to Nuvla: ", zap.Any("metricMap", metricMap))

	res, err := e.nuvlaApi.BulkOperation(e.cfg.ResourceId, "bulk-insert", *metricMap)
	if err != nil {
		e.settings.Logger.Error("Error in operation bulk insert: ", zap.Error(err))
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			e.settings.Logger.Error("Error closing the response body: ", zap.Error(err))
		}
	}(res.Body)

	dump, httperr := httputil.DumpResponse(res, true)
	if httperr != nil {
		e.settings.Logger.Error("Error dumping the response: ", zap.Error(httperr))
	}
	if res.StatusCode >= 400 {
		e.settings.Logger.Error("Error sending metrics to Nuvla: ", zap.String("dump", string(dump)))
		return fmt.Errorf("error sending metrics to Nuvla: %s", res.Status)
	}
	e.settings.Logger.Info("Metrics sent to Nuvla: ", zap.String("dump", string(dump)))
	return nil
}

func (e *NuvlaAPIExporter) updateMetric(serviceName *string, metric *pmetric.Metric,
	metricMap *[]map[string]interface{}, deploymentuuid *string) {

	var dp pmetric.NumberDataPointSlice

	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		dp = metric.Gauge().DataPoints()
	case pmetric.MetricTypeSum:
		dp = metric.Sum().DataPoints()
	default:
		panic("unhandled default case")
	}

	metricName := metric.Name()
	metricName, _ = strings.CutPrefix(metricName, *serviceName+"_")
	for i := 0; i < dp.Len(); i++ {
		var currMetricMap = make(map[string]interface{})
		datapoint := dp.At(i)
		e.settings.Logger.Info("Datapoint ", zap.Any("datapoint", datapoint))
		// TODO there could be situations of timestamps being very close or same.
		// Need to handle that.
		timestamp := datapoint.Timestamp().AsTime()
		timestamp = timestamp.Add(time.Millisecond * time.Duration(i))

		currMetricMap["timestamp"] = timestamp.Format(time.RFC3339)
		//time.Sleep(2 * time.Millisecond)
		currMetricMap["nuvla.deployment.uuid"] = *deploymentuuid
		currMetricMap["service.name"] = *serviceName

		switch datapoint.ValueType() {
		case pmetric.NumberDataPointValueTypeInt:
			currMetricMap[metricName] = datapoint.IntValue()
		case pmetric.NumberDataPointValueTypeDouble:
			currMetricMap[metricName] = datapoint.DoubleValue()
		default:
			panic("unhandled default case")
		}

		datapoint.Attributes().Range(func(k string, v pcommon.Value) bool {
			currMetricMap[k] = v.AsString()
			return true
		})
		// TODO this is added to make the metric unique. Need to find a better way.
		currMetricMap["metricInfo"] = metricName
		*metricMap = append(*metricMap, currMetricMap)
	}
	e.settings.Logger.Info("MetricMap in updateMetric", zap.Any("metricMap", metricMap))
}
