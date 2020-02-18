# Flaming Steve

> Because why not?

## TODO

- [x] Events bus scalfolding
- [x] Events bus autodiscovery (zeroconf)
- [x] Sensor: ak9753 hardware support
- [ ] Sensor: persistent state (during restart)
- [ ] Sensor: support other sensors model
- [ ] Sensor: unique id generator (base on MAC)
- [x] Sensor UI: display of sensor data
- [x] Sensor UI: auto add/remove sensors
- [ ] Sensor UI: Remote configuration
- [x] Sensor UI: Generic configs ui
- [x] Display: simulator
- [ ] Display: flame panel
- [ ] muthur: JSON-RPC for react based frontend 
- [ ] Seq: simple display sequencer
- [ ] Game Logic
- [ ] docker compose
- [ ] deamonise processes

## Sensor: `cmd/sensor`

Read outs the data from a [ak9753](http://wiki.seeedstudio.com/Grove-Human_Presence_Sensor-AK9753/) sensor.

```text
Usage of sensor:
  -i, --interval duration    interval for IR evaluation (default 5ms)
      --nats-server string   publish nats server where to push the sensor data
  -p, --publish              url for publish data push
      --ui                   display real time information on the terminal
```

## Sensor UI: `cmd/sensui`

Connects to muthur ([nats server](https://docs.nats.io/), display sensor data and set remote configuration.

## Sensor Simulator: `cmd/sensim`

Simulate one or more sensors.

## MUTHUR: `cmd/muthur`

Central service to be used for service (dis/re)covery. Also will tell you human are 
expendable if a xenomorph is present on your spaceship. 

Pretty much a embedded nats messaging server with a zeroconf service for discovery.

## Sequencer: `cmd/seq`

Small program which use a novation launchpad to display pixel's sequence

## Display Simulator: `cmd/dispsim`

Simulate a display device (flame panel).

## Used techs

- [periph.io](https://periph.io/project/library/)
- [nats.io](https://docs.nats.io/)
- [nucular ui](https://github.com/aarzilli/nucular)
