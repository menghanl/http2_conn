package httpconn

var constFakeAddr = &fakeAddr{}

type fakeAddr struct{}

func (a *fakeAddr) Network() string {
	return "fake_address_network"
}

func (a *fakeAddr) String() string {
	return "fake_address"
}

const magicHandshakeStr = "M_a_G_i_C\n"
