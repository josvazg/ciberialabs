package backend

import (
	/*	"github.com/stretchrcom/goweb"
		"github.com/stretchrcom/goweb/context"*/
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	//goweb.Map("restapi", restapi)
	http.HandleFunc("/user/", userHandler)
}

type SocialAccountProvider int

const (
	None SocialAccountProvider = iota
	Google
	Facebook
	Twitter
)

// Error is the error wrapper
type Error struct {
	Error string
}

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

// add adds or replaces a SocialAccount within User
func (u *User) add(saccount SocialAccount) {
	for i, sa := range u.Accounts {
		if sa.Provider == saccount.Provider && sa.Name == saccount.Name {
			sa.Email = saccount.Email
			u.Accounts[i] = sa
			return
		}
	}
	u.Accounts = append(u.Accounts, saccount)
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
	New(saccount SocialAccount) (*User, error)
	Get(id string) (*User, error) // id=User.Id or SocialAccount.Name
	Edit(u *User) (bool, error)
	Delete(id string) (bool, error) // id=User.Id or  User.Name (only)
}

// DataStoreUserManager implements UserManager on DataStore
type DataStoreUserManager struct {
	c appengine.Context
}

// New creates a new User from a SocialAccount of reference on the DataStore and returns it
func (dsum *DataStoreUserManager) New(saccount SocialAccount) (*User, error) {
	a := &User{Name: saccount.Name,
		Email:    saccount.Email,
		Accounts: []SocialAccount{saccount},
	}
	// create the account on the datastore and generate the auto-id
	k, err := datastore.Put(dsum.c, datastore.NewIncompleteKey(dsum.c, "User", nil), a)
	if err != nil {
		return nil, err
	}
	// store again with the generate id and save it again
	a.Id = k.IntID()
	k, err = datastore.Put(dsum.c, k, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Get returns an User by id, name or SocialAccount.Name
func (dsum *DataStoreUserManager) Get(id string) (u *User, err error) {
	u = &User{}
	intid, err := strconv.ParseInt(id, 10, 64)
	if err == nil && intid != 0 {
		u, err = dsum.getById(intid)
		if err != datastore.ErrNoSuchEntity {
			return
		}
	}
	q := datastore.NewQuery("User").Filter("Name =", id)
	t := q.Run(dsum.c)
	_, err = t.Next(u)
	if err != datastore.Done {
		return
	}
	q = datastore.NewQuery("User").Filter("Email =", id)
	t = q.Run(dsum.c)
	_, err = t.Next(u)
	if err != datastore.Done {
		return
	}
	q = datastore.NewQuery("User").Filter("Accounts.Name =", id)
	t = q.Run(dsum.c)
	_, err = t.Next(u)
	if err != datastore.Done {
		return
	}
	return nil, datastore.ErrNoSuchEntity
}

// getById gets a User by id only
func (dsum *DataStoreUserManager) getById(id int64) (u *User, err error) {
	u = &User{}
	err = datastore.Get(dsum.c, datastore.NewKey(dsum.c, "User", "", id, nil), u)
	return
}

// Edit modifies the user identified by id or name
func (dsum *DataStoreUserManager) Edit(u *User) (bool, error) {
	var oldUser *User
	var err error
	if u.Id != 0 {
		oldUser, err = dsum.getById(u.Id)
	} else {
		oldUser, err = dsum.Get(u.Name)
	}
	if err != nil {
		return false, err
	}
	u.Id = oldUser.Id // preserve Id whatever happens
	_, err = datastore.Put(dsum.c, datastore.NewKey(dsum.c, "User", "", u.Id, nil), u)
	return err == nil, err
}

// Delete removes the user by id or name
func (dsum *DataStoreUserManager) Delete(id string) (bool, error) {
	u, err := dsum.Get(id)
	if err != nil {
		return false, err
	}
	err = datastore.Delete(dsum.c, datastore.NewKey(dsum.c, "User", "", u.Id, nil))
	return err == nil, err
}

// userHandler handles rest calls to user resources
func userHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	dsum := &DataStoreUserManager{c}
	if r.Method == "POST" {
		if strings.HasPrefix(r.URL.Path, "/user/delete/") {
			delUser(dsum, w, r)
		} else {
			newUser(dsum, w, r)
		}
	} else {
		getUser(dsum, w, r)
	}
}

// getUser returns the rest JSON user
func getUser(dsum *DataStoreUserManager, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/user/"):]
	u, e := dsum.Get(id)
	returnJSON(errorOrData(u, e), w)
}

// newUser returns the rest JSON newly created user
func newUser(dsum *DataStoreUserManager, w http.ResponseWriter, r *http.Request) {
	argsPath := r.URL.Path[len("/user/"):]
	args := strings.SplitN(argsPath, "/", 3)
	if len(args) < 2 {
		returnJSON(fmt.Errorf(
			"Can't handle new user request, expecting '{SocialAccountType}/{Name}[/{email}]' but got '%s'!"), w)
		return
	}
	email := ""
	if len(args) > 2 {
		email = args[2]
	}
	saccount, e := buildSocialAccount(args[0], args[1], email)
	if e != nil {
		returnJSON(e, w)
		return
	}
	u, e := dsum.Get(saccount.Name)
	if e != nil && e != datastore.ErrNoSuchEntity {
		returnJSON(e, w)
		return
	}
	if u != nil {
		u.add(saccount)
		_, e = dsum.Edit(u)
	} else {
		u, e = dsum.New(saccount)
	}
	returnJSON(errorOrData(u, e), w)
}

// delUser removes the given user on a rest request returning true or error
func delUser(dsum *DataStoreUserManager, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/user/delete/"):]
	done, e := dsum.Delete(id)
	returnJSON(errorOrData(done, e), w)
}

// errorOrData returns the error (if not null) or the data
func errorOrData(i interface{}, e error) interface{} {
	if e != nil {
		return &Error{e.Error()}
	}
	return i
}

// returnJSON will return i as JSON on w response to request r
func returnJSON(i interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	e := enc.Encode(i)
	if e != nil {
		fmt.Println(w, e)
	}
}

// buildSocialAccount generates a SocialAccount type from the given arguments
func buildSocialAccount(satype, name, email string) (SocialAccount, error) {
	sap := None
	if strings.ToLower(satype) == "google" {
		sap = Google
	} else if strings.ToLower(satype) == "facebook" {
		sap = Facebook
	} else if strings.ToLower(satype) == "twitter" {
		sap = Twitter
	} else {
		return SocialAccount{}, fmt.Errorf("Invalid or Unsupported SocialAccount type %s!", satype)
	}
	return SocialAccount{Provider: sap, Name: name, Email: email}, nil
}
