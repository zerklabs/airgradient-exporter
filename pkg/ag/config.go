/*
 * MIT License
 *
 * Copyright (c)  2024 Robin Harper.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package ag

type ServerConfig struct {
	// Country code
	//
	// If set to "US" the device will use Fahrenheit
	Country string `json:"country"`

	// PM standard to use for the measurements
	//
	// If set to ugm3 the measurements will be in µg/m³ which will
	// set inUSAQI in the device to false. Any other value will set
	// inUSAQI to true.
	PMStandard string `json:"pmStandard"`

	// If set to true, the device will request a CO2 calibration
	CO2CalibrationRequested bool `json:"co2CalibrationRequested"`

	// If set to true, the device will request a LED bar test
	LEDBarTestRequested bool `json:"ledBarTestRequested"`

	// LED bar mode
	//
	// Can be set to "off", "co2", "pm"
	// If set to PM the LED bar will be relative to the PM level
	// If set to CO2 the LED bar will be relative to the CO2 level
	LEDBarMode string `json:"ledBarMode"`

	// Device model
	Model string `json:"model,omitempty"`

	// MQTT broker address
	MQTTBrokerURL string `json:"mqttBrokerUrl,omitempty"`
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Country:                 "US",
		PMStandard:              "",
		CO2CalibrationRequested: false,
		LEDBarTestRequested:     false,
		LEDBarMode:              "off",
	}
}
