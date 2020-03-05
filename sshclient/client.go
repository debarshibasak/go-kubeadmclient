package sshclient

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type SshConnection struct {
	Timeout     time.Duration
	Username    string
	IP          string
	KeyLocation string
	VerboseMode bool
	ClientID    string
}

func (sh *SshConnection) Collect(cmd string) (string, error) {
	var signer ssh.Signer
	var config *ssh.ClientConfig

	timeout := sh.Timeout

	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	if sh.KeyLocation != "" {
		d, err := ioutil.ReadFile(sh.KeyLocation)
		if err != nil {
			return "", err
		}

		signer, err = ssh.ParsePrivateKey(d)
		if err != nil {
			return "", err
		}

		config = &ssh.ClientConfig{
			User:    sh.Username,
			Timeout: timeout,
			Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	} else {
		config = &ssh.ClientConfig{
			User:    sh.Username,
			Timeout: timeout,
			Auth:    []ssh.AuthMethod{ssh.PublicKeys()},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", sh.IP), config)
	if err != nil {
		return "", err
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	writeCloudInfoOut, err := session.Output(fmt.Sprintf("sh -c '%v'", cmd))
	if err != nil {
		return "", err
	}

	return string(writeCloudInfoOut), nil

}

func (sh *SshConnection) ScpToWithData(data []byte, destination string) error {

	s := "/tmp/" + uuid.New().String()
	err := ioutil.WriteFile(s, data, os.FileMode(0777))
	if err != nil {
		return err
	}

	return sh.ScpTo(s, destination)
}

func (sh *SshConnection) ScpFrom(source string, destination string) error {
	cmd := exec.Command("sh", "-c", "scp -i "+sh.KeyLocation+" "+source+" "+sh.Username+"@"+sh.IP+":"+destination)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}

	if sh.VerboseMode {
		fmt.Println("[client_id: "+sh.ClientID+"]"+string(out))
	}

	return err
}

func (sh *SshConnection) ScpTo(source string, destination string) error {

	c := "scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i " + sh.KeyLocation + " " + source + " " + sh.Username + "@" + sh.IP + ":" + destination
	fmt.Println(c)
	cmd := exec.Command("sh", "-c", c)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}

	if sh.VerboseMode {
		fmt.Println("[client_id: "+sh.ClientID+"]"+string(out))
	}

	return err
}

func (sh *SshConnection) Run(cmd []string) error {

	var signer ssh.Signer
	var config *ssh.ClientConfig

	timeout := sh.Timeout

	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	if sh.KeyLocation != "" {
		d, err := ioutil.ReadFile(sh.KeyLocation)
		if err != nil {
			return err
		}

		signer, err = ssh.ParsePrivateKey(d)
		if err != nil {
			return err
		}

		config = &ssh.ClientConfig{
			User:    sh.Username,
			Timeout: timeout,
			Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	} else {
		config = &ssh.ClientConfig{
			User:    sh.Username,
			Timeout: timeout,
			Auth:    []ssh.AuthMethod{ssh.PublicKeys()},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	}

	client, err := ssh.Dial("tcp", sh.IP+":22", config)
	if err != nil {
		return err
	}

	defer client.Close()

	for _, ln := range cmd {

		session, err := client.NewSession()
		if err != nil {
			session.Close()
			return err
		}

		fmt.Println(ln)

		writeCloudInfoOut, err := session.Output(fmt.Sprintf("sh -c '%v'", ln))
		if err != nil {
			session.Close()
			log.Println(string(writeCloudInfoOut))
			return err
		}

		if sh.VerboseMode {
			fmt.Println("[client_id: "+sh.ClientID+"]"+string(writeCloudInfoOut))
		}

		session.Close()
	}

	return nil
}
