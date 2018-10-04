package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akaumov/cube-executor"
	"github.com/akaumov/cube-websocket-gateway"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Action = runServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "bus-host",
			EnvVar: "GATEWAY_BUS_HOST",
			Usage:  "bus host",
		},
		cli.IntFlag{
			Name:   "bus-port",
			EnvVar: "GATEWAY_BUS_PORT",
			Usage:  "bus port",
		},
		cli.StringFlag{
			Name:   "jwt-secret",
			EnvVar: "GATEWAY_JWT_SECRET",
			Usage:  "jwt secret",
		},
		cli.IntFlag{
			Name:   "max-connections",
			EnvVar: "GATEWAY_MAX_CONNECTIONS",
			Usage:  "maximum number of connections",
		},
		cli.StringFlag{
			Name:   "endpoints-map",
			EnvVar: "GATEWAY_ENDPOINTS_MAP",
			Usage:  "map url to endpoint",
		},
		cli.BoolTFlag{
			Name:   "only-authorized-requests",
			EnvVar: "GATEWAY_ONLY_AUTHORIZED_REQUESTS",
			Usage:  "handle only authorized requests",
		},
		cli.BoolFlag{
			Name:   "dev",
			EnvVar: "GATEWAY_DEV",
			Usage:  "log all requests",
		},
		cli.StringFlag{
			Name:   "port",
			EnvVar: "GATEWAY_PORT",
			Usage:  "port to listen",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func parseChannelsMap(rawMap string) (*map[cube_executor.CubeChannel]cube_executor.BusChannel, error) {
	if rawMap == "" {
		return &map[cube_executor.CubeChannel]cube_executor.BusChannel{}, nil
	}

	params := map[cube_executor.CubeChannel]cube_executor.BusChannel{}

	for _, rawMap := range strings.Split(rawMap, ";") {
		splittedMap := strings.Split(rawMap, ":")

		if len(splittedMap) != 2 {
			return nil, fmt.Errorf("Wrong params format: %v\n", rawMap)
		}

		key := splittedMap[0]
		value := splittedMap[1]

		params[cube_executor.CubeChannel(key)] = cube_executor.BusChannel(value)
	}

	return &params, nil
}

func runServer(c *cli.Context) error {

	busHost := c.String("bus-host")
	if busHost == "" {
		return fmt.Errorf("bus host is required")
	}

	busPort := c.Int("bus-port")
	if busPort == 0 {
		return fmt.Errorf("bus port is required")
	}

	jwtSecret := c.String("jwt-secret")
	if jwtSecret == "" {
		return fmt.Errorf("jwt secret is required")
	}

	maxConnections := c.String("max-connections")
	if maxConnections == "" {
		return fmt.Errorf("max connections is required")
	}

	port := c.String("port")

	endpointsMap := c.String("endpoints-map")
	channelsMapping, err := parseChannelsMap(endpointsMap)
	if err != nil {
		return fmt.Errorf("wrong channels mapping")
	}

	onlyAuthorizedRequests := "true"
	if c.Bool("only-authorized-requests") {
		onlyAuthorizedRequests = "true"
	} else {
		onlyAuthorizedRequests = "false"
	}

	dev := "false"
	if c.Bool("dev") {
		dev = "true"
	} else {
		dev = "false"
	}

	cube, err := cube_executor.NewCube(cube_executor.CubeConfig{
		BusPort: busPort,
		BusHost: busHost,
		Params: map[string]string{
			"jwtSecret":              jwtSecret,
			"maxConnections":         maxConnections,
			"endpointsMap":           endpointsMap,
			"onlyAuthorizedRequests": onlyAuthorizedRequests,
			"dev":                    dev,
			"port":                   port,
		},
		ChannelsMapping: *channelsMapping,
	}, &cube_websocket_gateway.Handler{})

	if err != nil {
		return fmt.Errorf("can't start: %v", err)
	}

	return cube.Start()
}
