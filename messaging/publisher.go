package messaging

import (
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
)

// Publisher defines an interface for message publisher
type Publisher interface {
	PublishMerchant(merchant *dto.MerchantPublish) error
}
