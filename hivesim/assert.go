package hivesim

import "fmt"

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

func (t *T) Nil(err error, msgAndArgs ...interface{}) {
	if err != nil {
		t.Fatalf("%s, err: %v", messageFromMsgAndArgs(msgAndArgs), err)
	}
}

func (t *T) True(ok bool, msgAndArgs ...interface{}) {
	if !ok {
		t.Fatal(messageFromMsgAndArgs(msgAndArgs))
	}
}
