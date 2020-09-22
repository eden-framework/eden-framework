package apollo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func dirExist(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	} else {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
}

//write config to file
func writeConfigFile(content []byte, confFullPath string) error {
	dir := filepath.Dir(confFullPath)
	exist, err := dirExist(dir)
	if err != nil {
		logrus.Errorf("dirExist failed![%+v][%+v]", confFullPath, err)
		return err
	}

	// create directory
	if !exist {
		if err := os.MkdirAll(dir, 0777); err != nil {
			logrus.Errorf("Mkdir failed![%+v][%+v]", dir, err)
			return err
		}
	}

	if err := ioutil.WriteFile(confFullPath, content, 0644); err != nil {
		logrus.Errorf("writeConfigFile fail,error:%+v", err)
		return err
	}

	return nil
}

//load config from file
func loadConfigFile(confFullPath string) ([]byte, error) {
	content, err := ioutil.ReadFile(confFullPath)
	if err != nil {
		logrus.Errorf("loadConfigFile fail,error:%+v", err)
		return nil, err
	}

	return content, nil
}
