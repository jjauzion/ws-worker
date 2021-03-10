package client

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CodeNoTaskinQueue = iota + 600
)

var (
	pebkac = func(s string, errCode codes.Code) error { return status.Error(errCode, s) }

	errBadReq         = pebkac("bad request data", codes.InvalidArgument)
	errForbidden      = pebkac("forbidden", codes.PermissionDenied)
	errNoTasksInQueue = pebkac("no tasks in queue", CodeNoTaskinQueue)
)

func exactlyOneOf(fields string) error {
	return pebkac("exactly one of "+fields+" must be set", codes.InvalidArgument)
}

func getErrorCode(err error) uint32 {
	return uint32(status.Code(err))
}
