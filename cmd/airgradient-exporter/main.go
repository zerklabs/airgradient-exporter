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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	cliName        = "airgradient-exporter"
	cliDescription = "AirGradent Metric destination and exporter for Prometheus"
)

var (
	logger *zap.Logger
)

func NewRootCommand(ctx context.Context) (*cobra.Command, *viper.Viper) {
	const (
		envPrefix = "AIRGRADIENT_EXPORTER"
	)
	vcfg := viper.New()

	cmd := cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{cliName},
		Run:        runExporter(ctx, vcfg),
	}

	cmd.PersistentFlags().String("address", "0.0.0.0", "Address to listen on")
	_ = vcfg.BindPFlag("address", cmd.PersistentFlags().Lookup("address"))

	cmd.PersistentFlags().String("port", "10000", "Port to listen on")
	_ = vcfg.BindPFlag("port", cmd.PersistentFlags().Lookup("port"))

	cmd.PersistentFlags().String("path", "/metrics", "HTTP path to expose metrics on")
	_ = vcfg.BindPFlag("path", cmd.PersistentFlags().Lookup("path"))

	cmd.PersistentFlags().String("ag-country", "US", "AirGradient config Country code")
	_ = vcfg.BindPFlag("ag_country", cmd.PersistentFlags().Lookup("ag-country"))

	cmd.PersistentFlags().String("ag-pm-standard", "", "AirGradient config PM standard. If left empty, the device will use US AQI value. Set to 'ugm3' for µg/m³.")
	_ = vcfg.BindPFlag("ag_pm_standard", cmd.PersistentFlags().Lookup("ag-pm-standard"))

	cmd.PersistentFlags().Bool("ag-co2-calibration-requested", false, "AirGradient config CO2 calibration requested")
	_ = vcfg.BindPFlag("ag_co2_calibration_requested", cmd.PersistentFlags().Lookup("ag-co2-calibration-requested"))

	cmd.PersistentFlags().Bool("ag-led-bar-test-requested", false, "AirGradient config LED bar test requested")
	_ = vcfg.BindPFlag("ag_led_bar_test_requested", cmd.PersistentFlags().Lookup("ag-led-bar-test-requested"))

	cmd.PersistentFlags().String("ag-led-bar-mode", "off", "AirGradient config LED bar mode. Can be set to 'off', 'co2', 'pm'")
	_ = vcfg.BindPFlag("ag_led_bar_mode", cmd.PersistentFlags().Lookup("ag-led-bar-mode"))

	cmd.PersistentFlags().String("ag-model", "", "AirGradient config Device model")
	_ = vcfg.BindPFlag("ag_model", cmd.PersistentFlags().Lookup("ag-model"))

	cmd.PersistentFlags().String("ag-mqtt-broker-url", "", "AirGradient config MQTT broker address")
	_ = vcfg.BindPFlag("ag_mqtt_broker_url", cmd.PersistentFlags().Lookup("ag-mqtt-broker-url"))

	return &cmd, vcfg
}

func init() {
	var err error
	// create default logger if user does not supply one.
	config := zap.NewProductionConfig()
	// set default time formatter to "2006-01-02T15:04:05.000Z0700"
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		config.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		config.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		config.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		config.Level.SetLevel(zapcore.ErrorLevel)
	case "fatal":
		config.Level.SetLevel(zapcore.FatalLevel)
	case "panic":
		config.Level.SetLevel(zapcore.PanicLevel)
	default:
		config.Level.SetLevel(zapcore.InfoLevel)
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	logger = logger.With(zap.String("app", cliName))
}

func Start() {
	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	cmd, _ := NewRootCommand(ctx)
	if err := cmd.Execute(); err != nil {
		logger.Error("error executing command", zap.Error(err))
		return
	}
}

func main() {
	Start()
}
