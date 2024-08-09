package cronjobs

import (
	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/services/actions"
)

func SendNotifications(extReq request.ExternalRequest, db storage.Database) {
	notificationRecord := models.NotificationRecord{}

	res, err := notificationRecord.PopFromQueue(db.Redis)

	if err != nil {
		extReq.Logger.Error("error getting notificatin records: ", err.Error())
		return
	}

	extReq.Logger.Error("Sending records found: ", res)

	err = actions.Send(extReq, db.Postgresql, db.Redis, &res)

	if err != nil {
		extReq.Logger.Error("error getting notificatin records: ", err.Error())
		return
	}

}
