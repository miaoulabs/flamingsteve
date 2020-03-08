# Flaming Steve

> Because why not?

## TODO

- [x] Events bus scalfolding
- [x] Events bus autodiscovery (zeroconf)
- [x] Sensor: ak9753 hardware support
- [ ] Sensor: persistent state (during restart)
- [x] Sensor: support other sensors model
- [x] Sensor UI: display of sensor data
- [x] Sensor UI: auto add/remove sensors
- [ ] Sensor UI: Remote configuration
- [x] Sensor UI: Generic configs ui
- [x] Display: simulator
- [x] Display: pico matrix 5x5
- [ ] Display: flame panel (rs285)
- [ ] JSON/GRPC webservice for imperative API 
- [x] Seq: simple display sequencer
- [ ] Game Logic
- [ ] docker compose
- [ ] deamonise processes

## Sensor: `cmd/sensor`

Read outs the data from a [ak9753](http://wiki.seeedstudio.com/Grove-Human_Presence_Sensor-AK9753/) sensor.

```text
      --mean int            number of sample to use for mean (default 6)
  -n, --name string         sensor name used for discovery
      --no-presence         disable presence detector
      --orphan              don't try to connect to muthur
  -t, --threshold float32   presence threshold (default 100)
      --type SensorType     sensor model [ak9753, amg8833] (default None)
      --ui                  ak9753_display real time information on the terminal
```

## Sensor UI: `cmd/sensui`

Connects to muthur ([nats server](https://docs.nats.io/), display sensor data and set remote configuration.

## Sensor Simulator: `cmd/sensim`

Simulate one or more sensors.

## Glue: `cmd/glue`

A webservice used to query the current state of all sensors and display.

## MUTHUR: `cmd/muthur`

Central service to be used for service (dis/re)covery. Also will tell you human are 
expendable if a xenomorph is present on your spaceship. 

Pretty much a embedded nats messaging server with a zeroconf service for discovery.

## Sequencer: `cmd/seq`

Small program which use a novation launchpad to display pixel's sequence

## Senspad: `cmd/senspad`

Visualised data from an amg8833 ir matrix camera unto a 8x8 launchpad midi controller.

## Matrix Display: `cmd/dispmatrix`

Run a display on a 2 small pimoroni 5x5 led matrix. 

## Display Simulator: `cmd/dispsim`

Simulate a display device (flame panel).

## Used techs

- [periph.io](https://periph.io/project/library/)
- [nats.io](https://docs.nats.io/)
- [nucular ui](https://github.com/aarzilli/nucular)

