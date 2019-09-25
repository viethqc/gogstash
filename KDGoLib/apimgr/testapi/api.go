package testapi

import "github.com/viethqc/gogstash/KDGoLib/apimgr"

type TestAPI struct{}

var Manager = apimgr.NewManager(TestAPI{})
