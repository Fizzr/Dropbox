package d7024e

type Net interface {
	SendPingMessage(contact *Contact)
	SendFindContactMessage(contact *Contact) CloseContacts
	SendFindDataMessage(hash string)
	SendStoreMessage(data []byte)
}

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// Query for k contacts closest to contact target
	// Should run synchronous (I guess)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
