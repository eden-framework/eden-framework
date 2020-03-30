package test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/profzone/eden-framework/pkg/courier/enumeration"
	"github.com/profzone/eden-framework/pkg/sqlx"
)

var db *sqlx.DB

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	db, _ = sqlx.Open("logger:mysql", "root:root@tcp(localhost:3306)/?charset=utf8&parseTime=true&interpolateParams=true&autocommit=true&loc=Local")
}

func TestUserCRUD(t *testing.T) {
	tt := assert.New(t)

	DBTest.MigrateTo(db, false)
	DBTest.MigrateTo(db, false)

	defer func() {
		err := db.Do(DBTest.Drop()).Err()
		tt.Nil(err)
	}()

	{
		user := User{}
		user.Name = uuid.New().String()

		err := user.Create(db)
		tt.Nil(err)
		tt.Equal(uint64(1), user.ID)

		user.Gender = GenderMale
		{
			err := user.CreateOnDuplicateWithUpdateFields(db, []string{"Gender"})
			tt.Nil(err)
		}
		{
			userForFetch := User{
				Name: user.Name,
			}
			err := userForFetch.FetchByName(db)
			tt.Nil(err)

			tt.Equal(user.Gender, userForFetch.Gender)
		}
		{
			{
				userForDelete := User{
					Name: user.Name,
				}
				err := userForDelete.SoftDeleteByName(db)
				tt.Nil(err)

				userForSelect := User{
					Name: user.Name,
				}
				stmt := userForSelect.T().Select().Where(userForSelect.Fields().Name.Eq(userForSelect.Name))
				errForSelect := db.Do(stmt).Scan(&userForSelect).Err()
				tt.Nil(errForSelect)
				tt.Equal(enumeration.BOOL__FALSE, userForSelect.Enabled)

				{
					err := user.Create(db)
					tt.Nil(err)
					tt.Equal(uint64(3), user.ID)

					userForDelete := User{}
					errForSoftDelete := userForDelete.SoftDeleteByName(db)
					tt.Nil(errForSoftDelete)

					users := []User{}
					fieldsOfUser := userForSelect.Fields()
					stmt := user.T().Select().Where(fieldsOfUser.Enabled.Eq(enumeration.BOOL__FALSE))
					errForSelect := db.Do(stmt).Scan(&users).Err()
					tt.Nil(errForSelect)
					tt.Len(users, 1)
					tt.Equal(uint64(1), users[0].ID)
				}
			}
		}
	}
}

func TestUserList(t *testing.T) {
	tt := assert.New(t)

	DBTest.MigrateTo(db, false)
	defer func() {
		err := db.Do(DBTest.Drop()).Err()
		tt.Nil(err)
	}()

	createUser := func() {
		user := User{}
		user.Name = uuid.New().String()
		err := user.Create(db)
		tt.Nil(err)
	}

	for i := 0; i < 10; i++ {
		createUser()
	}

	list := UserList{}
	count, err := list.FetchList(db, -1, -1)
	tt.Nil(err)
	tt.Equal(int32(10), count)
	tt.Len(list, 10)

	names := []string{}
	for _, user := range list {
		names = append(names, user.Name)
	}

	{

		list2 := UserList{}
		err := list2.BatchFetchByNameList(db, names)
		tt.Nil(err)
		tt.Len(list, 10)
	}

}
