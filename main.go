package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/dineshsuryaprakash/go-rest-api/api"
	"github.com/dineshsuryaprakash/go-rest-api/config"
	"k8s.io/klog/v2"
)

func main() {
	flag.Parse()

	//create Api server
	HolidayRestApiServer, err := api.NewServer(config.ListenAddress, config.BankHolidayJsonEndpointURL)
	if err != nil {
		klog.Fatalf("Error creating HolidayRestApiServer: %v", err)
		return
	}

	//init and run Api server
	HolidayRestApiServer.Initialize()
	HolidayRestApiServer.Run()

	// Listen for OS signals to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	err = HolidayRestApiServer.Shutdown()
	if err != nil {
		klog.Fatalf("Error shutting down HolidayRestApiServer: %v", err)
		return
	}

}

func init() {
	klog.InitFlags(nil)
	config.InitFlags()
}
