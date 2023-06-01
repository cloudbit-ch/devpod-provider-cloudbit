package options

import (
	"fmt"
	"os"
)

type Options struct {
	MachineID     string
	MachineFolder string

	Token    string
	Location string
	Image    string
	Product  string
	Network  string
}

func FromEnv(skipMachine bool) (*Options, error) {
	retOptions := &Options{}

	var err error
	if !skipMachine {
		if retOptions.MachineID, err = fromEnvOrError("MACHINE_ID"); err != nil {
			return nil, err
		}
		// prefix with devpod-
		retOptions.MachineID = "devpod-" + retOptions.MachineID

		if retOptions.MachineFolder, err = fromEnvOrError("MACHINE_FOLDER"); err != nil {
			return nil, err
		}
	}

	if retOptions.Token, err = fromEnvOrError("TOKEN"); err != nil {
		return nil, err
	}

	if retOptions.Location, err = fromEnvOrError("LOCATION"); err != nil {
		return nil, err
	}

	if retOptions.Image, err = fromEnvOrError("IMAGE"); err != nil {
		return nil, err
	}

	if retOptions.Product, err = fromEnvOrError("PRODUCT"); err != nil {
		return nil, err
	}

	retOptions.Network = os.Getenv("NETWORK")
	if retOptions.Network == "" {
		retOptions.Network = "Default Network"
	}

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("couldn't find option %s in environment, please make sure %s is defined", name, name)
	}

	return val, nil
}
