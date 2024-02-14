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

package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-json-experiment/json"
	promClient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zerklabs/airgradient-exporter/pkg/ag"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	atmp *promClient.GaugeVec
	rco2 *promClient.GaugeVec
	rhum *promClient.GaugeVec
	pm01 *promClient.GaugeVec
	pm25 *promClient.GaugeVec
	pm10 *promClient.GaugeVec
	pm03 *promClient.GaugeVec
	tvoc *promClient.GaugeVec
	nox  *promClient.GaugeVec
	boot *promClient.CounterVec
	rssi *promClient.GaugeVec
)

var shutdownTimeout = time.Second * 5

func setupMetrics(registry *promClient.Registry) {
	atmp = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "atmp",
		Help:      "Ambient temperature",
		Namespace: "airgradient",
	}, []string{"mac"})

	rco2 = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "rco2",
		Help:      "Relative CO2 level",
		Namespace: "airgradient",
	}, []string{"mac"})

	rhum = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "rhum",
		Help:      "Relative humidity",
		Namespace: "airgradient",
	}, []string{"mac"})

	pm01 = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "pm01",
		Help:      "Ultrafine particles",
		Namespace: "airgradient",
	}, []string{"mac"})

	pm25 = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "pm25",
		Help:      "Fine particles",
		Namespace: "airgradient",
	}, []string{"mac"})

	pm10 = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "pm10",
		Help:      "Coarse particles",
		Namespace: "airgradient",
	}, []string{"mac"})

	pm03 = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "pm03",
		Help:      "Number of particles with diameter > 0.3 µm in 100 cm³ of air",
		Namespace: "airgradient",
	}, []string{"mac"})

	tvoc = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "tvoc",
		Help:      "TVOC (Total Volatile Organic Compounds) index",
		Namespace: "airgradient",
	}, []string{"mac"})

	nox = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "nox",
		Help:      "NOx index (Nitrogen Oxide and dioxide)",
		Namespace: "airgradient",
	}, []string{"mac"})

	boot = promauto.NewCounterVec(promClient.CounterOpts{
		Name:      "boot",
		Help:      "Number of loop iterations since boot",
		Namespace: "airgradient",
	}, []string{"mac"})

	rssi = promauto.NewGaugeVec(promClient.GaugeOpts{
		Name:      "rssi",
		Help:      "WiFi RSSI",
		Namespace: "airgradient",
	}, []string{"mac"})

	registry.MustRegister(atmp, rco2, rhum, pm01, pm25, pm10, pm03, tvoc, nox, boot, rssi)
}

func runExporter(ctx context.Context, vcfg *viper.Viper) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		logger.Info(
			"starting exporter",
			zap.String("address", vcfg.GetString("address")),
			zap.String("port", vcfg.GetString("port")),
			zap.String("path", vcfg.GetString("path")),
		)

		// Set up the prometheus metrics registry
		registry := promClient.NewPedanticRegistry()

		setupMetrics(registry)

		// Set up the prometheus metrics handler
		metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

		// Set up the HTTP server mux
		mux := http.NewServeMux()
		mux.Handle("GET "+vcfg.GetString("path"), metricsHandler)
		mux.Handle("GET /sensors/{mac}/one/config", agConfigHandler(vcfg))
		mux.Handle("POST /sensors/{mac}/measures", agHandleMeasurements())

		server := &http.Server{
			Addr:    vcfg.GetString("address") + ":" + vcfg.GetString("port"),
			Handler: mux,
		}

		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			err := server.ListenAndServe()
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
			return nil
		})
		g.Go(func() error {
			// Wait until g ctx canceled, then try to shut down server.
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			return server.Shutdown(ctx)
		})

		err := g.Wait()
		if err != nil {
			logger.Error("error running exporter", zap.Error(err))
		}
	}
}

func agConfigHandler(vcfg *viper.Viper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.String("content-type", r.Header.Get("Content-Type")),
			zap.String("accept", r.Header.Get("Accept")),
		)

		w.Header().Set("Content-Type", "application/json")
		config := ag.DefaultServerConfig()
		config.Country = vcfg.GetString("ag_country")
		config.PMStandard = vcfg.GetString("ag_pm_standard")
		config.CO2CalibrationRequested = vcfg.GetBool("ag_co2_calibration_requested")
		config.LEDBarTestRequested = vcfg.GetBool("ag_led_bar_test_requested")
		config.LEDBarMode = vcfg.GetString("ag_led_bar_mode")
		config.Model = vcfg.GetString("ag_model")
		config.MQTTBrokerURL = vcfg.GetString("ag_mqtt_broker_url")

		b, err := json.Marshal(config, json.DefaultOptionsV2())
		if err != nil {
			logger.Error("failed to marshal config", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			logger.Error("failed to write response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func agHandleMeasurements() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.String("content-type", r.Header.Get("Content-Type")),
			zap.String("accept", r.Header.Get("Accept")),
		)

		var m ag.Measurement
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &m, json.DefaultOptionsV2())
		if err != nil {
			logger.Error(
				"failed to unmarshal measurement",
				zap.Error(err),
				zap.String("body", string(body)),
			)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mac := strings.TrimPrefix(r.PathValue("mac"), "airgradient:")

		atmp.WithLabelValues(mac).Set(float64(m.Atmp))
		rco2.WithLabelValues(mac).Set(float64(m.Rco2))
		rhum.WithLabelValues(mac).Set(float64(m.Rhum))
		pm01.WithLabelValues(mac).Set(float64(m.Pm01))
		pm25.WithLabelValues(mac).Set(float64(m.Pm02))
		pm10.WithLabelValues(mac).Set(float64(m.Pm10))
		pm03.WithLabelValues(mac).Set(float64(m.Pm003Count))
		tvoc.WithLabelValues(mac).Set(float64(m.TvocIndex))
		nox.WithLabelValues(mac).Set(float64(m.NoxIndex))
		boot.WithLabelValues(mac).Add(float64(m.Boot))
		rssi.WithLabelValues(mac).Set(float64(m.WiFiRssi))

		if ce := logger.Check(zap.DebugLevel, "received measurement"); ce != nil {
			ce.Write(
				zap.String("mac", mac),
				zap.Any("measurement", m),
			)
		} else {
			logger.Info("received measurement", zap.String("mac", mac))
		}

		_, err = w.Write([]byte("OK"))
		if err != nil {
			logger.Error("failed to write response", zap.Error(err))
		}
	}
}
