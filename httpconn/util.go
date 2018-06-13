package httpconn

var constFakeAddr = &fakeAddr{}

type fakeAddr struct{}

func (a *fakeAddr) Network() string {
	return "fake_address_network"
}

func (a *fakeAddr) String() string {
	return "fake_address"
}
