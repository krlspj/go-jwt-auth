package domain

type Config struct {
	id        string
	refresh   bool
	createdAt uint32
}

func (c *Config) ID() string {
	return c.id
}
func (c *Config) SetID(thisId string) {
	c.id = thisId
}

func (c *Config) Refresh() bool {
	return c.refresh
}
func (c *Config) SetRefresh(b bool) {
	c.refresh = b
}

func (c *Config) CreatedAt() uint32 {
	return c.createdAt
}
func (c *Config) SetCreatedAt(t uint32) {
	c.createdAt = t
}
