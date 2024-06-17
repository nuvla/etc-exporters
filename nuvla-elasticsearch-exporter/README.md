# nuvla-elasticsearch-exporter
An OTC exporter which would convert data into appropriate format and send to elasticsearch endpoint

For elasticsearch, the exporter would create a timeseries resource (<index_prefix>-<application_name>) if not present.
Then all the metric datapoints are stored. We take the data attributes along with metrics data and nuvla.deployment.uuid.
So, the exporter cannot in the current version create timeseries resource dynamically based on the data sent by the
application. We should provide in the configuration the list of metrics along with the data attributes.


## Configuration options

- elasticsearch
  - enabled: Enable/Disable the exporter (default: false)
  - endpoint: Elastic search endpoint (default: http://localhost:9200)
  - insecure: Skip SSL verification (default: false)
  - ca_file: CA certificate file for verifying the server certificate (default: "")
  - index_prefix: Prefix for the index patterns and templates used for timeseries resource
                   (default: "nuvla-opentelemetry-")
  - MetricsTobeExported: List of metrics to be exported to elastic search 
    Format : <metricName>,<metricType>,<isDimension>
       Dimension metrics along with the timestamp would define the unique document in elastic search
       Example : "cpu.usage,counter,false"

- sending_queue: Queue Settings (enabled by ) (
  Follow the Queue Settings config in https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)
- retry_on_failure : Retry on failure (Follow the retry on failure configuration in
  https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)