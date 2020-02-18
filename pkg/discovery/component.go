package discovery

import (
	"flamingsteve/pkg/logger"
	"flamingsteve/pkg/muthur"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

type Component struct {
	cfg IdentifierConfig
	sub *nats.Subscription
	log logger.Logger

	cfgIn  *nats.Subscription
	cfgOut *nats.Subscription
}

type IdentifierConfig struct {
	Name  string
	Model string
	Type  EntryType
}

func NewComponent(cfg IdentifierConfig) *Component {
	i := &Component{
		cfg: cfg,
		log: logger.Dummy(),
	}

	i.log = logFactory(cfg.Name)

	return i
}

func (c *Component) Connect() {
	var err error

	c.log.Infof("subscribing to topic '%s'", TopicScan)
	c.sub, err = muthur.EncodedConnection().Subscribe(TopicScan, c.whoAmI)

	if err != nil {
		c.log.Errorf("error subscribing to topic '%s'", TopicScan)
	}

	// send a who am c to tell this object exists
	c.whoAmI(nil)
}

func (c *Component) Close() {
	var err error
	e := c.entry()
	topic := fmt.Sprintf("%s.%s.", TopicDeviceOff, c.Id())
	err = muthur.EncodedConnection().Publish(topic, &e)
	if err != nil {
		c.log.Errorf("error sending message on topic '%s'", topic)
	}
	_ = c.sub.Unsubscribe()

	if c.cfgOut != nil {
		_ = c.cfgOut.Unsubscribe()
	}

	if c.cfgIn != nil {
		_ = c.cfgIn.Unsubscribe()
	}
}

func (c *Component) PushData(data interface{}) error {
	return muthur.EncodedConnection().Publish(c.DataTopic(), data)
}

/*
	Register a callback when configuration are requested
*/
func (c *Component) OnConfigRequest(getConfig func() []byte) {
	c.cfgOut, _ = muthur.Connection().Subscribe(c.readConfigTopic(), func(msg *nats.Msg) {
		data := getConfig()
		err := msg.Respond(data)
		if err != nil {
			c.log.Errorf("%v", err)
		}
	})
}

/*
	Register a callback when configuration are written
*/
func (c *Component) OnConfigWrite(setConfig func([]byte)) {
	c.cfgIn, _ = muthur.Connection().Subscribe(c.writeConfigTopic(), func(msg *nats.Msg) {
		setConfig(msg.Data)
	})
}

func (c *Component) entry() Entry {
	hostname, _ := os.Hostname()

	return Entry{
		Type:        c.cfg.Type,
		Name:        c.cfg.Name,
		Hostname:    hostname,
		Model:       c.cfg.Model,
		Id:          c.Id(),
		DataTopic:   c.DataTopic(),
		ConfigTopic: c.readConfigTopic(),
	}
}

func (c *Component) whoAmI(*nats.Msg) {
	topic := fmt.Sprintf("%s.%s.%s", TopicDeviceOn, c.cfg.Type, c.Id())
	c.log.Debugf("responding to a who broadcast on topic '%s'", topic)

	e := c.entry()
	err := muthur.EncodedConnection().Publish(topic, &e)
	if err != nil {
		c.log.Errorf("error sending message on topic '%s'", err)
	}
}

func (c *Component) Id() string {
	return c.cfg.Name // TODO: generate a unique id based on mac addr
}

func (c *Component) DataTopic() string {
	return fmt.Sprintf("%s.%s.%s.data", c.cfg.Type, c.cfg.Model, c.Id())
}

func (c *Component) readConfigTopic() string {
	return c.configTopic() + ".read"
}

func (c *Component) writeConfigTopic() string {
	return c.configTopic() + ".write"
}

func (c *Component) configTopic() string {
	return fmt.Sprintf("%s.%s.%s.configs", c.cfg.Type, c.cfg.Model, c.Id())
}
