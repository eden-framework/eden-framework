package generator

import (
	"fmt"
	"github.com/profzone/eden-framework/internal/generator/files"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

type ServiceOption struct {
	Name            string
	PackageName     string      `survey:"package_name"`
	DatabaseSupport expressBool `survey:"database_support"`
	ApolloSupport   expressBool `survey:"apollo_support"`

	Group           string
	Owner           string
	Desc            string
	Version         string
	ProgramLanguage string `survey:"project_language"`
	Workflow        string
}

type ServiceGenerator struct {
	opt ServiceOption
}

func NewServiceGenerator(opt ServiceOption) *ServiceGenerator {
	s := &ServiceGenerator{
		opt: opt,
	}

	return s
}

func (s *ServiceGenerator) Load(path string) {
}

func (s *ServiceGenerator) Pick() {
}

func (s *ServiceGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	// create service directory
	p := path.Join(outputPath, s.opt.Name)
	createPath(p)

	// go.mod file
	mod := files.NewModFile(s.opt.PackageName, "1.14")
	mod.AddReplace("k8s.io/client-go", "", "k8s.io/client-go", "v0.18.8")
	mod.AddRequired("github.com/profzone/eden-framework", "v1.0.9")
	mod.AddRequired("github.com/sirupsen/logrus", "v1.6.0")
	mod.AddRequired("github.com/spf13/cobra", "v0.0.5")
	outputs.WriteFile(path.Join(p, "go.mod"), mod.String())

	err := os.Chdir(p)
	if err != nil {
		logrus.Panicf("os.Chdir failed: %v", err)
	}

	// apollo config file
	if s.opt.ApolloSupport {
		apolloFile := s.createApolloFile()
		outputs.WriteFile(path.Join(p, "internal/global/apollo.go"), apolloFile.String())
	}

	// db config file
	if s.opt.DatabaseSupport {
		dbFile := s.createDbConfigFile()
		outputs.WriteFile(path.Join(p, "internal/databases/db.go"), dbFile.String())
	}

	// general config file
	configFile := s.createConfigFile(p)
	outputs.WriteFile(path.Join(p, "internal/global/config.go"), configFile.String())

	// router v0 root files
	routerV0RootFile := s.createRouterV0RootFile()
	outputs.WriteFile(path.Join(p, "internal/routers/v0/root.go"), routerV0RootFile.String())

	// router root files
	routerRootFile := s.createRouterRootFile(p)
	outputs.WriteFile(path.Join(p, "internal/routers/root.go"), routerRootFile.String())

	// main file
	mainFile := s.createMainFile(p)
	outputs.Add(path.Join(p, "cmd/main.go"), mainFile.String())

	return outputs
}

func createPath(p string) {
	if !PathExist(p) {
		err := os.Mkdir(p, 0755)
		if err != nil {
			logrus.Panicf("os.Mkdir failed: %v, path: %s", err, p)
		}
		return
	}
	logrus.Panicf("os.Stat exist: %s", p)
}

func (s *ServiceGenerator) createApolloFile() *files.GoFile {
	file := files.NewGoFile("global")
	file.WithBlock(fmt.Sprintf(`
var ApolloConfig = {{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/conf/apollo" "" }}.ApolloBaseConfig{
	AppId:            "%s",
	Host:             "localhost:8080",
	BackupConfigPath: "./apollo_config",
	Cluster:          "default",
}
`, s.opt.Name))

	return file
}

func (s *ServiceGenerator) createDbConfigFile() *files.GoFile {
	file := files.NewGoFile("databases")
	file.WithBlock(`
var Config = struct {
	DBTest *{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/sqlx" "" }}.Database
}{
	DBTest: &{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/sqlx" "" }}.Database{},
}
`)

	return file
}

func (s *ServiceGenerator) createConfigFile(cwd string) *files.GoFile {
	file := files.NewGoFile("global")

	file.WithBlock(`
var Config = struct {
	LogLevel {{ .UseWithoutAlias "github.com/sirupsen/logrus" "" }}.Level
`)
	if s.opt.DatabaseSupport {
		file.WithBlock(`
	// db
	MasterDB *{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/client/mysql" "" }}.MySQL
	SlaveDB  *{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/client/mysql" "" }}.MySQL
`)
	}
	file.WithBlock(`
	// administrator
	GRPCServer *{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier/transport_grpc" "" }}.ServeGRPC
	HTTPServer *{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier/transport_http" "" }}.ServeHTTP
}{
	LogLevel: {{ .UseWithoutAlias "github.com/sirupsen/logrus" "" }}.DebugLevel,
`)
	if s.opt.DatabaseSupport {
		dbUse := path.Join(s.opt.PackageName, "internal/databases")
		dbPath := path.Join(cwd, "internal/databases")
		file.WithBlock(fmt.Sprintf(`
	MasterDB: &{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/client/mysql" "" }}.MySQL{Database: {{ .UseWithoutAlias "%s" "%s" }}.Config.DBTest},
	SlaveDB:  &{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/client/mysql" "" }}.MySQL{Database: {{ .UseWithoutAlias "%s" "%s" }}.Config.DBTest},
`, dbUse, dbPath, dbUse, dbPath))
	}
	file.WithBlock(`
	GRPCServer: &{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier/transport_grpc" "" }}.ServeGRPC{
		Port: 8900,
	},
	HTTPServer: &{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier/transport_http" "" }}.ServeHTTP{
		Port:     8800,
		WithCORS: true,
	},
}
`)

	return file
}

func (s *ServiceGenerator) createRouterV0RootFile() *files.GoFile {
	file := files.NewGoFile("v0")
	file.WithBlock(`
var Router = {{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier" "" }}.NewRouter(V0Router{})

type V0Router struct {
	{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier" "" }}.EmptyOperator
}

func (V0Router) Path() string {
	return "/v0"
}
`)

	return file
}

func (s *ServiceGenerator) createRouterRootFile(cwd string) *files.GoFile {
	pkgPath := path.Join(s.opt.PackageName, "internal/routers/v0")
	filePath := path.Join(cwd, "internal/routers/v0")

	file := files.NewGoFile("routers")
	file.WithBlock(fmt.Sprintf(`
var Router = {{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier" "" }}.NewRouter(RootRouter{})

func init() {
	Router.Register({{ .UseWithoutAlias "%s" "%s" }}.Router)
}

type RootRouter struct {
	{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/courier" "" }}.EmptyOperator
}

func (RootRouter) Path() string {
	return "/%s"
}
`, pkgPath, filePath, GetServiceName(s.opt.Name)))

	return file
}

func (s *ServiceGenerator) createMainFile(cwd string) *files.GoFile {
	globalPkgPath := path.Join(s.opt.PackageName, "internal/global")
	globalFilePath := path.Join(cwd, "internal/global")

	file := files.NewGoFile("main")
	file.WithBlock(`
func main() {
	app := application.NewApplication(runner`)

	if s.opt.ApolloSupport {
		file.WithBlock(fmt.Sprintf(`, {{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/application" "" }}.WithApollo(&{{ .UseWithoutAlias "%s" "%s" }}.ApolloConfig)`, globalPkgPath, globalFilePath))
	}

	file.WithBlock(fmt.Sprintf(`, application.WithConfig(&{{ .UseWithoutAlias "%s" "%s" }}.Config`, globalPkgPath, globalFilePath))

	if s.opt.DatabaseSupport {
		pkgPath := path.Join(s.opt.PackageName, "internal/databases")
		filePath := path.Join(cwd, "internal/databases")
		file.WithBlock(fmt.Sprintf(`, &{{ .UseWithoutAlias "%s" "%s" }}.Config`, pkgPath, filePath))
	}

	file.WithBlock(`))
	app.AddCommand(&{{ .UseWithoutAlias "github.com/spf13/cobra" "" }}.Command{
		Use: "migrate",
		Run: func(cmd *{{ .UseWithoutAlias "github.com/spf13/cobra" "" }}.Command, args []string) {
			migrate(args)
		},
	})

	app.Start()
}
`)

	routerPkgPath := path.Join(s.opt.PackageName, "internal/routers")
	routerFilePath := path.Join(cwd, "internal/routers")

	file.WithBlock(fmt.Sprintf(`
func runner(ctx *{{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/context" "" }}.WaitStopContext) error {
	{{ .UseWithoutAlias "github.com/sirupsen/logrus" "" }}.SetLevel({{ .UseWithoutAlias "%s" "%s" }}.Config.LogLevel)
	go {{ .UseWithoutAlias "%s" "%s" }}.Config.GRPCServer.Serve({{ .UseWithoutAlias "%s" "%s" }}.Router)
	return {{ .UseWithoutAlias "%s" "%s" }}.Config.HTTPServer.Serve({{ .UseWithoutAlias "%s" "%s" }}.Router)
}
`, globalPkgPath, globalFilePath, globalPkgPath, globalFilePath, routerPkgPath, routerFilePath, globalPkgPath, globalFilePath, routerPkgPath, routerFilePath))

	file.WithBlock(fmt.Sprintf(`
func migrate(args []string) {
	if err := {{ .UseWithoutAlias "github.com/profzone/eden-framework/pkg/sqlx/migration" "" }}.Migrate({{ .UseWithoutAlias "%s" "%s" }}.Config.MasterDB, nil); err != nil {
		panic(err)
	}
}
`, globalPkgPath, globalFilePath))

	return file
}
