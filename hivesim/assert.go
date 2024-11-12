package hivesim

func (t *T) Nil(err error, msg string) {
	if err != nil {
		t.Fatalf("%s, err: %v", msg, err)
	}
}

func (t *T) True(ok bool, msg string) {
	if !ok {
		t.Fatalf("%s", msg)
	}
}
