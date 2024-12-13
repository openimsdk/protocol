//go:build mage
// +build mage

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
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
}

// Generate is the Mage target to generate Go code from proto files
func Generate() error {
	fmt.Println("Generating service...")

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
		return Service{}, fmt.Errorf("打开 proto 文件失败: %v", err)
	}
	defer file.Close()

	filePath := filepath.Dir(protoFilePath)
	fileName := strings.TrimSuffix(filepath.Base(protoFilePath), filepath.Ext(protoFilePath))
	scanner := bufio.NewScanner(file)
	var (
		goPackage   string
		serviceName string
		methods     []ServiceMethod
	)

	//packageRegex := regexp.MustCompile(`^\s*package\s+([a-zA-Z0-9_.]+);`)
	goPackageRegex := regexp.MustCompile(`^\s*option\s+go_package\s*=\s*"([^"]+)";`)
	serviceRegex := regexp.MustCompile(`^\s*service\s+([A-Za-z0-9_]+)\s*{`)
	rpcRegex := regexp.MustCompile(`^\s*rpc\s+([A-Za-z0-9_]+)\s*\(\s*([A-Za-z0-9_]+)\s*\)\s+returns\s*\(\s*([A-Za-z0-9_]+)\s*\);`)

	inService := false

	for scanner.Scan() {
		line := scanner.Text()

		//if matches := packageRegex.FindStringSubmatch(line); matches != nil {
		//	protoPackage = matches[1]
		//	continue
		//}

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
		FilePath:    filePath,
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
	}{
		ServiceNameCamel: toCamelCase(service.ServiceName),
		GoPackage:        service.GoPackage,
		Methods:          service.Methods,
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

	// 写入文件
	if err := os.WriteFile(goFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	log.Printf("gen %s success", goFilePath)
	return nil
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
