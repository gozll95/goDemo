package dao

import (
	"go-kit-gin/entity"

	"github.com/sirupsen/logrus"
)

//CreateAccount 创建账户
func CreateAccount(account *entity.UserAccount) error {
	if err := createModel(account); err != nil {
		logrus.Errorln("CreateAccount ", err)
		return err
	}
	return nil
}

const accountSQL = "account = ?"

//Account 获取账户
func Account(account string) (entity.UserAccount, error) {
	var myAccount entity.UserAccount
	if mydb := DBInst().Where(accountSQL, account).First(&myAccount); mydb.Error != nil {
		logrus.Errorln("Account ", mydb.Error)
		return myAccount, mydb.Error
	}
	return myAccount, nil
}

//AccountById 通过id查账户
func AccountById(id uint) (entity.UserAccount, error) {
	var myAccount entity.UserAccount
	if mydb := DBInst().Where("id = ?", id).First(&myAccount); mydb.Error != nil {
		return myAccount, mydb.Error
	}
	return myAccount, nil
}

const updateAccountSQL = "account = ?"

//UpdateAccount 更新账户
func UpdateAccount(account string, userAccount entity.UserAccount) error {
	var myAccount entity.UserAccount
	if mydb := DBInst().Where(updateAccountSQL, account).First(&myAccount); mydb.Error != nil {
		logrus.Errorln("UpdateAccount First", mydb.Error)
		return mydb.Error
	}
	if mydb := DBInst().Model(&myAccount).Updates(userAccount); mydb.Error != nil {
		logrus.Errorln("UpdateAccount Updates", mydb.Error)
		return mydb.Error
	}
	return nil
}
