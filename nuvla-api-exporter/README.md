# nuvla-api-exporter
An OTC exporter which would send data to Nuvla API. The users need to give api keys and secrets to authenticate with the Nuvla API.
The timeseries resource needs to be created before using this exporter and resource_id should be provided in the configuration.


## Configuration options

- nuvla
  - enabled: Enable/Disable the exporter (default: false)
  - insecure: Skip server certificate check (default: true)
  - endpoint: Nuvla endpoint (default: https://nuvla.io)
  - api_key: Nuvla API key
  - api_secret: Nuvla API secret
  - resource_id: Resource ID of the timeseries resource to which the data should be sent

- sending_queue: Queue Settings (enabled by ) (
  Follow the Queue Settings config in https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)
- retry_on_failure : Retry on failure (Follow the retry on failure configuration in
  https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)