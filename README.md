# modbus_http_exporter

This is for receiving datapoints from the Data to server HTTP service from RutOS (Teltonika network devices).

## Tested devices:
- RUT240

## RutOS configuration:
Assuming that a device is set up that is producing Modbus datapoints to the Network device.

1. Set up an HTTP 'Data to server' configuration.
2. Set the address to this exporter address.
3. Make sure the format is set to json
4. Set 'Send as object' to True.

Datapoints sent from the HTTP service will be stored in memory and exposed at '/metrics' 
