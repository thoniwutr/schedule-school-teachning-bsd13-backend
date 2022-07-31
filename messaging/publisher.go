package messaging

import (
	"github.com/Beam-Data-Company/merchant-config-svc/server/dto"
)

// Publisher defines an interface for message publisher
type Publisher interface {
	PublishMerchant(merchant *dto.MerchantPublish) error
}
