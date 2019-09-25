package testmethod21

import (
	"github.com/viethqc/gogstash/KDGoLib/apimgr"
	"github.com/viethqc/gogstash/KDGoLib/apimgr/testapi"
)

func init() {
	testapi.Manager.Add(
		apimgr.Definition{
			Description: `
				Test api 2 method 21
			`,
			Method:  "GET",
			Pattern: "/1/testmethod21/testres2",
			Request: TestAPI21{},
		},
	)
}

type TestAPI21 struct{}
