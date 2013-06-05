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
	Provider SocialAccountProvider
	Name     string
	Email    string // if available
}

// WebWallet User (linked to social accounts)
type User struct {
	Id       int // accound inmutable internal id
	Name     string
	Email    string // if available
	Accounts []SocialAccount
}

// WebWallet Address
type Address struct {
	Addr     string // 1... bitcoin formatted address or master address
	Label    string // label as used by bitcoin clients
	Master   bool   // is this a master address, false means a plain?
	walletId Wallet // wallet to which this address belongs to
}

// Wallet
type Wallet struct {
	Id        int       // wallet inmutable internal id
	Name      string    // The address group can have a name, for instance, "Long Term Savings"
	Addr      []Address // addresses on this group
	Balance   float64   // this is not stored but calculated everytime
	Timestamp time.Time // tells the time the balance was calculated
	UserId    int       // User id to which this Wallet belongs to
}
