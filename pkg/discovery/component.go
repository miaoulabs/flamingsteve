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
}

type IdentifierConfig struct {
	Name  string
	Model string
	Type  EntryType
}

func NewIdentifier(cfg IdentifierConfig, log logger.Logger) *Component {
	i := &Component{
		cfg: cfg,
		log: logger.Dummy(),
	}

	if log != nil {
		i.log = log
	}

	return i
}

func (c *Component) Connect() {
	var err error

	c.log.Infof("subscribing to topic '%s'", TopicScan)
	c.sub, err = muthur.Connection().Subscribe(TopicScan, c.whoAmI)

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
	err = muthur.Connection().Publish(topic, &e)
	if err != nil {
		c.log.Errorf("error sending message on topic '%s'", topic)
	}
	_ = c.sub.Unsubscribe()
}

func (c *Component) PushData(data interface{}) error {
	return muthur.Connection().Publish(c.dataTopic(), data)
}

func (c *Component) entry() Entry {
	hostname, _ := os.Hostname()

	return Entry{
		Type:        c.cfg.Type,
		Name:        c.cfg.Name,
		Hostname:    hostname,
		Model:       c.cfg.Model,
		Id:          c.Id(),
		DataTopic:   c.dataTopic(),
		ConfigTopic: c.configTopic(),
	}
}

func (c *Component) whoAmI(*nats.Msg) {
	topic := fmt.Sprintf("%s.%s.%s", TopicDeviceOn, c.cfg.Type, c.Id())
	c.log.Debugf("responding to a who broadcast on topic '%s'", topic)

	e := c.entry()
	err := muthur.Connection().Publish(topic, &e)
	if err != nil {
		c.log.Errorf("error sending message on topic '%s'", err)
	}
}

func (c *Component) Id() string {
	return c.cfg.Name // TODO: generate a unique id based on mac addr
}

func (c *Component) dataTopic() string {
	return fmt.Sprintf("%s.%s.%s.data", c.cfg.Type, c.cfg.Model, c.Id())
}

func (c *Component) configTopic() string {
	return fmt.Sprintf("%s.%s.%s.configs", c.cfg.Type, c.cfg.Model, c.Id())
}
