package modloader

import (
	codecjson "github.com/viethqc/gogstash/codec/json"
	"github.com/viethqc/gogstash/config"
	filteraddfield "github.com/viethqc/gogstash/filter/addfield"
	filtercond "github.com/viethqc/gogstash/filter/cond"
	filterdate "github.com/viethqc/gogstash/filter/date"
	filtergeoip2 "github.com/viethqc/gogstash/filter/geoip2"
	filtergonx "github.com/viethqc/gogstash/filter/gonx"
	filtergrok "github.com/viethqc/gogstash/filter/grok"
	filterjson "github.com/viethqc/gogstash/filter/json"
	filtermutate "github.com/viethqc/gogstash/filter/mutate"
	filterratelimit "github.com/viethqc/gogstash/filter/ratelimit"
	filterremovefield "github.com/viethqc/gogstash/filter/removefield"
	filtertypeconv "github.com/viethqc/gogstash/filter/typeconv"
	filterurlparam "github.com/viethqc/gogstash/filter/urlparam"
	filteruseragent "github.com/viethqc/gogstash/filter/useragent"
	inputbeats "github.com/viethqc/gogstash/input/beats"
	inputdockerlog "github.com/viethqc/gogstash/input/dockerlog"
	inputdockerstats "github.com/viethqc/gogstash/input/dockerstats"
	inputexec "github.com/viethqc/gogstash/input/exec"
	inputfile "github.com/viethqc/gogstash/input/file"
	inputhttp "github.com/viethqc/gogstash/input/http"
	inputhttplisten "github.com/viethqc/gogstash/input/httplisten"
	inputlorem "github.com/viethqc/gogstash/input/lorem"
	inputredis "github.com/viethqc/gogstash/input/redis"
	inputsocket "github.com/viethqc/gogstash/input/socket"
	outputamqp "github.com/viethqc/gogstash/output/amqp"
	outputcond "github.com/viethqc/gogstash/output/cond"
	outputelastic "github.com/viethqc/gogstash/output/elastic"
	outputemail "github.com/viethqc/gogstash/output/email"
	outputfile "github.com/viethqc/gogstash/output/file"
	outputhttp "github.com/viethqc/gogstash/output/http"
	outputprometheus "github.com/viethqc/gogstash/output/prometheus"
	outputredis "github.com/viethqc/gogstash/output/redis"
	outputreport "github.com/viethqc/gogstash/output/report"
	outputstdout "github.com/viethqc/gogstash/output/stdout"

	inputrabbitmq "github.com/viethqc/gogstash/input/rabbitmq"
)

func init() {
	config.RegistInputHandler(inputbeats.ModuleName, inputbeats.InitHandler)
	config.RegistInputHandler(inputdockerlog.ModuleName, inputdockerlog.InitHandler)
	config.RegistInputHandler(inputdockerstats.ModuleName, inputdockerstats.InitHandler)
	config.RegistInputHandler(inputexec.ModuleName, inputexec.InitHandler)
	config.RegistInputHandler(inputfile.ModuleName, inputfile.InitHandler)
	config.RegistInputHandler(inputhttp.ModuleName, inputhttp.InitHandler)
	config.RegistInputHandler(inputhttplisten.ModuleName, inputhttplisten.InitHandler)
	config.RegistInputHandler(inputlorem.ModuleName, inputlorem.InitHandler)
	config.RegistInputHandler(inputredis.ModuleName, inputredis.InitHandler)
	config.RegistInputHandler(inputsocket.ModuleName, inputsocket.InitHandler)
	config.RegistInputHandler(inputrabbitmq.ModuleName, inputrabbitmq.InitHandler)

	config.RegistFilterHandler(filteraddfield.ModuleName, filteraddfield.InitHandler)
	config.RegistFilterHandler(filtercond.ModuleName, filtercond.InitHandler)
	config.RegistFilterHandler(filterdate.ModuleName, filterdate.InitHandler)
	config.RegistFilterHandler(filtergeoip2.ModuleName, filtergeoip2.InitHandler)
	config.RegistFilterHandler(filtergonx.ModuleName, filtergonx.InitHandler)
	config.RegistFilterHandler(filtergrok.ModuleName, filtergrok.InitHandler)
	config.RegistFilterHandler(filterjson.ModuleName, filterjson.InitHandler)
	config.RegistFilterHandler(filtermutate.ModuleName, filtermutate.InitHandler)
	config.RegistFilterHandler(filterratelimit.ModuleName, filterratelimit.InitHandler)
	config.RegistFilterHandler(filterremovefield.ModuleName, filterremovefield.InitHandler)
	config.RegistFilterHandler(filtertypeconv.ModuleName, filtertypeconv.InitHandler)
	config.RegistFilterHandler(filteruseragent.ModuleName, filteruseragent.InitHandler)
	config.RegistFilterHandler(filterurlparam.ModuleName, filterurlparam.InitHandler)

	config.RegistOutputHandler(outputamqp.ModuleName, outputamqp.InitHandler)
	config.RegistOutputHandler(outputcond.ModuleName, outputcond.InitHandler)
	config.RegistOutputHandler(outputelastic.ModuleName, outputelastic.InitHandler)
	config.RegistOutputHandler(outputemail.ModuleName, outputemail.InitHandler)
	config.RegistOutputHandler(outputhttp.ModuleName, outputhttp.InitHandler)
	config.RegistOutputHandler(outputprometheus.ModuleName, outputprometheus.InitHandler)
	config.RegistOutputHandler(outputredis.ModuleName, outputredis.InitHandler)
	config.RegistOutputHandler(outputreport.ModuleName, outputreport.InitHandler)
	config.RegistOutputHandler(outputstdout.ModuleName, outputstdout.InitHandler)
	config.RegistOutputHandler(outputfile.ModuleName, outputfile.InitHandler)

	config.RegistCodecHandler(config.DefaultCodecName, config.DefaultCodecInitHandler)
	config.RegistCodecHandler(codecjson.ModuleName, codecjson.InitHandler)
}
