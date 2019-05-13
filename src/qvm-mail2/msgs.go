package main

var (
	ecs_upgrade_msg                         = `{"data":{"instanceId":"i-xxxxxxxxxx","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_upgrade","messageSource":"console","source":"console","timeStamp":1516800705630,"timestamp":1516800741,"uniqueID":"3aa1fcea-7dce-4ef0-a28e-c6425619e36c","userID":"1103909446200972"}`
	ecs_bandwidth_downgrade_msg             = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","internetTx":102400,"intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","vswitchInstanceId":"v-2ze1f71lXXXjxh7cowvc"},"eventID":"ecs_bandwidth_downgrade","messageSource":"console","source":"console","timeStamp":1516953245439,"timestamp":1516953247,"uniqueID":"f39e58d3-a83b-4099-9124-5c05ac582d73","userID":"1103909446200972"}`
	ecs_renewal_upgrade_msg                 = `{"data":{"expectedRestartTime":"2016/12/22 10:10","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_renewal_upgrade","messageSource":"console","source":"console","timeStamp":1516953281087,"timestamp":1516953290,"uniqueID":"fa240627-7b12-46e7-b6aa-297cce7b0804","userID":"1103909446200972"}`
	disk_create_msg                         = `{"data":{"mcUserName":"abc***@aliyun.com","mergeCount":"8"},"eventID":"disk_create","messageSource":"console","source":"console","timeStamp":1516953323796,"timestamp":1516953334,"uniqueID":"a1424af3-c1b8-4024-bcad-841c02c4acda","userID":"1103909446200972"}`
	disk_release_msg                        = `{"data":{"mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"disk_release","messageSource":"console","source":"console","timeStamp":1516953422077,"timestamp":1516953424,"uniqueID":"309b7c2b-fc39-48b3-a06b-1a33b5f97f33","userID":"1103909446200972"}`
	vm_bandwidth_change_billing_msg         = `{"data":{"billingTag":"4","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"vm_bandwidth_change_billing","messageSource":"console","source":"console","timeStamp":1516953582915,"timestamp":1516953585,"uniqueID":"341c4097-67b9-40ec-8bf1-05ce56413520","userID":"1103909446200972"}`
	ecs_about_to_release_by_usersetting_msg = `{"data":{"mcUserName":"abc***@aliyun.com","mergeCount":"7","willReleaseTime":"2016-12-22 10:10"},"eventID":"ecs_about_to_release_by_usersetting","messageSource":"console","source":"console","timeStamp":1517552080725,"timestamp":1517552255,"uniqueID":"55c4644b-d905-4b60-8ddc-04e916673a7d","userID":"1103909446200972"}`
	ecs_resource_auto_restart_msg           = `{"data":{"expectedRestartTime":"2016-12-22 10:10","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_resource_auto_restart","messageSource":"console","source":"console","timeStamp":1516960573527,"timestamp":1516960581,"uniqueID":"ddbd7efe-4043-439b-95be-038c3d7c26c2","userID":"1103909446200972"}`
	ecs_about_to_restart_msg                = `{"data":{"expectedRestartTime":"2016-12-22 10:10","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_about_to_restart","messageSource":"console","source":"console","timeStamp":1516955063937,"timestamp":1516955076,"uniqueID":"d5248a4f-e6e2-48cd-af11-d07ec7f1ffc5","userID":"1103909446200972"}`
	ecs_about_to_release_afterpay_msg       = `{"data":{"diskCount":"7","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7","willReleaseTime":"2016/12/22 10:10"},"eventID":"ecs_about_to_release_afterpay","messageSource":"console","source":"console","timeStamp":1516955927362,"timestamp":1516955929,"uniqueID":"d8d166ca-05d2-48dc-9439-dab8550269c7","userID":"1103909446200972"}`
	ecs_release_by_usersetting_msg          = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7","willReleaseTime":"2016/12/22 10:10"},"eventID":"ecs_release_by_usersetting","messageSource":"console","source":"console","timeStamp":1516956126031,"timestamp":1516956128,"uniqueID":"1c54e057-8529-4706-b2ba-42945cc3b674","userID":"1103909446200972"}`
	disk_about_to_release_msg               = `{"data":{"mcUserName":"abc***@aliyun.com","mergeCount":"7","willReleaseTime":"2016-12-22 10:10"},"eventID":"disk_about_to_release","messageSource":"console","source":"console","timeStamp":1516956302381,"timestamp":1516956304,"uniqueID":"48ca691a-f3a6-4140-8152-e8f0db2fff2d","userID":"1103909446200972"}`
	ecs_release_msg                         = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_release","messageSource":"console","source":"console","timeStamp":1516957590161,"timestamp":1516957594,"uniqueID":"652b6702-5bec-4775-8b4a-c517974d382a","userID":"1103909446200972"}`
	ecs_release_afterpay_msg                = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com"},"eventID":"ecs_release_afterpay","messageSource":"console","source":"console","timeStamp":1516957632676,"timestamp":1516957651,"uniqueID":"a91482cb-a86d-4611-953d-401e761acac1","userID":"1103909446200972"}`
	ecs_renewal_msg                         = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_renewal","messageSource":"console","source":"console","timeStamp":1516957876471,"timestamp":1516957890,"uniqueID":"147f8a43-9f8b-450e-ae78-a3b8fd46a64e","userID":"1103909446200972"}`
	ecs_expired_msg                         = `{"data":{"diskCount":"7","endTime":"2016-12-22 10:10","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","keepDays":"8","lastDay":"2016-12-22 10:10","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_expired","messageSource":"console","source":"console","timeStamp":1516958044280,"timestamp":1516958053,"uniqueID":"d2e81664-e888-43e8-8eed-400f0417e6c1","userID":"1103909446200972"}`
	ecs_about_to_release_msg                = `{"data":{"diskCount":"7","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7","willReleaseTime":"2016/12/22 10:10"},"eventID":"ecs_about_to_release","messageSource":"console","source":"console","timeStamp":1516958087132,"timestamp":1516958089,"uniqueID":"8dc72d66-a21c-4c11-b03e-941cc8fa362c","userID":"1103909446200972"}`
	commodity_will_auto_renewal_notify_msg  = `{"data":{"mcUserName":"abc***@aliyun.com","productName":"ECS","renewalTime":"2016/12/22"},"eventID":"commodity_will_auto_renewal_notify","messageSource":"console","source":"console","timeStamp":1516958122882,"timestamp":1516958154,"uniqueID":"e3ba8f21-0a8e-4859-9503-2247f261d568","userID":"1103909446200972"}`
	sp_ecs_create_msg                       = `{"data":{"dataList":[{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","privateIp":"10.44.41.140"}],"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","isWin":true,"mcUserName":"abc***@aliyun.com","mergeCount":"7","vswitchInstanceId":"v-2ze1f71lXXXjxh7cowvc"},"eventID":"sp_ecs_create","messageSource":"console","source":"console","timeStamp":1516958178851,"timestamp":1516958217,"uniqueID":"a18f05a3-ab1d-4cb1-bcdd-92fc2857aa22","userID":"1103909446200972"}`
	ecs_cleanup_msg                         = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"8"},"eventID":"ecs_cleanup","messageSource":"console","source":"console","timeStamp":1516958277145,"timestamp":1516958279,"uniqueID":"d1f77f7c-da78-4916-adb6-f558c542f030","userID":"1103909446200972"}`
	ecs_afterpay_expire_msg                 = `{"data":{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_afterpay_expire","messageSource":"console","source":"console","timeStamp":1516958311916,"timestamp":1516958332,"uniqueID":"f8eb2857-bc7e-4bfb-8d8c-6c2c78c60051","userID":"1103909446200972"}`
	ecs_about_to_expire_15_msg              = `{"data":{"atDays":"15","dataList":[{"endTime":"2016/12/22 10:10","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","privateIp":"10.44.41.140","vswitchInstanceId":"v-2ze1f71lXXXjxh7cowvc"}],"endTime":"2016/12/22 10:10","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_about_to_expire_15","messageSource":"console","source":"console","timeStamp":1516958362938,"timestamp":1516958366,"uniqueID":"401cfb14-8063-4b9b-916d-55e9af01e5f0","userID":"1103909446200972"}`
	ecs_about_to_expire_30_msg              = `{"data":{"atDays":"5","dataList":[{"endTime":"2016/12/22 10:10","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","privateIp":"10.44.41.140","vswitchInstanceId":"v-2ze1f71lXXXjxh7cowvc"}],"endTime":"2016/12/22 10:10","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_about_to_expire_30","messageSource":"console","source":"console","timeStamp":1516958398543,"timestamp":1516958400,"uniqueID":"6482a424-8c8c-4278-87d0-69a42a79bd03","userID":"1103909446200972"}`
	ecs_about_to_expire_msg                 = `{"data":{"atDays":"2016/12/22 10:10","dataList":[{"instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","vswitchInstanceId":"v-2ze1f71lXXXjxh7cowvc"}],"endTime":"2016/12/22 AM 10:10:00","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_about_to_expire","messageSource":"console","source":"console","timeStamp":1516958554699,"timestamp":1516958561,"uniqueID":"97566cbf-a1ce-40d5-ba7b-ea5bb637bf82","userID":"1103909446200972"}`
	ecs_free_trial_about_to_expire_msg      = `{"data":{"atDays":"5","diskCount":"7","endTime":"2016/12/22 AM 10:10:00","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","mcUserName":"abc***@aliyun.com","mergeCount":"7"},"eventID":"ecs_free_trial_about_to_expire","messageSource":"console","source":"console","timeStamp":1516958673750,"timestamp":1516958675,"uniqueID":"0b0e5b82-0f15-44f5-94a9-6884f4a67fdd","userID":"1103909446200972"}`
	ecs_compensate_renew_msg                = `{"data":{"compensatoryReason":"1 houre outrage compensatory","endTime":"2016/12/22","instanceId":"i-2ze1f71lXXXjxh7cowvc","instanceName":"tb1133563_2012","internetIp":"101.200.152.51","intranetIp":"10.44.41.140","mcUserName":"abc***@aliyun.com","mergeCount":7,"startTime":"2016/11/22","vswitchInstanceId":"v-2ze1f71lXXXjxh7cowvc"},"eventID":"ecs_compensate_renew","messageSource":"console","source":"console","timeStamp":1516960305678,"timestamp":1516960307,"uniqueID":"2b91934c-a453-455a-bf9e-f8dacf09dd44","userID":"1103909446200972"}`
	fenghuotai_vm_down_start_msg            = `{"data":{"instanceId":"iZp7rbtfvwnwx6iv3cio8tZ/11.239.171.142","mcUserName":"abc***@aliyun.com"},"eventID":"fenghuotai_vm_down_start","messageSource":"console","source":"console","timeStamp":1516960454664,"timestamp":1516960463,"uniqueID":"4e4d4d16-705c-429e-9a74-d33287a7dfdb","userID":"1103909446200972"}`
	fenghuotai_vm_down_end_msg              = `{"data":{"instanceId":"iZp7rbtfvwnwx6iv3cio8tZ/11.239.171.142","mcUserName":"abc***@aliyun.com"},"eventID":"fenghuotai_vm_down_end","messageSource":"console","source":"console","timeStamp":1516960484589,"timestamp":1516960493,"uniqueID":"fd4920c1-c7f1-4d5d-b59e-f25f36e1a269","userID":"1103909446200972"}`
)