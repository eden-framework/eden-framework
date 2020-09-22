package apollo

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"go/ast"
)

var (
	Branch           string
	MasterBranch     = "master"
	DefaultNamespace = "application"
)

// 取分支名字
func fetchBranchLastName(branch string) string {
	nameList := strings.Split(branch, "/")
	return nameList[len(nameList)-1]
}

type Initializer interface {
	Init()
}

type Refresher interface {
	Refresh()
}

func InitialConf(conf interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(conf))
	tpe := rv.Type()
	for i := 0; i < tpe.NumField(); i++ {
		value := rv.Field(i)
		if conf, ok := value.Interface().(Initializer); ok {
			conf.Init()
		}
	}
}

func RefreshConf(conf interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(conf))
	tpe := rv.Type()
	for i := 0; i < tpe.NumField(); i++ {
		value := rv.Field(i)
		if conf, ok := value.Interface().(Refresher); ok {
			conf.Refresh()
		}
	}
}

func AssignConf(conf interface{}) {
	AssignConfWithDefault(conf, ApolloBaseConfig{})
}

func AssignConfWithDefault(conf interface{}, defaultBaseConf ApolloBaseConfig) {
	// LOCAL environment doesn't fetch config from apollo
	if os.Getenv("APOLLO_ENV") == "LOCAL" {
		InitialConf(conf)
		return
	}

	if Branch == "" {
		panic(fmt.Sprintf("Namespece[Branch] is empty!"))
	}

	Branch = fetchBranchLastName(Branch)
	if Branch == MasterBranch {
		Branch = DefaultNamespace
	}

	workerCh := make(chan bool)
	apolloConfig := NewApolloConfig(Branch, conf, defaultBaseConf)
	if err := apolloConfig.Start(workerCh); err != nil {
		panic(err)
	}

	InitialConf(conf)

	marshalConfStruct(conf)

	// wait for update action
	go worker(workerCh, conf)
}

func worker(refreshWork chan bool, conf interface{}) {
	for {
		<-refreshWork
		logrus.Infof("worker:%+v", conf)
		RefreshConf(conf)
	}
}

func marshalConfStruct(conf interface{}) {
	rv := reflect.ValueOf(conf)
	if rv.Kind() != reflect.Ptr {
		logrus.Errorf("conf is not a pointer")
		return
	}

	rve := rv.Elem()
	if rve.Kind() != reflect.Struct {
		logrus.Errorf("conf is not a pointer to struct")
		return
	}

	// get a copy of conf
	tmpConf := reflect.New(rve.Type())
	tmpConf.Elem().Set(rve)
	tempConfIntf := tmpConf.Interface()

	// hide sensitive information
	hideSensitiveInfo(tempConfIntf)

	// print apollo config json string
	jsonConf, err := json.Marshal(tempConfIntf)
	if err != nil {
		logrus.Errorf("json marshal err: %v", err)
		return
	}
	logrus.Infof("apollo config json: %s", string(jsonConf))
}

func hideSensitiveInfo(v interface{}) {
	// sensitive keys are case insensitive
	sensitiveKeys := []string{"Password", "Secret"}
	for _, sk := range sensitiveKeys {
		doHideSensitiveInfo(v, sk)
	}
}

// v should be a pointer to struct
func doHideSensitiveInfo(v interface{}, sk string) {
	rve := reflect.ValueOf(v).Elem()
	rt := rve.Type()

	for i := 0; i < rt.NumField(); i++ {
		rtfn := rt.Field(i).Name
		rvf := rve.Field(i)
		if !ast.IsExported(rtfn) {
			continue
		}

		switch rvf.Kind() {
		case reflect.String:
			if strings.Contains(strings.ToLower(rtfn), strings.ToLower(sk)) && rvf.CanSet() {
				rvf.SetString("******")
			}

		case reflect.Ptr:
			rvfe := rvf.Elem()
			if rvfe.Kind() == reflect.Struct && rvf.CanSet() {
				tmpRvf := reflect.New(rvfe.Type())
				tmpRvf.Elem().Set(rvfe)
				doHideSensitiveInfo(tmpRvf.Interface(), sk)
				rvf.Set(tmpRvf)
			}

		case reflect.Struct:
			newRvf := reflect.New(rvf.Type())
			newRvf.Elem().Set(rvf)
			doHideSensitiveInfo(newRvf.Interface(), sk)
			rvf.Set(newRvf.Elem())

		default:

		}
	}
}
