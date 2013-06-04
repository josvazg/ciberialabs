package backend

import (
	"time"
)

type SocialAccountProvider int

const (
	Google SocialAccountProvider = iota
	Facebook
	Twitter
)

// Social Account from Google, Facebook, Twitter...
type SocialAccount struct {
	provider SocialAccountProvider
	name     string
	email    string // if available
}

// WebWallet Account (linked to social accounts)
type Account struct {
	name   string
	email  string // if available
	social []SocialAccount
}

// WebWallet Address
type Address struct {
	addr    string // 1... bitcoin formatted address or master address
	account string // address name or alias as used by bitcoin clients
	master  bool   // is this a master address, false means a plain?
}

// Address Group
type AddressGroup struct {
	name      string    // The address group can have a name, for instance, "Long Term Savings"
	addr      []Address // addresses on this group
	balance   float64   // this is not stored but calculated everytime
	timestamp time.Time // tells the time the balance was calculated
}
