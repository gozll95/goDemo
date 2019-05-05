package tasks

import (
	"fmt"
	"time"

	"MyApp/models"                     // you probably want to access your models in the task
	"github.com/astaxie/beego/toolbox" // the toolbox package
	// to keep the business logic at the right place
)

func init() {
	first_task := toolbox.NewTask("first_task", "0/30 * * * * *", func() error {
		// this task will run every 30 seconds
		campaigns, err := models.GetFinishedCampaigns()
		if err != nil {
			fmt.Println("Could not load finished campaigns with error:", err)
			return err
		}

		if len(campaigns) == 0 {
			fmt.Println("No campaigns finished yet. Exiting.")
			return nil
		}

		for _, campaign := range campaigns {
			result, err := models.SendCampaignNotification(campaign)
			if err != nil {
				fmt.Printf("\nCampaign %d could not be notified!\n", campaign.Id)
			}

			if result {
				fmt.Printf("\nCampaign %d was notified!\n", campaign.Id)
			}
		}

		fmt.Printf("\nNotification task ran at: %s\n", time.Now())
		return nil
	})

	toolbox.AddTask("first_task", first_task)
	toolbox.StartTask()
	defer toolbox.StopTask()
}
