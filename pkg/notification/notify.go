package notification

import "sync"

type Notifier interface {
	Subscribe(channel chan<- bool)
	Unsubscribe(channel chan<- bool)
	//UnsubscribeAll()
}

type NotifierImpl struct {
	suscribers []chan<- bool
	mutex      sync.Mutex
}

func (n *NotifierImpl) UnsubscribeAll() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, c := range n.suscribers {
		close(c)
	}
}

func (n *NotifierImpl) Subscribe(channel chan<- bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.suscribers = append(n.suscribers, channel)
}

func (n *NotifierImpl) Unsubscribe(channel chan<- bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for idx, ch := range n.suscribers {
		if ch == channel {
			if len(n.suscribers) > 1 {
				n.suscribers = append(n.suscribers[:idx], n.suscribers[idx+1:]...)
			} else {
				n.suscribers = []chan<- bool{}
			}
			return
		}
	}
}

func (n *NotifierImpl) Notify() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, c := range n.suscribers {
		select {
		case c <- true:
		default:
			println("skipping notification")
		}
	}
}
