package d2client

import (
	"errors"
	"fmt"
	"net"
)

var (
	// ErrUnableToWrite is used when we try to write on the connection and fail.
	ErrUnableToWrite = errors.New("Unable to write over the connection")

	// ErrNotConnected will be used when a write is being performed on a closed connection.
	ErrNotConnected = errors.New("Connection has not been opened")
)

// Client provides operations on the diablo game client.
type Client struct {
	connection net.Conn
}

// Open opens a connection to the given host.
func (c *Client) Open(host string) error {
	// Resolve the host.
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}

	// Dial the host for a connection.
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	// Set connection on the client.
	c.connection = conn

	return nil
}

// Close will close the socket if it's open.
func (c *Client) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}

// Login needs to run before you can actually send any message, will
// return an error if we can't connect to the server.
func (c *Client) Login(account string, password string) error {
	if c.connection == nil {
		return ErrNotConnected
	}

	// Write First carriage return.
	_, err := c.connection.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	// Write account name.
	accountMsg := fmt.Sprintf("%s\n\n", account)
	_, err = c.connection.Write([]byte(accountMsg))
	if err != nil {
		return err
	}

	// Write password.
	passwordMsg := fmt.Sprintf("%s\n\n", password)
	_, err = c.connection.Write([]byte(passwordMsg))
	if err != nil {
		return err
	}

	return nil
}

// Write will take the message and write it over the connection.
func (c *Client) Write(message string) error {
	if c.connection == nil {
		return ErrNotConnected
	}

	msg := fmt.Sprintf("%s\n\n", message)

	// Write message.
	_, err := c.connection.Write([]byte(msg))
	if err != nil {
		return ErrUnableToWrite
	}

	return nil
}

// Whisper is a helper function to whisper a specific account with a message.
func (c *Client) Whisper(account string, message string) error {
	if c.connection == nil {
		return ErrNotConnected
	}

	msg := fmt.Sprintf("/msg %s %s\n\n", account, message)

	// Whisper message to the account.
	_, err := c.connection.Write([]byte(msg))
	if err != nil {
		return ErrUnableToWrite
	}

	return nil
}

// Reads all the output of the tcp connection on a channel.
func (c *Client) Read(ch chan []byte, errors chan error) error {
	if c.connection == nil {
		return ErrNotConnected
	}

	// Make a byte slice to read data into.
	buf := make([]byte, 1024)

	go func(ch chan []byte, errs chan error) {
		for {
			// Read the amount of bytes into the buf.
			bytes, err := c.connection.Read(buf)
			if err != nil {
				errs <- err
				return
			}

			// Only send the newly written bytes to the channel.
			ch <- buf[:bytes]
		}
	}(ch, errors)

	return nil
}

// New creates a new client with all dependencies.
func New() *Client {
	return &Client{}
}
