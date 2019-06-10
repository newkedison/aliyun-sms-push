package main

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	//   "go.mongodb.org/mongo-driver/mongo"
	//   "go.mongodb.org/mongo-driver/bson/primitive"
	//   "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/kardianos/service"
)

var (
	dump     = spew.Dump
	now      = time.Now
	print    = log.Println
	ctxEmpty = context.TODO()
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
	port uint16
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		print("Running in terminal.")
	} else {
		print("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	print("I'm running", service.Platform())
	tickerAlive := time.NewTicker(60 * time.Second)
	tickerUpdateConfig := time.NewTicker(30 * time.Second)
	err := connectToDB(globalConfig.MongoDB.URI)
	if err != nil {
		return err
	}
	test_db()
	updateWhiteList(getContextWithTimeout(200))
	go func() {
		for {
			select {
			case tm := <-tickerAlive.C:
				print("Still running at ", tm)
			case <-tickerUpdateConfig.C:
				updateWhiteList(getContextWithTimeout(1000))
			case <-p.exit:
				tickerAlive.Stop()
				tickerUpdateConfig.Stop()
				print("disconnecting")
				err = dbClient.Disconnect(getContextWithTimeout(1000))
				if err != nil {
					panic(err)
				}
				print("disconnected")
			}
		}
	}()
	return StartServer(globalConfig.BaseAddr)
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	print("I'm Stopping!")
	close(p.exit)
	return nil
}

func main() {
	prg := &program{}
	svcConfig := &service.Config{
		Name:        "sms-push",
		DisplayName: "SMS Push",
		Description: "Send sms to user",
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	globalConfig = readConfig(DefaultConfigFile)
	if globalConfig == nil {
		panic("读取配置文件失败")
	}

	err = s.Run()
	if err != nil {
		print(err)
	}
}

func updateAllowedPhoneListFromDb(ctx context.Context) {
	cursor, err := colPhoneList.Find(ctx, bson.D{})
	if err != nil {
		dump(err)
	}
	var result []allowedPhone
	for cursor.Next(ctx) {
		var r allowedPhone
		err := cursor.Decode(&r)
		if err != nil {
			dump(err)
			break
		}
		result = append(result, r)
	}
	for _, s := range result {
		allowPhoneList = append(allowPhoneList, s.Phone)
	}
}

func updateAllowedIpListFromDb(ctx context.Context) {
	cursor, err := colIpList.Find(ctx, bson.D{})
	if err != nil {
		dump(err)
	}
	var result []allowedIp
	for cursor.Next(ctx) {
		var r allowedIp
		err := cursor.Decode(&r)
		if err != nil {
			dump(err)
			break
		}
		result = append(result, r)
	}
	for _, s := range result {
		allowIpList = append(allowIpList, s.Ip)
	}
}

func updateWhiteList(ctx context.Context) {
	mutexAllowList.Lock()
	defer mutexAllowList.Unlock()
	allowPhoneList = readLines(getCurrentPath() + "/phone.txt")
	allowIpList = readLines(getCurrentPath() + "/ip.txt")
	updateAllowedPhoneListFromDb(ctx)
	updateAllowedIpListFromDb(ctx)
}
