package predicates

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func HasLocation() th.Predicate {
	return func(update telego.Update) bool {
		if update.Message != nil {
			return update.Message.Location != nil
		}
		return false
	}
}
