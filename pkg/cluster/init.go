package cluster

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

func Init(ccfg Config) error {
	var err error

	err = initToolsServer(ccfg)
	if err != nil {
		return fmt.Errorf("couldn't initiate toolsServer: %s", err)
	}

	return nil
}

func initToolsServer(ccfg Config) error {
	cToolsServer := ccfg.Cluster.Nodes.ToolsServer
	if reflect.ValueOf(cToolsServer).IsZero() {
		log.Println("skip toolsServer initialisation, given empty value(s)")
		return nil
	}

	if reflect.ValueOf(cToolsServer.ProviderMachineTemplate).IsZero() {
		msg := "for the toolsServer you need to provide a concrete ProviderMachineTemplate"
		return errors.New(msg)
	}

	v := reflect.ValueOf(cToolsServer.ProviderMachineTemplate)
	countNonEmpty := 0
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			countNonEmpty++
		}
	}

	if countNonEmpty > 1 {
		return errors.New("for the toolsServer you provided more than 1 ProviderMachineTemplate")
	}

	provider, err := extractFirstFound(v)
	if err != nil {
		return fmt.Errorf("couldn't extract ProviderMachineTemplate: %s", err)
	}

	switch provider.Interface().(type) {
	case hetznerCloudMachineProvider:
		//todo: continue
		fmt.Println("lets use the hetznerCloudMachineProvider")
	}

	return nil
}

func extractFirstFound(v reflect.Value) (reflect.Value, error) {
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			return v.Field(i), nil
		}
	}

	return reflect.Value{}, errors.New("coudn't extract/found even one")
}

//plan
////gather info
////create plan
//applay
////exec plan

//reconsile (currentState, desiredState) error
//// <- method or gets hetzner client

// interface cloud provide
// cloudProvider.reconcile(desiredState)

// how to notice vms ?
