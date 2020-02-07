package discovery

import (
	"flamingsteve/pkg/logger"
	"flamingsteve/pkg/muthur"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

type Identifier struct {
	cfg IdentifierConfig
	sub *nats.Subscription
	log logger.Logger
}

type IdentifierConfig struct {
	Name  string
	Model string
	Type  EntryType
}

func NewIdentifier(cfg IdentifierConfig, log logger.Logger) *Identifier {
	i := &Identifier{
		cfg: cfg,
		log: logger.Dummy(),
	}

	if log != nil {
		i.log = log
	}

	return i
}

func (i *Identifier) Connect() {
	var err error

	i.log.Infof("subscribing to topic '%s'", TopicScan)
	i.sub, err = muthur.Connection().Subscribe(TopicScan, i.whoAmI)

	if err != nil {
		i.log.Errorf("error subscribing to topic '%s'", TopicScan)
	}

	// send a who am i to tell this object exists
	i.whoAmI(nil)
}

func (i *Identifier) Close() {
	var err error
	e := i.entry()
	topic := fmt.Sprintf("%s.%s.", TopicDeviceOff, i.Id())
	err = muthur.Connection().Publish(topic, &e)
	if err != nil {
		i.log.Errorf("error sending message on topic '%s'", topic)
	}
	_ = i.sub.Unsubscribe()
}

func (i *Identifier) PushRaw(data interface{}) error {
	return muthur.Connection().Publish(i.dataTopic(), data)
}

func (i *Identifier) PushPresence(data interface{}) error {
	//return muthur.Connection().Publish(i.Pre(), data)
	return nil
}

func (i *Identifier) entry() Entry {
	hostname, _ := os.Hostname()

	return Entry{
		Type:        i.cfg.Type,
		Name:        i.cfg.Name,
		Hostname:    hostname,
		Model:       i.cfg.Model,
		Id:          i.Id(),
		DataTopic:   i.dataTopic(),
		ConfigTopic: i.configTopic(),
	}
}

func (i *Identifier) whoAmI(*nats.Msg) {
	topic := fmt.Sprintf("%s.%s.%s", TopicDeviceOn, i.cfg.Type, i.Id())
	i.log.Infof("sending message on topic '%s'", topic)

	e := i.entry()
	err := muthur.Connection().Publish(topic, &e)
	if err != nil {
		i.log.Errorf("error sending message on topic '%s'", err)
	}
}

func (i *Identifier) Id() string {
	return i.cfg.Name // TODO: generate a unique id based on mac addr
}

func (i *Identifier) dataTopic() string {
	return fmt.Sprintf("%s.%s.%s.data", i.cfg.Type, i.cfg.Model, i.Id())
}

func (i *Identifier) configTopic() string {
	return fmt.Sprintf("%s.%s.%s.configs", i.cfg.Type, i.cfg.Model, i.Id())
}
