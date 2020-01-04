package chat

import (
	"gp-websoket/impl"
	"gp-websoket/model"
)

type Client struct {
	User model.User
	Cone impl.Connection `json:"-"`
}
