package postgresql

import "gorm.io/gorm"

func DeleteRecordFromDb(db *gorm.DB, record interface{}) error {
	tx := db.Delete(record)
	return tx.Error
}
