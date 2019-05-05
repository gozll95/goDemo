package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type InitMgoReplSetOpts struct {
	ReplSetName string      `json:"replset_name"`
	Members     []MgoMember `json:"members"`
}
type MgoMember struct {
	Host        string `json:"host"`
	ArbiterOnly bool   `json:"arbiter_only"`
}

func InitMgoReplSet(mgoURL, setName string, members []MgoMember) error {
	mgo.SetLogger(mgoLogger)
	info, err := mgo.ParseURL(mgoURL)
	if err != nil {
		return err
	}
	info.Direct = true
	info.Timeout = time.Second * 3
	sess, err := mgo.DialWithInfo(info)
	if err != nil {
		return err
	}
	sess.SetMode(mgo.Monotonic, true)

	mgoMembers := make([]bson.M, len(members))
	for i, member := range members {
		mgoMembers[i] = bson.M{
			"_id":         i,
			"host":        member.Host,
			"arbiterOnly": member.ArbiterOnly,
		}
	}
	result := make(bson.M)
	defer log.Printf("%+v\n", result)
	return sess.Run(bson.D{{
		Name: "replSetInitiate",
		Value: bson.M{
			"_id":     setName,
			"members": mgoMembers,
		},
	}}, &result)
}

func InitMgoReplSetCmd(args []string) {
	var (
		conf string
	)
	fl := flag.NewFlagSet("InitMgoReplSet", flag.ExitOnError)
	fl.StringVar(&conf, "c", "mongo.json", "mongo replication set config file")
	fl.Parse(args)

	f, err := os.Open(conf)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	initOpts := make([]InitMgoReplSetOpts, 0)
	err = json.NewDecoder(f).Decode(&initOpts)
	if err != nil {
		panic(err)
	}
	for _, opts := range initOpts {
		if len(opts.Members) < 2 {
			log.Printf("[ERROR] count of replication set members %s must be greater than 1, skip\n", opts.ReplSetName)
			continue
		}
		mgoURL := fmt.Sprintf("mongodb://%s/", opts.Members[0].Host)
		err := InitMgoReplSet(mgoURL, opts.ReplSetName, opts.Members)
		if err != nil {
			log.Printf("[ERROR] init replication set %s through %s error: %s\n", opts.ReplSetName, mgoURL, err)
		}
	}
}
