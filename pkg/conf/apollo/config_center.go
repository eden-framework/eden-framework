package apollo

import (
	"encoding/json"
	"fmt"
	"github.com/profzone/eden-framework/pkg/reflectx"
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

type Refresher interface {
	Refresh()
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

func AssignConfWithDefault(defaultBaseConf ApolloBaseConfig, conf ...interface{}) {
	if Branch == "" {
		panic(fmt.Sprintf("Namespece[Branch] is empty!"))
	}

	Branch = fetchBranchLastName(Branch)
	if Branch == MasterBranch {
		Branch = DefaultNamespace
	}

	workerCh := make(chan bool)
	apolloConfig := NewApolloConfig(Branch, defaultBaseConf, conf...)
	if err := apolloConfig.Start(workerCh); err != nil {
		panic(err)
	}

	marshalConfStruct(conf...)

	// wait for update action
	go worker(workerCh, conf...)
}

func worker(refreshWork chan bool, conf ...interface{}) {
	for {
		<-refreshWork
		logrus.Infof("worker:%+v", conf)

		for _, c := range conf {
			RefreshConf(c)
		}
	}
}

func marshalConfStruct(conf ...interface{}) {
	for _, c := range conf {
		rv := reflect.ValueOf(c)
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
		jsonConf, err := json.MarshalIndent(tempConfIntf, "", "    ")
		if err != nil {
			logrus.Errorf("json marshal err: %v", err)
			return
		}

		logrus.Infof("apollo config json from [%s]: \n%s", reflectx.FullTypeName(reflectx.FromRType(reflect.TypeOf(c))), string(jsonConf))
	}
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
