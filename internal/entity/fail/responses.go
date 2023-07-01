package fail

import "github.com/mephistolie/chefbook-backend-common/responses/fail"

var (
	GrpcInvalidCode = fail.CreateGrpcClient(fail.TypeInvalidBody, "invalid code")
)
