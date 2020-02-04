# Flaming Steve

> Because why not?

## `cmd/sensor`

Read outs the data from a [ak9753](http://wiki.seeedstudio.com/Grove-Human_Presence_Sensor-AK9753/) sensor.

```text
Usage of sensor:
  -i, --interval duration    interval for IR evaluration (default 30ms)
      --nats-server string   publish nats server where to push the sensor data
  -p, --publish              url for publish data push
      --remote               connect to a remote sensor
  -s, --smoothing float32    0.3 very steep, 0.1 less steep, 0.05 less steep (default 0.05)
  -t, --threshold float32    presence threshold (default 10)
      --ui                   display real time informatio on the terminal
```

## `cmd/sensor-ui`

Connects to a [nats server](https://docs.nats.io/) and display sensor data.

```text
Usage TBD
```

## `cmd/muthur`

Central service to be used for service (dis/re)covery. Also will tell you human are 
expendable if a xenomorph is present on your spaceship. 

Pretty much a embedded nats messaging server with a zeronconf service for discovery.

## Used techs

- [periph.io](https://periph.io/project/library/)
- [nats.io](https://docs.nats.io/)
- [nucular ui](https://github.com/aarzilli/nucular)
