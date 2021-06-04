package admin

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/qor/assetfs"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"github.com/qor/roles"
)

var (
	globalViewPaths []string
	globalAssetFSes []assetfs.Interface
	gomod           []string
)

// HasPermissioner has permission interface
type HasPermissioner interface {
	HasPermission(roles.PermissionMode, *qor.Context) bool
}

// ResourceNamer is an interface for models that defined method `ResourceName`
type ResourceNamer interface {
	ResourceName() string
}

// I18n define admin's i18n interface
type I18n interface {
	Scope(scope string) I18n
	Default(value string) I18n
	T(locale string, key string, args ...interface{}) template.HTML
}

// RegisterViewPath register view path for all assetfs
func RegisterViewPath(pth string) {
	globalViewPaths = append(globalViewPaths, pth)

	for _, assetFS := range globalAssetFSes {
		if assetFS.RegisterPath(filepath.Join(utils.AppRoot, "vendor", pth)) != nil {
			for _, gopath := range utils.GOPATH() {
				if assetFS.RegisterPath(filepath.Join(gopath, "src", pth)) == nil {
					break
				}
				pth = strings.TrimSuffix(pth, "/views")
				if assetFS.RegisterPath(filepath.Join(gopath, "pkg/mod", pth, "@", getDepVersionFromMod(pth), "views")) == nil {
					break
				}
			}
		}
	}
}

func equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func getDepVersionFromMod(pth string) (v string) {
	if len(gomod) == 0 {
		if cont, err := ioutil.ReadFile("go.mod"); err == nil {
			gomod = strings.Split(string(cont), "\n")
		}
	}

	for _, val := range gomod {
		if txt := strings.Trim(val, "\t\r"); strings.HasPrefix(txt, pth) {
			v = strings.Split(txt, " ")[1]
			return
		}
	}
	return
}
