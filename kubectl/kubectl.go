package kubectl

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Kubectl struct {
	kubeconfig []byte
	timeout time.Duration
}

func New(kubeconfig []byte) *Kubectl {
	return &Kubectl{
		kubeconfig:kubeconfig,
		timeout:20*time.Second,
	}
}

func (k *Kubectl) ApplyFile(file string) error {
	location := mount(k.kubeconfig)
	defer unmount(location)

	cmd := exec.Command("sh","-c", "kubectl apply -f "+file)
	cmd.Env = []string{
		"KUBECONFIG="+location,
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v - %v", string(out), err.Error())
	}

	return nil
}

func mount(kubeconfig []byte) string {
	filename := uuid.New().String()
	fileLoc := "/tmp/"+ filename
	if err := ioutil.WriteFile(fileLoc, kubeconfig, os.FileMode(0777)); err != nil {
		panic(err)
	}
	return fileLoc
}

func unmount(fileLoc string) {
	err := os.RemoveAll(fileLoc)
	if err != nil {
		panic(err)
	}
}

func (k *Kubectl) SetLabel(node string, annotation string) error {
	location := mount(k.kubeconfig)
	defer unmount(location)

	cmd := exec.Command("sh","-c", "kubectl label node "+node+" node-role.kubernetes.io/"+annotation+"="+annotation)
	cmd.Env = []string{
		"KUBECONFIG="+location,
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v - %v", string(out), err.Error())
	}

	return nil
}

func (k *Kubectl) TaintAllNodes(annotation string) error {

	location := mount(k.kubeconfig)
	defer unmount(location)

	cmd := exec.Command("sh","-c", "kubectl taint nodes --all "+annotation)
	cmd.Env = []string{
		"KUBECONFIG="+location,
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v - %v", string(out), err.Error())
	}

	return nil
}

