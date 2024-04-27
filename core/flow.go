package core

import "time"

func GetConfigUnAuth(cli *Nacos) error {
	if err := cli.GetNameSpace(); err != nil {
		return err
	}
	if err := cli.GetConfig(); err != nil {
		return err
	}
	return nil
}

func GetConfigWithAuth(cli *Nacos) error {
	if err := cli.AddUser(); err != nil {
		return err
	}
	time.Sleep(4 * time.Second)
	if err := cli.SetJwt(); err != nil {
		return err
	}
	time.Sleep(4 * time.Second)
	if err := cli.GetNameSpace(); err != nil {
		return err
	}
	if err := cli.GetConfig(); err != nil {
		return err
	}
	return nil
}
