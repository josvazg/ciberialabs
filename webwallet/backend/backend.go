package backend

import (
	/*	"github.com/stretchrcom/goweb"
		"github.com/stretchrcom/goweb/context"*/
	"appengine"
	"appengine/datastore"
	"fmt"
	"strings"
	"time"
)

func init() {
	//goweb.Map("restapi", restapi)
}

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
	Id       int64 // accound inmutable internal id
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

//func restApi(c context.Context) error {
//return goweb.Respond.With(c, 200, []byte("Welcome to the Goweb restapi!"))
//}
/*
Cuentas:
(POST) Alta dados los datos de una cuenta social (twitter, facebook, google, etc) [Tu te encargas en el navegador de obtener los datos de la otra cuenta y a cambio no te tienes que preocupar de datastores ni nada de eso]
(GET) Consulta dado el login de otra cuenta social.
(PUT?) Modificación, añadir o quitar cuentas sociales y modificar el email de contacto.
(DELETE) Baja y borrado de una cuenta y sus datos asociados.
Direcciones (normales y maestras ¿debe el backend distinguirlas? creo que para javascript tenemos el código para manejar las direcciones de carbonwallet ¿no?):
(POST) Alta de una dirección (o varias) en una cuenta dada.
(GET) Consulta de direcciones de una cuenta.
(DELETE) Baja de una dirección (o varias) de una cuenta.
Saldos:
(GET) Saldo de una cuenta o de un conjunto de direcciones.
*/

// UserManager contains the operations on Users
type UserManager interface {
	New(saccount *SocialAccount) (*Account, error)
	Get(id string) (*Account, error) // id=Account.Id o SocialAccount.Name
	Edit(account *Account) (bool, error)
	Delete(id string) (bool, error) // id=Account.Id (only)
}

// DataStoreUserManager implements UserManager on DataStore
type DataStoreUserManager struct {
	c appengine.Context
}

// New creates a new Account from a SocialAccount of reference on the DataStore and returns it
func (dsum *DataStoreUserManager) New(saccount *SocialAccount) (*Account, error) {
	if saccount == nil {
		return nil, fmt.Errorf("Need a non nil SocialAccount!")
	}
	a := &Account{Name: saccount.Name,
		Email:    saccount.Email,
		Accounts: []SocialAccount{saccount},
	}
	k := datastore.NewIncompleteKey(dsum.c, "Account", nil)
	// create the account on the datastore and generate the auto-id
	k, err := datastore.Put(dsum.c, k, a)
	if !err {
		return nil, err
	}
	// store again with the generate id and save it again
	a.Id = k.IntID()
	k, err = datastore.Put(dsum.c, k, a)
	if !err {
		return nil, err
	}
	return a
}

// Get returns an Account by id, name or SocialAccount.Name
func (dsum *DataStoreUserManager) Get(id string) (*Account, error) {
	a := &Account{}
	intid := toint64(id)
	if intid != nil {
		err := datastore.Get(dsum.c, datastore.NewKey(dsum.c, "Account", nil, intid, nil))
		if !err {
			return nil, err
		}
	}
}
