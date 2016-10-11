package comm

import (
	"bufio"
	"fmt"
)

type TelnetAuthenticator struct {
	Auth Auth
}

func (a *TelnetAuthenticator) Authenticate(c Connection) error {
	r := bufio.NewReader(c)
	_, err := r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate login failed: %s", err)
	}

	_, err = c.Write([]byte(a.Auth.Login + "\r\n"))
	if err != nil {
		return fmt.Errorf("authenticate write login failed: %s", err)
	}

	_, err = r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate password failed: %s", err)
	}

	_, err = c.Write([]byte(a.Auth.Password + "\r\n"))
	if err != nil {
		return fmt.Errorf("authenticate password failed: %s", err)
	}
	return nil
}
