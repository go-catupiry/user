package user

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-catupiry/catu"
	"github.com/go-catupiry/catu/helpers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserModel struct {
	ID uint64 `gorm:"primary_key;column:id;" json:"id"`

	Username string `gorm:"unique;column:username;" json:"username"`
	Email    string `gorm:"unique;column:email;" json:"email"`

	DisplayName string `gorm:"column:displayName;" json:"displayName"`
	FullName    string `gorm:"column:fullName;" json:"fullName"`
	Biography   string `gorm:"column:biography;type:TEXT;" json:"biography"`
	Gender      string `gorm:"column:gender;" json:"gender"`

	Active  bool `gorm:"column:active;" json:"active"`
	Blocked bool `gorm:"column:blocked;" json:"blocked"`

	Language     string `gorm:"column:language;" json:"language"`
	ConfirmEmail string `gorm:"column:confirmEmail;" json:"confirmEmail"`

	AcceptTerms bool   `gorm:"column:acceptTerms;" json:"acceptTerms"`
	Birthdate   string `gorm:"column:birthdate;" json:"birthdate"`
	Phone       string `gorm:"column:phone;" json:"phone"`

	Roles     []string `gorm:"-" json:"roles"`
	RolesText string   `gorm:"column:roles;" json:"-"`

	CreatedAt time.Time `gorm:"column:createdAt;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;" json:"updatedAt"`
}

func (r *UserModel) GetID() string {
	return strconv.FormatUint(r.ID, 10)
}

func (r *UserModel) SetID(id string) error {
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	r.ID = n

	return nil
}

func (r *UserModel) SetRoles(v []string) error {
	r.Roles = v
	return nil
}

func (r *UserModel) AddRole(roleName string) error {
	roles := r.GetRoles()
	roles = append(roles, roleName)

	jsonString, _ := json.Marshal(roles)
	r.RolesText = string(jsonString)

	r.Roles = roles

	rolesByte, _ := json.Marshal(&r.Roles)
	r.RolesText = string(rolesByte)

	return nil
}

func (r *UserModel) RemoveRole(role string) error {
	// r.Roles.
	r.Roles, _ = helpers.SliceRemove(r.Roles, role)
	return nil
}

func (r *UserModel) GetEmail() string {
	return r.Email
}

func (r *UserModel) SetEmail(v string) error {
	// TODO! Validate email format!
	r.Email = v
	return nil
}

func (r *UserModel) GetUsername() string {
	return r.Username
}

func (r *UserModel) SetUsername(v string) error {
	r.Username = v
	return nil
}

func (r *UserModel) SetDisplayName(v string) error {
	r.DisplayName = v
	return nil
}

func (r *UserModel) SetFullName(v string) error {
	r.FullName = v
	return nil
}

func (r *UserModel) GetLanguage() string {
	return r.Language
}

func (r *UserModel) SetLanguage(v string) error {
	// TODO! Validate if this land is valid
	r.Language = v
	return nil
}

func (r *UserModel) IsActive() bool {
	return r.Active
}

func (r *UserModel) SetActive(v bool) error {
	r.Active = v
	return nil
}

func (r *UserModel) SetBlocked(blocked bool) error {
	r.Blocked = blocked
	return nil
}

func (UserModel) TableName() string {
	return "users"
}

func (r *UserModel) FillById(id string) error {
	return UserFindOne(id, r)
}

func (r *UserModel) GetRoles() []string {
	if r.RolesText != "" {
		_ = json.Unmarshal([]byte(r.RolesText), &r.Roles)
	}

	return r.Roles
}

func (r *UserModel) GetDisplayName() string {
	return r.DisplayName
}

func (r *UserModel) GetFullName() string {
	return r.FullName
}

func (r *UserModel) IsBlocked() bool {
	return r.Blocked
}

func (r *UserModel) GetBiography() string {
	return r.Biography
}

func (r *UserModel) GetGender() string {
	return r.Gender
}

func (r *UserModel) GetActiveString() string {
	return strconv.FormatBool(r.Active)
}

func (r *UserModel) GetBlockedString() string {
	return strconv.FormatBool(r.Blocked)
}

func (r *UserModel) GetAcceptTermsString() string {
	return strconv.FormatBool(r.AcceptTerms)
}

func (r *UserModel) GetBirthdate() string {
	return r.Birthdate
}

func (r *UserModel) GetPhone() string {
	return r.Phone
}

func (r *UserModel) GetCreatedAtString() string {
	return r.CreatedAt.UTC().String()
}

func (r *UserModel) GetUpdateAtString() string {
	return r.UpdatedAt.UTC().String()
}

func (m *UserModel) Save() error {
	var err error
	db := catu.GetDefaultDatabaseConnection()

	if m.ID == 0 {
		// create ....
		err = db.Create(&m).Error
		if err != nil {
			return err
		}
	} else {
		// update ...
		err = db.Save(&m).Error
		if err != nil {
			return err
		}
	}

	// TODO! re-set url alias

	return nil
}

func (m *UserModel) LoadTeaserData() error {
	m.GetRoles()
	return nil
}

func (m *UserModel) LoadData() error {
	m.GetRoles()
	return nil
}

func (r *UserModel) Delete() error {
	db := catu.GetDefaultDatabaseConnection()
	return db.Unscoped().Delete(&r).Error
}

func UsersQuery(userList *[]UserModel, limit int) error {

	db := catu.GetDefaultDatabaseConnection()

	if err := db.
		Limit(limit).
		Find(userList).Error; err != nil {
		return err
	}
	return nil
}

// FindOne - Find one user record
func UserFindOne(id string, record *UserModel) error {
	db := catu.GetDefaultDatabaseConnection()

	return db.First(record, id).Error
}

func UserFindOneByUsername(username string, record *UserModel) error {
	db := catu.GetDefaultDatabaseConnection()

	return db.
		Where(
			db.Where("username = ?", username).
				Or(db.Where("email = ?", username)),
		).
		First(record).Error
}

func UserFindOneByEmail(email string, record *UserModel) error {
	db := catu.GetDefaultDatabaseConnection()

	return db.
		Where("email = ?", email).
		First(record).Error
}

func LoadAllUsers(userList *[]UserModel) error {
	db := catu.GetDefaultDatabaseConnection()

	if err := db.
		Limit(99999).
		Order("displayName ASC, id ASC").
		Find(userList).Error; err != nil {
		return err
	}
	return nil

}

type QueryAndCountFromRequestCfg struct {
	Records *[]*UserModel
	Count   *int64
	Limit   int
	Offset  int
	C       echo.Context
	IsHTML  bool
}

func QueryAndCountFromRequest(opts *QueryAndCountFromRequestCfg) error {
	db := catu.GetDefaultDatabaseConnection()

	c := opts.C

	q := c.QueryParam("q")
	query := db
	ctx := c.(*catu.RequestContext)

	can := ctx.Can("find_user")
	if !can {
		return nil
	}

	queryI, err := ctx.Query.SetDatabaseQueryForModel(query, &UserModel{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("QueryAndCountFromRequest error")
	}
	query = queryI.(*gorm.DB)

	if q != "" {
		query = query.Where(
			db.Where("displayName LIKE ?", "%"+q+"%").Or(db.Where("fullName LIKE ?", "%"+q+"%")),
		)
	}

	orderColumn, orderIsDesc, orderValid := helpers.ParseUrlQueryOrder(c.QueryParam("order"))

	if orderValid {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Table: clause.CurrentTable, Name: orderColumn},
			Desc:   orderIsDesc,
		})
	} else {
		query = query.Order("createdAt DESC").
			Order("id DESC")
	}

	query = query.Limit(opts.Limit).
		Offset(opts.Offset)

	r := query.Find(opts.Records)
	if r.Error != nil {
		return r.Error
	}

	return CountQueryFromRequest(opts)
}

func CountQueryFromRequest(opts *QueryAndCountFromRequestCfg) error {
	db := catu.GetDefaultDatabaseConnection()

	c := opts.C
	q := c.QueryParam("q")
	ctx := c.(*catu.RequestContext)

	// Count ...
	queryCount := db

	if q != "" {
		queryCount = queryCount.Or(
			db.Where("displayName LIKE ?", "%"+q+"%"),
			db.Where("fullName LIKE ?", "%"+q+"%"),
		)
	}

	can := ctx.Can("find_user")
	if !can {
		return nil
	}

	queryICount, err := ctx.Query.SetDatabaseQueryForModel(queryCount, &UserModel{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("QueryAndCountFromRequest count error")
	}
	queryCount = queryICount.(*gorm.DB)

	return queryCount.
		Table("users").
		Count(opts.Count).Error
}

type UserModelOpts struct {
}

func NewUserModel(opts *UserModelOpts) (*UserModel, error) {
	return &UserModel{}, nil
}
