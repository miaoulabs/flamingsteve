package notify

import "sync"

type Notifier struct {
	suscribers []chan<- bool

	mutex sync.Mutex
}

func (n *Notifier) UnsubscribeAll() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, c := range n.suscribers {
		close(c)
	}
}

func (n *Notifier) Subscribe(channel chan<- bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.suscribers = append(n.suscribers, channel)
}

func (n *Notifier) Notify() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, c := range n.suscribers {
		select {
		case c <- true:
		default:
		}
	}
}
