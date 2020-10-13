package files

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/eden-framework/internal/generator/api"
	"github.com/eden-framework/eden-framework/internal/generator/importer"
	str "github.com/eden-framework/strings"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type ClientFile struct {
	ClientName  string
	PackageName string
	Name        string
	Importer    *importer.PackageImporter
	a           *api.Api
}

func NewClientFile(name string, a *api.Api) *ClientFile {
	return &ClientFile{
		Name:        str.ToLowerLinkCase(name),
		PackageName: str.ToLowerSnakeCase("client-" + name),
		ClientName:  str.ToUpperCamelCase("client-" + name),
		Importer:    importer.NewPackageImporter(""),
		a:           a,
	}
}

func (c *ClientFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", c.PackageName))
	return
}

func (c *ClientFile) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, c.Importer.String())
	return
}

func (c *ClientFile) WriteTypeInterface(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("type %sInterface interface {\n", c.ClientName))
	if err != nil {
		return err
	}
	for groupName, group := range c.a.Operators {
		for methodName, method := range group.Methods {
			req, resp := make([]string, 0), make([]string, 0)
			var typeName string
			for _, modelName := range method.Inputs {
				model, ok := c.a.Models[modelName]
				if !ok {
					logrus.Panic(fmt.Errorf("%s not exist in model definations", modelName))
				}
				if model.NeedAlias {
					typeName = c.Importer.Use(modelName)
				} else {
					typeName = model.Name
				}
				req = append(req, fmt.Sprintf("%s *%s", str.ToLowerCamelCase(model.Name), typeName))
			}
			for _, modelName := range method.Outputs {
				model, ok := c.a.Models[modelName]
				if !ok {
					logrus.Panic(fmt.Errorf("%s not exist in model definations", modelName))
				}
				if model.NeedAlias {
					typeName = c.Importer.Use(modelName)
				} else {
					typeName = model.Name
				}
				resp = append(resp, fmt.Sprintf("%s *%s", str.ToLowerCamelCase(model.Name), typeName))
			}
			resp = append(resp, "err error")
			methodString := fmt.Sprintf("%s(%s) (%s)\n", str.ToUpperCamelCase(groupName+methodName), strings.Join(req, ", "), strings.Join(resp, ", "))
			_, err = io.WriteString(w, methodString)
			if err != nil {
				return err
			}
		}
	}

	_, err = io.WriteString(w, "}\n")
	return
}

func (c *ClientFile) WriteClientGeneral(w io.Writer) (err error) {
	typeDef := `type ` + c.ClientName + ` struct {
	` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.PeerConfig") + `
	peer       ` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.Peer") + `
	session    ` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.Session") + `
	RemoteAddr string` + " `yaml:\"remoteAddr\" ini:\"remoteAddr\" comment:\"Remote address with port\"`" + `
}

func (c *` + c.ClientName + `) Init() {
	c.peer = ` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.NewPeer") + `(c.PeerConfig)

	var stat *` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.Status") + `
	c.session, stat = c.peer.Dial(c.RemoteAddr)
	if !stat.OK() {
		panic(` + c.Importer.Use("fmt.Errorf") + `("connection err, status: %v", stat.String()))
	}
}

func (c *` + c.ClientName + `) DockerDefaults() map[string]string {
	return map[string]string{}
}

func (c *` + c.ClientName + `) RegisterCallRouter(route interface{}, plugins ...` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.Plugin") + `) []string {
	return c.peer.RouteCall(route, plugins...)
}

func (c *` + c.ClientName + `) RegisterPushRouter(route interface{}, plugins ...` + c.Importer.Use("github.com/henrylee2cn/erpc/v6.Plugin") + `) []string {
	return c.peer.RoutePush(route, plugins...)
}

`

	_, err = io.WriteString(w, typeDef)
	return
}

func (c *ClientFile) WriteMethods(w io.Writer) (err error) {
	var methodsDef string
	c.a.WalkOperators(func(g *api.OperatorGroup) {
		g.WalkMethods(func(m *api.OperatorMethod) {
			req, resp := make([][]string, 0), make([][]string, 0)
			var typeName string
			m.WalkInputs(func(i string) {
				model, ok := c.a.Models[i]
				if !ok {
					logrus.Panic(fmt.Errorf("%s not exist in model definations", i))
				}
				if model.NeedAlias {
					typeName = c.Importer.Use(i)
				} else {
					typeName = model.Name
				}
				req = append(req, []string{str.ToLowerCamelCase(model.Name), "*" + typeName})
			})
			m.WalkOutputs(func(i string) {
				model, ok := c.a.Models[i]
				if !ok {
					logrus.Panic(fmt.Errorf("%s not exist in model definations", i))
				}
				if model.NeedAlias {
					typeName = c.Importer.Use(i)
				} else {
					typeName = model.Name
				}
				resp = append(resp, []string{str.ToLowerCamelCase(model.Name), "*" + typeName})
			})
			resp = append(resp, []string{"err", "error"})

			var requestStr string
			if len(req) == 0 {
				requestStr = "nil"
			} else {
				requestStr = req[0][0]
			}
			if !g.IsPush {
				methodsDef += `func (c *` + c.ClientName + `) ` + strings.Join([]string{g.Name, m.Name}, "") + `(` + str.RecursiveJoin(req, " ", ", ") + `) (` + str.RecursiveJoin(resp, " ", ", ") + `) {
	stat := c.session.Call("` + strings.Join([]string{g.Path, m.Path}, "") + `", ` + requestStr + `, &` + resp[0][0] + `).Status()
	if !stat.OK() {
		err = stat.Cause()
	}
	return
}

`
			} else {
				methodsDef += `func (c *` + c.ClientName + `) ` + strings.Join([]string{g.Name, m.Name}, "") + `(` + str.RecursiveJoin(req, " ", ", ") + `) (err error) {
	stat := c.session.Push("` + strings.Join([]string{g.Path, m.Path}, "") + `", ` + requestStr + `)
	if !stat.OK() {
		err = stat.Cause()
	}
	return
}

`
			}
		})
	})

	_, err = io.WriteString(w, methodsDef)
	return
}

func (c *ClientFile) WriteAll() string {
	w := bytes.NewBuffer([]byte{})
	err := c.WriteTypeInterface(w)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteClientGeneral(w)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteMethods(w)
	if err != nil {
		logrus.Panic(err)
	}

	return w.String()
}

func (c *ClientFile) String() string {
	buf := bytes.NewBuffer([]byte{})

	content := c.WriteAll()

	err := c.WritePackage(buf)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteImports(buf)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = io.WriteString(buf, content)
	if err != nil {
		logrus.Panic(err)
	}

	return buf.String()
}
