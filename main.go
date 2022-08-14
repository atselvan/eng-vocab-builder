package main

import (
	"github.com/privatesquare/bkst-go-utils/utils/httputils"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
)

//go:generate oapi-codegen -old-config-style -generate gin -package main -o gen.go api/eng-vocab-builder_v1.yaml
func main() {
	err := SetupRoutes().Run(":8000")
	if err != nil {
		logger.Error(httputils.ServerStartupErrMsg, err)
	}
}
