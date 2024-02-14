AirGradient Prometheus Exporter
===============================

Receive data from AirGradient sensors and expose it as Prometheus metrics.

# AirGradient Device Configuration
For the device to send data to the exporter, it needs to be configured to send data to the exporter's IP address and port. This will require flashing the device with a custom sketch. 

The following steps will guide you through the process of preparing your computer to flash the device.
1. [Install the Arduino Software for the ESP32-C3](https://www.airgradient.com/blog/install-arduino-c3-mini/)
2. [Install ESP8266 Arduino Core](https://arduino-esp8266.readthedocs.io/en/latest/installing.html#instructions)
3. [Install the required libraries](firmware/install-required-libraries.md)

## Flashing the Device
1. [Download the AirGradient sketch](firmware/ONE_I_9PSL.custom.ino)
2. Open the sketch in the Arduino IDE
3. Find and replace `YOUR_IP` with the IP address system that will host the exporter
4. Find and replace `YOUR_PORT` with the port the exporter will listen on (Default: 10000)
5. Select the correct board and port
6. Upload the sketch to the device
7. The device should now send data to the exporter

# Exporter Configuration
| Name                      | Description                                                                                                                   | Default  | Options                                | CLI Flag                         | Environment Variable                                |
|---------------------------|-------------------------------------------------------------------------------------------------------------------------------|----------|----------------------------------------|----------------------------------|-----------------------------------------------------|
| Address                   | The address the exporter will listen on                                                                                       | 0.0.0.0  | * or a valid IP                        | `--address`                      | `AIRGRADIENT_EXPORTER_ADDRESS`                      |
| Port                      | The port the exporter will listen on                                                                                          | 10000    | 1-65536                                | `--port`                         | `AIRGRADIENT_EXPORTER_PORT`                         |
| Path                      | The path the exporter will expose metrics on                                                                                  | /metrics |                                        | `--path`                         | `AIRGRADIENT_EXPORTER_PATH`                         |
| Log Level                 | The log level of the exporter                                                                                                 | info     | debug, info, warn, error, fatal, panic | `--log-level`                    | `LOG_LEVEL`                                         |
| Country                   | The country the device is located in. <br/><br/>Specifying `US` will cause the device to report in Fahrenheit                 | US       | Valid Country Code                     | `--ag-country`                   | `AIRGRADIENT_EXPORTER_AG_COUNTRY`                   |
| PM Standard               | The PM standard to use. If left empty, the US AQI is used.                                                                    | not set  | empty, ugm3                            | `--ag-pm-standard`               | `AIRGRADIENT_EXPORTER_AG_PM_STANDARD`               |
| CO2 Calibration Requested | Whether CO2 calibration is requested                                                                                          | false    | true, false                            | `--ag-co2-calibration-requested` | `AIRGRADIENT_EXPORTER_AG_CO2_CALIBRATION_REQUESTED` |
| LED Bar Test Requested    | Whether the LED bar test is requested                                                                                         | false    | true, false                            | `--ag-led-bar-test-requested`    | `AIRGRADIENT_EXPORTER_AG_LED_BAR_TEST_REQUESTED`    |
| LED Bar Mode              | The LED bar mode. If set to something other than `off` then the LED bar will represent the level of the measurement specified | off      | off, pm, co2                           | `--ag-led-bar-mode`              | `AIRGRADIENT_EXPORTER_AG_LED_BAR_MODE`              |
| Model                     | The model of the device                                                                                                       | not set  |                                        | `--ag-model`                     | `AIRGRADIENT_EXPORTER_AG_MODEL`                     |
| MQTT Broker URL           | The URL of the MQTT broker                                                                                                    | not set  |                                        | `--ag-mqtt-broker-url`           | `AIRGRADIENT_EXPORTER_MQTT_BROKER_URL`              |

# Inspiration
- [AirGradient](https://airgradient.com/)
- [Monitoring my home's air quality (CO2, PM2.5, Temp/Humidity) with AirGradient's DIY sensor](https://www.jeffgeerling.com/blog/2021/airgradient-diy-air-quality-monitor-co2-pm25)

# License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.