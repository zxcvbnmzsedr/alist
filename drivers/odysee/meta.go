package odysee

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	driver.RootID
	AuthToken string `json:"authToken" required:"true"`
}

var config = driver.Config{
	Name:        "Odysee",
	DefaultRoot: "root",
}

func New() driver.Driver {
	return &Odysee{}
}

func init() {
	op.RegisterDriver(config, New)
}
