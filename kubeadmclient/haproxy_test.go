package kubeadmclient

import (
	"log"
	"testing"
)

func TestHaProxyNode_Install(t *testing.T) {
	t.SkipNow()

	var haproxy = NewHaProxyNode("ubuntu", "192.168.64.218", "/Users//.ssh/id_rsa")
	if err := haproxy.install([]string{"192.1.1.1", "192.11.1.1", "192.11.1.2"}); err != nil {
		log.Fatal(err)
	}
}
