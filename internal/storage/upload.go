package storage

import (
	"os"

	"github.com/asdine/storm/q"
	log "github.com/sirupsen/logrus"

	"github.com/systemli/ticker/internal/model"
)

//FindUploadsByMessage returns all uploads for a Message.
func FindUploadsByMessage(message *model.Message) []*model.Upload {
	var uploads []*model.Upload

	if len(message.Attachments) > 0 {
		var uuids []string
		for _, attachment := range message.Attachments {
			uuids = append(uuids, attachment.UUID)
		}
		err := DB.Select(q.In("UUID", uuids)).Find(&uploads)
		if err != nil {
			log.WithField("error", err).Error("failed to find uploads for message")
		}
	}

	return uploads
}

//DeleteUploadsByTicker removes all connected uploads with the given Ticker.
func DeleteUploadsByTicker(ticker *model.Ticker) error {
	err := DB.Select(q.Eq("TickerID", ticker.ID)).Delete(&model.Upload{})
	if err != nil && err.Error() == "not found" {
		return nil
	}

	return err
}

//DeleteUpload remove the given Upload.
func DeleteUpload(upload *model.Upload) error {
	//TODO: Rework with afero.FS from Config
	if err := os.Remove(upload.FullPath()); err != nil {
		log.WithField("error", err).WithField("upload", upload).Error("failed to delete upload file")
	}

	if err := DB.DeleteStruct(upload); err != nil {
		log.WithField("error", err).WithField("upload", upload).Error("failed to delete upload")
		return err
	}

	return nil
}

//DeleteUploads removes a map of Upload.
func DeleteUploads(uploads []*model.Upload) {
	if len(uploads) > 0 {
		for _, upload := range uploads {
			_ = DeleteUpload(upload)
		}
	}
}
