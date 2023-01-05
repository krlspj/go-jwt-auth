package domain

type EmbeddedStruct struct {
	name    string
	surname string
}

func (es *EmbeddedStruct) Name() string {
	return es.name
}
func (es *EmbeddedStruct) SetName(s string) {
	es.name = s
}
func (es *EmbeddedStruct) Surname() string {
	return es.surname
}
func (es *EmbeddedStruct) SetSurname(s string) {
	es.surname = s
}

type Config struct {
	id        string
	refresh   string
	createdAt int64
	refreshB  *bool
	EmbeddedStruct
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

func (c *Config) RefreshB() *bool {
	return c.refreshB
}
func (c *Config) SetRefreshB(b bool) {
	c.refreshB = &b
}
func (c *Config) CreatedAt() int64 {
	return c.createdAt
}
func (c *Config) SetCreatedAt(t int64) {
	c.createdAt = t
}
