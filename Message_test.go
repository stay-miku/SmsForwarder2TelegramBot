package main

import "testing"

func TestMessageGenerate(t *testing.T) {
	t.Log(messageGenerate("1234567890", "Your sms verify code is 182734"))
}
