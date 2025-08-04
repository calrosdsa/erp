package userservice

import (
	"erp/internal/app/connection"
	"time"
)

type AdminService struct {
	conn    *connection.Connection
	timeout time.Duration
}



