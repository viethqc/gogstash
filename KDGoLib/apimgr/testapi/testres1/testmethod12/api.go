package testmethod12

import (
	"github.com/viethqc/gogstash/KDGoLib/apimgr"
	"github.com/viethqc/gogstash/KDGoLib/apimgr/testapi"
)

func init() {
	testapi.Manager.Add(
		apimgr.Definition{
			Description: `
				Test api 1 method 12
			`,
			Method:  "GET",
			Pattern: "/1/testmethod12/testres1",
			Request: TestAPI12{},
		},
	)
}

type TestAPI12 struct{}
