package domain

type Config struct {
	id        string
	refresh   string
	createdAt int64
	//refresh   bool
}

func (c *Config) ID() string {
	return c.id
}
func (c *Config) SetID(thisId string) {
	c.id = thisId
}

func (c *Config) Refresh() string {
	return c.refresh
}
func (c *Config) SetRefresh(b string) {
	c.refresh = b
}

func (c *Config) CreatedAt() int64 {
	return c.createdAt
}
func (c *Config) SetCreatedAt(t int64) {
	c.createdAt = t
}
