//go:build mage
// +build mage

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var Default = RpcCaller

var (
	//packageRegex = regexp.MustCompile(`^\s*package\s+([a-zA-Z0-9_.]+);`)
	importRegex    = regexp.MustCompile(`^\s*import\s+"([a-zA-Z0-9_.]+)";`)
	goPackageRegex = regexp.MustCompile(`^\s*option\s+go_package\s*=\s*"([^"]+)";`)
	serviceRegex   = regexp.MustCompile(`^\s*service\s+([A-Za-z0-9_]+)\s*{`)
	rpcRegex       = regexp.MustCompile(`^\s*rpc\s+([A-Za-z0-9_]+)\s*\(\s*([A-Za-z0-9_]+\.?[A-Za-z0-9_]*)\s*\)\s+returns\s*\(\s*([A-Za-z0-9_]+\.?[A-Za-z0-9_]*)\s*\);`)
)

func RpcCaller() {
	if err := Generate(); err != nil {
		fmt.Println(err)
	}
}

// ServiceMethod holds information about a single RPC method
type ServiceMethod struct {
	Name           string
	RequestType    string
	ResponseType   string
	FullMethodName string
}

// Service holds information about a service and its methods
type Service struct {
	FilePath    string
	FileName    string
	ServiceName string
	GoPackage   string
	Methods     []ServiceMethod
	Imports     []string
}

// Generate is the Mage target to generate Go code from proto files
func Generate() error {
	fmt.Println("Generating rpc_caller...")

	projectRoot, err := os.Getwd()
	if err != nil {
		return errors.New("get root directory failed")
	}

	entries, err := os.ReadDir(projectRoot)
	if err != nil {
		return fmt.Errorf("read root directory failed: %v", err)
	}

	var services []Service

	for _, entry := range entries {
		if entry.IsDir() {
			serviceName := entry.Name()
			subDir := filepath.Join(projectRoot, serviceName)

			protoFiles, err := filepath.Glob(filepath.Join(subDir, "*.proto"))
			if err != nil {
				continue
			}

			if len(protoFiles) == 0 {
				continue
			}

			for _, protoFile := range protoFiles {
				service, err := parseProtoFile(protoFile)
				if err != nil {
					//fmt.Printf("warn: parse %s failed: %v\n", protoFile, err)
					continue
				}

				services = append(services, service)
			}
		}
	}
	// // for all protoFile
	//err = filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		log.Printf("warn: err in begin walk %s: %v\n", path, err)
	//		return nil
	//	}
	//
	//	if info.IsDir() {
	//		return nil
	//	}
	//
	//	if strings.HasSuffix(info.Name(), ".proto") {
	//		service, err := parseProtoFile(path)
	//		if err != nil {
	//			//log.Printf("warn: parse %s failed: %v\n", path, err)
	//			return nil
	//		}
	//
	//		services = append(services, service)
	//	}
	//
	//	return nil
	//})
	// gen Go file for each service

	for _, service := range services {
		if err := generateGoFile(service); err != nil {
			return fmt.Errorf("gen Go file failed: %v", err)
		}
	}

	return nil
}

func parseProtoFile(protoFilePath string) (Service, error) {
	file, err := os.Open(protoFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("open proto file failed: %v", err)
	}
	defer file.Close()

	filePath_ := filepath.Dir(protoFilePath)
	fileName := strings.TrimSuffix(filepath.Base(protoFilePath), filepath.Ext(protoFilePath))
	scanner := bufio.NewScanner(file)
	var (
		goPackage      string
		serviceName    string
		methods        []ServiceMethod
		imports        []string
		alreadyImports map[string]struct{}
		allImports     map[string]string
	)

	inService := false

	for scanner.Scan() {
		line := scanner.Text()

		//if matches := packageRegex.FindStringSubmatch(line); matches != nil {
		//	protoPackage = matches[1]
		//	continue
		//}

		if matches := importRegex.FindStringSubmatch(line); matches != nil {
			pkg := strings.TrimSuffix(filepath.Base(matches[1]), filepath.Ext(matches[1]))
			allImports[pkg] = matches[1]
			continue
		}

		if matches := goPackageRegex.FindStringSubmatch(line); matches != nil {
			goPackage = matches[1]

			continue
		}

		if matches := serviceRegex.FindStringSubmatch(line); matches != nil {
			serviceName = matches[1]
			inService = true
			continue
		}

		if inService {
			if matches := rpcRegex.FindStringSubmatch(line); matches != nil {
				methodName := matches[1]
				requestType := matches[2]
				responseType := matches[3]
				fullMethodName := fmt.Sprintf("%s_%s_FullMethodName", toCamelCase(serviceName), toCamelCase(methodName))
				method := ServiceMethod{
					Name:           toCamelCase(methodName) + "Caller",
					RequestType:    toCamelCase(requestType),
					ResponseType:   toCamelCase(responseType),
					FullMethodName: fullMethodName,
				}
				methods = append(methods, method)

				if strings.Contains(requestType, ".") {
					imp := strings.Split(requestType, ".")[0]
					if f, ok := allImports[imp]; ok {
						if _, ok = alreadyImports[imp]; !ok {
							alreadyImports[imp] = struct{}{}
							gopkg, err := getGoPackage(filepath.Join(protoFilePath, f))
							if err != nil {
								fmt.Printf("get go package failed: %v", err)
							} else {
								imports = append(imports, gopkg)
							}
						}
					}
				}
				continue
			}

			if strings.Contains(line, "}") {
				inService = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Service{}, fmt.Errorf("scan proto file failed: %v", err)
	}

	if serviceName == "" {
		return Service{}, fmt.Errorf("no serviceName in proto")
	}

	if goPackage == "" {
		return Service{}, fmt.Errorf("to goPackage in proto")
	}

	sp := strings.Split(goPackage, "/")
	service := Service{
		FilePath:    filePath_,
		FileName:    fileName,
		ServiceName: serviceName,
		GoPackage:   sp[len(sp)-1],
		Methods:     methods,
	}

	return service, nil
}

func generateGoFile(service Service) error {
	if len(service.Methods) == 0 {
		return nil
	}

	tmpl := `package {{.GoPackage}}

import (
	{{- range .Imports }}
	"{{ .Imports }}"
	{{- end }}
	"github.com/openimsdk/protocol/rpccall"
	"google.golang.org/grpc"
)

func Init{{.ServiceNameCamel}}(conn *grpc.ClientConn) {
	{{- range .Methods }}
	{{ .Name }}.SetConn(conn)
	{{- end }}
}

var (
	{{- range .Methods }}
	{{ .Name }} = rpccall.NewRpcCaller[{{ .RequestType }}, {{ .ResponseType }}]({{ .FullMethodName }})
	{{- end }}
)
`

	data := struct {
		ServiceNameCamel string
		GoPackage        string
		Methods          []ServiceMethod
		Imports          []string
	}{
		ServiceNameCamel: toCamelCase(service.ServiceName),
		GoPackage:        service.GoPackage,
		Methods:          service.Methods,
		Imports:          service.Imports,
	}

	t, err := template.New("goFile").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("parse template failed: %v", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template failed: %v", err)
	}

	goFileName := strings.ToLower(service.FileName) + "_caller.go"
	goFilePath := filepath.Join(service.FilePath, goFileName)

	if err := os.WriteFile(goFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	formatFile(goFilePath)
	return nil
}

func getGoPackage(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open proto file failed: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := goPackageRegex.FindStringSubmatch(line); matches != nil {
			return matches[1], nil
		}
	}
	return "", errors.New("goPackage not found")
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func formatFile(filePath string) {
	_ = runGoImports(filePath)
	_ = runGoFmt(filePath)
}

func runGoImports(filePath string) error {
	cmd := exec.Command("goimports", "-w", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goimports error: %v, stderr: %s", err, stderr.String())
	}
	return nil
}

func runGoFmt(filePath string) error {
	cmd := exec.Command("gofmt", "-w", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gofmt error: %v, stderr: %s", err, stderr.String())
	}
	return nil
}
