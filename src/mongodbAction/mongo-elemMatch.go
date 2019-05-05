//db.rds_network.find({"region_id":"cn-beijing","db_network_attribute.db_instance_net_info.connection_string":"a4c7ae5d-b934-4fa4-9888-b397482ab6bd"},{"db_network_attribute.db_instance_net_info.$":1}).pretty()
func (_ *_RdsNetwork) FindByUidRegionIdDBInstanceIdAndQuery(uid uint32, regionId, dBInstanceId string, subQuery bson.M) (rdsNetwork *RdsNetworkModel, err error) {
	if uid == 0 || regionId == "" || dBInstanceId == "" {
		return nil, ErrInvalidParams
	}
	RdsNetwork.Query(func(c *mgo.Collection) {
		query := bson.M{
			"uid":            uid,
			"region_id":      regionId,
			"db_instance_id": dBInstanceId,
			"db_network_attribute.db_instance_net_info": bson.M{
				// "$elemMatch": bson.M{
				// 	"connection_string": "a4c7ae5d-b934-4fa4-9888-b397482ab6bd",
				// },
				"$elemMatch": subQuery,
			},
		}
		projection := bson.M{
			"db_network_attribute.db_instance_net_info.$": 1,
		}
		err = c.Find(query).Select(projection).One(&rdsNetwork)
	})

	return
}

/*
db.rds_network.find({


"region_id":"cn-beijing",
"db_network_attribute.db_instance_net_info":{"$elemMatch":{"connection_string":"a4c7ae5d-b934-4fa4-9888-b397482ab6bd"}},
"db_network_attribute.db_instance_net_info.$":1


});

db.getCollection('web_mem_favorites').find


(

{"_id":NumberLong(1181675746),"favorite_shards.sid":NumberLong(577)},{"favorite_shards.$":1})


.pretty()


db.rds_network.find({"region_id":"cn-beijing","db_network_attribute.db_instance_net_info.connection_string":"a4c7ae5d-b934-4fa4-9888-b397482ab6bd"},{"db_network_attribute.db_instance_net_info.$":1}).pretty()





*/