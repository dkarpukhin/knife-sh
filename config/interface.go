package config

func (c *Config) Hosts() []string {
	return c.hosts
}

func (c *Config) ConnectTimeout() int64 {
	return c.timeoutSshConnect
}

func (c *Config) ExecTimeout() int64 {
	return c.timeoutExec
}

func (c *Config) Command() string {
	return c.command
}

func (c *Config) Concurrency() int64 {
	return c.concurrency
}

func (c *Config) SshKeyContent() string {
	return c.sshKey
}

func (c *Config) SshUser() string {
	return c.sshUser
}
