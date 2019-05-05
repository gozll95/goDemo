# 分布式crontab


crontab

[cron]
enable=true
[cron.specs]
"0 0 0 * * *"=["UpdateSystemImage"]
"0 0 1 * * *"=["UpdateOms"]


type CronConfig struct {
	Enable bool                `toml:"enable"`
	Specs  map[string][]string `toml:"specs"` // one spec to multi events
}

func (config *CronConfig) Events() (res map[string]string) {
	res = make(map[string]string)
	for spec, events := range config.Specs {
		for _, event := range events {
			res[event] = spec
		}
	}

	return
}

func NewCron(config *conf.CronConfig) (c *Cron) {
	c = &Cron{
		croner: cron.New(),
		config: config,
	}

	// config enable cron or need manual operation start cron
	if config != nil && config.Enable {
		c.Start()
	}

	return
}



# config + cron.New()-->CronJob

addEvent:
func (c *Cron) addEvent(name string, spec string, event Event) {
	c.croner.AddFunc(spec, func() {
		// this make cron event
		now := time.Now()
		nextTime, err := c.nextTime(spec, now)
		if err != nil {
			utils.StdLog.Errorf("parse spec %s error %v", spec, err)
			return
		}

		cronModel := model.NewCronModel(name, now, nextTime)
		err = cronModel.Upsert()
		if err != nil {
			utils.StdLog.Warnf("cron event %s blocked by error %v", name, err)
			return
		}
		event()
	})

}