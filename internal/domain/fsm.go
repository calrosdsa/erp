package domain

import "erp/api/common"

type State[T any] func(req *common.RequestContext, args T) (T, State[T], error)

func Run[T any](req *common.RequestContext, args T, start State[T]) (T, error) {
	var err error
	current := start
	for {
		if req.Ctx.Err() != nil {
			return args, req.Ctx.Err()
		}
		args, current, err = current(req, args)
		if err != nil {
			return args, err
		}
		if current == nil {
			return args, nil
		}
	}
}