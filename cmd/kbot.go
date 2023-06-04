/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	//	"github.com/stianeikeland/go-rpio"
	"github.com/hirosassa/zerodriver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	telebot "gopkg.in/telebot.v3"
)

var (
	// TeleToken bot
	TokenFile   = os.Getenv("TOKEN_FILE")
	TraceHost   = os.Getenv("TRACE_HOST")
	MetricsHost = os.Getenv("METRICS_HOST")
	// Load Telegram token from file
	tokenBytes, _ = ioutil.ReadFile(TokenFile)

	TeleToken = string(tokenBytes)
)

// Initialize OpenTelemetry
func initMetrics(ctx context.Context) {

	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)

	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(appVersion),
		attribute.String("some-attribute", "some-value"),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
		),
	)
	otel.SetMeterProvider(mp)

}

// Initialize OpenTelemetry
func initTracer(ctx context.Context) {

	// Set up OTLP exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(TraceHost), otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	// Set up tracer provider with OTLP exporter and resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appVersion),
			attribute.String("some-attribute", "some-value"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerodriver.NewProductionLogger()

		fmt.Printf("kbot %s started", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Plaese check TOKEN_FILE. %s", err)
			return
		}

		// err = rpio.Open()
		// if err != nil {
		// 	log.Printf("Unable to open gpio: %s", err.Error())
		// }

		// defer rpio.Close()

		trafficSignal := make(map[string]map[string]int8)

		trafficSignal["red"] = make(map[string]int8)
		trafficSignal["amber"] = make(map[string]int8)
		trafficSignal["green"] = make(map[string]int8)

		trafficSignal["red"]["pin"] = 12
		//default on/off
		//trafficSignal["red"]["on"]=0
		trafficSignal["amber"]["pin"] = 27
		trafficSignal["green"]["pin"] = 22

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {

			// var (
			// 	err error
			// 	pin = rpio.Pin(0)
			// )
			log.Print(m.Message().Payload, m.Text())
			payload := m.Message().Payload

			switch payload {
			case "hello":
				err = m.Send(fmt.Sprintf("Hello I'm Kbot %s!", appVersion))

			case "red", "amber", "green":

				// Start a new span
				ctx := context.Background()
				tracer := otel.Tracer("kbot")
				ctx, span := tracer.Start(
					ctx,
					"OnText",
					trace.WithAttributes(attribute.String("component", "kbot")),
					trace.WithAttributes(attribute.String("TraceID", trace.TraceID{1, 2, 3, 4}.String())),
				)
				defer span.End()

				trace_id := span.SpanContext().TraceID().String()
				//span_id := span.SpanContext().SpanID().String()
				logger.Info().Str("TraceID", trace_id).Msg(payload)

				meter := otel.GetMeterProvider().Meter("example")
				counter, _ := meter.Int64Counter("telebot_OnText")
				// Bind the counter to some labels

				counter.Add(ctx, 1)
				// pin = rpio.Pin(trafficSignal[payload]["pin"])
				if trafficSignal[payload]["on"] == 0 {
					// pin.Output()
					trafficSignal[payload]["on"] = 1
				} else {
					// pin.Input()
					trafficSignal[payload]["on"] = 0
				}

				err = m.Send(fmt.Sprintf("Switch %s light signal to %d", payload, trafficSignal[payload]["on"]))

			default:
				err = m.Send("Usage: /s red|amber|green")

			}

			return err

		})

		kbot.Start()
	},
}

func init() {
	ctx := context.Background()
	initTracer(ctx)
	initMetrics(ctx)
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Initialize OpenTelemetry tracer

}
