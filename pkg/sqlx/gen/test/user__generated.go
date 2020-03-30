package test

import (
	fmt "fmt"
	time "time"

	golib_tools_courier_enumeration "github.com/profzone/eden-framework/pkg/courier/enumeration"
	golib_tools_sqlx "github.com/profzone/eden-framework/pkg/sqlx"
	golib_tools_sqlx_builder "github.com/profzone/eden-framework/pkg/sqlx/builder"
	golib_tools_timelib "github.com/profzone/eden-framework/pkg/timelib"
)

var UserTable *golib_tools_sqlx_builder.Table

func init() {
	UserTable = DBTest.Register(&User{})
}

func (user *User) D() *golib_tools_sqlx.Database {
	return DBTest
}

func (user *User) T() *golib_tools_sqlx_builder.Table {
	return UserTable
}

func (user *User) TableName() string {
	return "t_user"
}

type UserFields struct {
	Name       *golib_tools_sqlx_builder.Column
	Username   *golib_tools_sqlx_builder.Column
	Nickname   *golib_tools_sqlx_builder.Column
	Gender     *golib_tools_sqlx_builder.Column
	Birthday   *golib_tools_sqlx_builder.Column
	Boolean    *golib_tools_sqlx_builder.Column
	CreateTime *golib_tools_sqlx_builder.Column
	UpdateTime *golib_tools_sqlx_builder.Column
	ID         *golib_tools_sqlx_builder.Column
	Enabled    *golib_tools_sqlx_builder.Column
}

var UserField = struct {
	Name       string
	Username   string
	Nickname   string
	Gender     string
	Birthday   string
	Boolean    string
	CreateTime string
	UpdateTime string
	ID         string
	Enabled    string
}{
	Name:       "Name",
	Username:   "Username",
	Nickname:   "Nickname",
	Gender:     "Gender",
	Birthday:   "Birthday",
	Boolean:    "Boolean",
	CreateTime: "CreateTime",
	UpdateTime: "UpdateTime",
	ID:         "ID",
	Enabled:    "Enabled",
}

func (user *User) Fields() *UserFields {
	table := user.T()

	return &UserFields{
		Name:       table.F(UserField.Name),
		Username:   table.F(UserField.Username),
		Nickname:   table.F(UserField.Nickname),
		Gender:     table.F(UserField.Gender),
		Birthday:   table.F(UserField.Birthday),
		Boolean:    table.F(UserField.Boolean),
		CreateTime: table.F(UserField.CreateTime),
		UpdateTime: table.F(UserField.UpdateTime),
		ID:         table.F(UserField.ID),
		Enabled:    table.F(UserField.Enabled),
	}
}

func (user *User) IndexFieldNames() []string {
	return []string{"ID", "Name", "Nickname", "Username"}
}

func (user *User) ConditionByStruct() *golib_tools_sqlx_builder.Condition {
	table := user.T()

	fieldValues := golib_tools_sqlx.FieldValuesFromStructByNonZero(user)

	conditions := []*golib_tools_sqlx_builder.Condition{}

	for _, fieldName := range user.IndexFieldNames() {
		if v, exists := fieldValues[fieldName]; exists {
			conditions = append(conditions, table.F(fieldName).Eq(v))
			delete(fieldValues, fieldName)
		}
	}

	if len(conditions) == 0 {
		panic(fmt.Errorf("at least one of field for indexes has value"))
	}

	for fieldName, v := range fieldValues {
		conditions = append(conditions, table.F(fieldName).Eq(v))
	}

	condition := golib_tools_sqlx_builder.And(conditions...)

	condition = golib_tools_sqlx_builder.And(condition, table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	return condition
}

func (user *User) PrimaryKey() golib_tools_sqlx.FieldNames {
	return golib_tools_sqlx.FieldNames{"ID"}
}
func (user *User) Indexes() golib_tools_sqlx.Indexes {
	return golib_tools_sqlx.Indexes{
		"I_nickname": golib_tools_sqlx.FieldNames{"Nickname"},
		"I_username": golib_tools_sqlx.FieldNames{"Username"},
	}
}
func (user *User) UniqueIndexes() golib_tools_sqlx.Indexes {
	return golib_tools_sqlx.Indexes{"I_name": golib_tools_sqlx.FieldNames{"Name", "Enabled"}}
}

func (user *User) Create(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	if user.CreateTime.IsZero() {
		user.CreateTime = golib_tools_timelib.MySQLTimestamp(time.Now())
	}
	user.UpdateTime = user.CreateTime

	stmt := user.D().
		Insert(user).
		Comment("User.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		user.ID = uint64(lastInsertID)
	}

	return err
}

func (user *User) DeleteByStruct(db *golib_tools_sqlx.DB) (err error) {
	table := user.T()

	stmt := table.Delete().
		Comment("User.DeleteByStruct").
		Where(user.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (user *User) CreateOnDuplicateWithUpdateFields(db *golib_tools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	if user.CreateTime.IsZero() {
		user.CreateTime = golib_tools_timelib.MySQLTimestamp(time.Now())
	}
	user.UpdateTime = user.CreateTime

	table := user.T()

	fieldValues := golib_tools_sqlx.FieldValuesFromStructByNonZero(user, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range user.UniqueIndexes() {
		for _, field := range fieldNames {
			delete(m, field)
		}
	}

	if len(m) == 0 {
		panic(fmt.Errorf("no fields for updates"))
	}

	for field := range fieldValues {
		if !m[field] {
			delete(fieldValues, field)
		}
	}

	stmt := table.
		Insert().Columns(cols).Values(vals...).
		OnDuplicateKeyUpdate(table.AssignsByFieldValues(fieldValues)...).
		Comment("User.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (user *User) FetchByID(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByID").
		Where(golib_tools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) FetchByIDForUpdate(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByIDForUpdate").
		Where(golib_tools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) DeleteByID(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Delete().
		Comment("User.DeleteByID").
		Where(golib_tools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) UpdateByIDWithMap(db *golib_tools_sqlx.DB, fieldValues golib_tools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = golib_tools_timelib.MySQLTimestamp(time.Now())
	}

	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("User.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(golib_tools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return user.FetchByID(db)
	}
	return nil
}

func (user *User) UpdateByIDWithStruct(db *golib_tools_sqlx.DB, zeroFields ...string) error {
	fieldValues := golib_tools_sqlx.FieldValuesFromStructByNonZero(user, zeroFields...)
	return user.UpdateByIDWithMap(db, fieldValues)
}

func (user *User) SoftDeleteByID(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()

	fieldValues := golib_tools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = golib_tools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = golib_tools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("User.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(golib_tools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		dbErr := golib_tools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return user.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (user *User) FetchByName(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByName").
		Where(golib_tools_sqlx_builder.And(
			table.F("Name").Eq(user.Name),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) FetchByNameForUpdate(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByNameForUpdate").
		Where(golib_tools_sqlx_builder.And(
			table.F("Name").Eq(user.Name),
			table.F("Enabled").Eq(user.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) DeleteByName(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Delete().
		Comment("User.DeleteByName").
		Where(golib_tools_sqlx_builder.And(
			table.F("Name").Eq(user.Name),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) UpdateByNameWithMap(db *golib_tools_sqlx.DB, fieldValues golib_tools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = golib_tools_timelib.MySQLTimestamp(time.Now())
	}

	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("User.UpdateByNameWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(golib_tools_sqlx_builder.And(
			table.F("Name").Eq(user.Name),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return user.FetchByName(db)
	}
	return nil
}

func (user *User) UpdateByNameWithStruct(db *golib_tools_sqlx.DB, zeroFields ...string) error {
	fieldValues := golib_tools_sqlx.FieldValuesFromStructByNonZero(user, zeroFields...)
	return user.UpdateByNameWithMap(db, fieldValues)
}

func (user *User) SoftDeleteByName(db *golib_tools_sqlx.DB) error {
	user.Enabled = golib_tools_courier_enumeration.BOOL__TRUE

	table := user.T()

	fieldValues := golib_tools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = golib_tools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = golib_tools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("User.SoftDeleteByName").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(golib_tools_sqlx_builder.And(
			table.F("Name").Eq(user.Name),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		dbErr := golib_tools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return user.DeleteByName(db)
		}
		return err
	}
	return nil
}

type UserList []User

// deprecated
func (userList *UserList) FetchList(db *golib_tools_sqlx.DB, size int32, offset int32, conditions ...*golib_tools_sqlx_builder.Condition) (count int32, err error) {
	*userList, count, err = (&User{}).FetchList(db, size, offset, conditions...)
	return
}

func (user *User) FetchList(db *golib_tools_sqlx.DB, size int32, offset int32, conditions ...*golib_tools_sqlx_builder.Condition) (userList UserList, count int32, err error) {
	userList = UserList{}

	table := user.T()

	condition := golib_tools_sqlx_builder.And(conditions...)

	condition = golib_tools_sqlx_builder.And(condition, table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(golib_tools_sqlx_builder.Count(golib_tools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

func (user *User) List(db *golib_tools_sqlx.DB, condition *golib_tools_sqlx_builder.Condition) (userList UserList, err error) {
	userList = UserList{}

	table := user.T()

	condition = golib_tools_sqlx_builder.And(condition, table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.List").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

func (user *User) ListByStruct(db *golib_tools_sqlx.DB) (userList UserList, err error) {
	userList = UserList{}

	table := user.T()

	condition := user.ConditionByStruct()

	condition = golib_tools_sqlx_builder.And(condition, table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

// deprecated
func (userList *UserList) BatchFetchByIDList(db *golib_tools_sqlx.DB, idList []uint64) (err error) {
	*userList, err = (&User{}).BatchFetchByIDList(db, idList)
	return
}

func (user *User) BatchFetchByIDList(db *golib_tools_sqlx.DB, idList []uint64) (userList UserList, err error) {
	if len(idList) == 0 {
		return UserList{}, nil
	}

	table := user.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

// deprecated
func (userList *UserList) BatchFetchByNameList(db *golib_tools_sqlx.DB, nameList []string) (err error) {
	*userList, err = (&User{}).BatchFetchByNameList(db, nameList)
	return
}

func (user *User) BatchFetchByNameList(db *golib_tools_sqlx.DB, nameList []string) (userList UserList, err error) {
	if len(nameList) == 0 {
		return UserList{}, nil
	}

	table := user.T()

	condition := table.F("Name").In(nameList)

	condition = condition.And(table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.BatchFetchByNameList").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

// deprecated
func (userList *UserList) BatchFetchByNicknameList(db *golib_tools_sqlx.DB, nicknameList []string) (err error) {
	*userList, err = (&User{}).BatchFetchByNicknameList(db, nicknameList)
	return
}

func (user *User) BatchFetchByNicknameList(db *golib_tools_sqlx.DB, nicknameList []string) (userList UserList, err error) {
	if len(nicknameList) == 0 {
		return UserList{}, nil
	}

	table := user.T()

	condition := table.F("Nickname").In(nicknameList)

	condition = condition.And(table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.BatchFetchByNicknameList").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

// deprecated
func (userList *UserList) BatchFetchByUsernameList(db *golib_tools_sqlx.DB, usernameList []string) (err error) {
	*userList, err = (&User{}).BatchFetchByUsernameList(db, usernameList)
	return
}

func (user *User) BatchFetchByUsernameList(db *golib_tools_sqlx.DB, usernameList []string) (userList UserList, err error) {
	if len(usernameList) == 0 {
		return UserList{}, nil
	}

	table := user.T()

	condition := table.F("Username").In(usernameList)

	condition = condition.And(table.F("Enabled").Eq(golib_tools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.BatchFetchByUsernameList").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}
