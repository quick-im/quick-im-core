package messaging

type Message interface {
	Publish(subject string, msg []byte)
	Subscribe(subject string, callback func(msg []byte))
	SubscribeGroup(subject string, callback func(msg []byte))
}
