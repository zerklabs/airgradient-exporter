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

type Measurement struct {
	WiFiRssi int `json:"wifi"`
	// Relative CO2 level
	Rco2 int `json:"rco2"`
	// Ultrafine particles
	Pm01 int `json:"pm01"`
	// Fine particles
	Pm02 int `json:"pm02"`
	// Coarse particles
	Pm10 int `json:"pm10"`
	// Number of particles with diameter > 0.3 µm in 100 cm³ of air
	Pm003Count int `json:"pm003_count"`
	// TVOC (Total Volatile Organic Compounds) index
	TvocIndex int `json:"tvoc_index"`
	// NOx index (Nitrogen Oxide and dioxide)
	NoxIndex int `json:"nox_index"`
	// Ambient temperature
	Atmp float32 `json:"atmp"`
	// Relative humidity
	Rhum int `json:"rhum"`
	// Number of loop iterations since boot
	Boot int64 `json:"boot"`
}
