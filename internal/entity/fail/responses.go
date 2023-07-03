package fail

import "github.com/mephistolie/chefbook-backend-common/responses/fail"

var (
	GrpcInvalidPublicKey           = fail.CreateGrpcClient(fail.TypeInvalidBody, "unable to parse public key")
	GrpcPrivateKeyLengthOutOfRange = fail.CreateGrpcClient(fail.TypeInvalidBody, "encrypted private key length is out of acceptable length")

	GrpcInvalidCode = fail.CreateGrpcClient(fail.TypeInvalidBody, "invalid code")

	GrpcNoVault = fail.CreateGrpcClient(fail.TypeAccessDenied, "user haven't encrypted vault")

	GrpcRecipeKeyLengthOutOfRange = fail.CreateGrpcClient(fail.TypeInvalidBody, "encrypted recipe key length is out of acceptable length")
	GrpcOwnedRecipeKeyDeletion    = fail.CreateGrpcClient(fail.TypeAccessDenied, "owned recipe key can't be deleted; delete entire recipe")
)
