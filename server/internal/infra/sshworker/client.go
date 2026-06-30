package sshworker

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	Host     string
	User     string
	KeyPath  string
	BasePath string
}

func New(host, user, keyPath, basePath string) *Client {
	return &Client{
		Host:     host,
		User:     user,
		KeyPath:  keyPath,
		BasePath: basePath,
	}
}

func (c *Client) connect() (*ssh.Client, error) {
	key, err := os.ReadFile(c.KeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", c.Host+":22", config)
}

func (c *Client) UploadFile(localPath, remotePath string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	src, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (c *Client) RunCommand(cmd string) (string, error) {
	conn, err := c.connect()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out
	session.Stderr = &out

	err = session.Run(cmd)
	return out.String(), err
}

func (c *Client) SendTaskAndExecute(jsonData []byte) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	remoteFile := c.BasePath + "/mission_test.json"

	// 1. создаём файл на удалённой машине
	writeCmd := fmt.Sprintf(
		"cat > %s << 'EOF'\n%s\nEOF\n",
		remoteFile,
		string(jsonData),
	)

	// 2. запускаем python после записи
	runCmd := fmt.Sprintf(
		"cd %s && ./venv/bin/python -u drone.py /dev/ttyACM0 0.45",
		c.BasePath,
	)

	fullCmd := writeCmd + runCmd
	fmt.Println(fullCmd)
	var out bytes.Buffer
	session.Stdout = &out
	session.Stderr = &out

	err = session.Run(fullCmd)
	if err != nil {
		return fmt.Errorf("%v | output: %s", err, out.String())
	}

	return nil
}
