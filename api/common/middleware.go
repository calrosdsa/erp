package common

import "github.com/danielgtaylor/huma/v2"

type Middleware func(ctx huma.Context, next func(huma.Context))
