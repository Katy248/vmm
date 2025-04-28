package qmp

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"vmm/vm"
)

type Connection struct {
	conn net.Conn
}

func New(socketFile string) (*Connection, error) {
	if _, err := os.Stat(socketFile); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	dial, err := net.Dial("unix", socketFile)
	if err != nil {
		return nil, err
	}
	connection := &Connection{conn: dial}

	connection.getGreetingMsg()

	return connection, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) read() ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf[:])
	if err != nil {
		return nil, err
	}
	// log.Println("Readed string: ", string(buf[0:n]))
	return buf[0:n], nil
}

type QmpGreetingMsg struct {
	QMP struct {
		Version struct {
			Qemu struct {
				Micro string `json:"micro"`
				Minor string `json:"minor"`
				Major string `json:"major"`
			} `json:"qemu"`
			Package string `json:"package"`
		} `json:"version"`
		Capabilities []string `json:"capabilities"`
	} `json:"QMP"`
}
type VMStatus struct {
	Status     string `json:"status"`
	Running    bool   `json:"running"`
	Singlestep bool   `json:"singlestep"`
}
type QmpVMStatus struct {
	Return VMStatus `json:"return"`
}

func (c *Connection) getGreetingMsg() (*QmpGreetingMsg, error) {
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	var msg QmpGreetingMsg
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (c *Connection) SendQmpCapabilities() error {
	command, err := json.Marshal(map[string]string{"execute": "qmp_capabilities"})
	if err != nil {
		return err
	}
	_, err = c.conn.Write(command)
	if err != nil {
		return err
	}
	c.read()
	return nil
}

func (c *Connection) GetVMStatus(info vm.VM) (*VMStatus, error) {
	command, err := json.Marshal(map[string]string{"execute": "query-status"})
	if err != nil {
		return nil, err
	}
	_, err = c.conn.Write(command)
	if err != nil {
		return nil, err
	}

	response, err := c.read()
	if err != nil {
		return nil, err
	}
	var status QmpVMStatus
	err = json.Unmarshal(response, &status)
	if err != nil {
		return nil, err
	}
	return &status.Return, nil
}

func (c *Connection) SendStop() error {
	command, err := json.Marshal(map[string]string{"execute": "quit"})
	if err != nil {
		return err
	}
	_, err = c.conn.Write(command)
	if err != nil {
		return err
	}

	c.read()
	return nil
}
