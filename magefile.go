//go:build mage
// +build mage

package main

import (
	"bufio"
	"bytes"
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
	Name           string
	GoPackage      string
	QuoteGoPackage string
	ProtoPackage   string
	Methods        []ServiceMethod
	ProtoFilePath  string
}

// Generate is the Mage target to generate Go code from proto files
func Generate() error {
	fmt.Println("Generating service...")
	// 获取项目根目录
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	//// rpccall 目录
	//rpccallDir := filepath.Join(projectRoot, "rpccall")
	//if _, err := os.Stat(rpccallDir); os.IsNotExist(err) {
	//	if err := os.Mkdir(rpccallDir, 0755); err != nil {
	//		return fmt.Errorf("创建 rpccall 目录失败: %v", err)
	//	}
	//}
	fmt.Println("ReadDir...")
	// 遍历项目根目录下的所有子目录，查找 .proto 文件
	entries, err := os.ReadDir(projectRoot)
	if err != nil {
		return fmt.Errorf("读取项目根目录失败: %v", err)
	}
	var services []Service
	for _, entry := range entries {
		if entry.IsDir() {
			serviceName := entry.Name()
			protoFile := filepath.Join(projectRoot, serviceName, serviceName+".proto")
			if _, err := os.Stat(protoFile); os.IsNotExist(err) {
				log.Printf("警告: 找不到 %s 目录下的 %s.proto 文件，跳过...", serviceName, serviceName)
				continue
			}
			fmt.Println("parseProtoFile...")
			service, err := parseProtoFile(protoFile)
			if err != nil {
				//return fmt.Errorf("解析 %s 失败: %v", protoFile, err)
				continue
			}
			service.Name = serviceName
			service.ProtoFilePath = protoFile
			services = append(services, service)
		}
	}
	fmt.Println("generateGoFile...")
	// 为每个服务生成 Go 文件
	for _, service := range services {
		if err := generateGoFile(filepath.Join(projectRoot, service.Name), service); err != nil {
			return fmt.Errorf("生成 Go 文件失败: %v", err)
		}
	}

	return nil
}

// parseProtoFile 解析单个 proto 文件，提取服务和方法信息
func parseProtoFile(protoFilePath string) (Service, error) {
	file, err := os.Open(protoFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("打开 proto 文件失败: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var protoPackage string
	var goPackage string
	var serviceName string
	var methods []ServiceMethod

	// 正则表达式
	packageRegex := regexp.MustCompile(`^\s*package\s+([a-zA-Z0-9_.]+);`)
	goPackageRegex := regexp.MustCompile(`^\s*option\s+go_package\s*=\s*"([^"]+)";`)
	serviceRegex := regexp.MustCompile(`^\s*service\s+([A-Za-z0-9_]+)\s*{`)
	rpcRegex := regexp.MustCompile(`^\s*rpc\s+([A-Za-z0-9_]+)\s*\(\s*([A-Za-z0-9_]+)\s*\)\s+returns\s+\(\s*([A-Za-z0-9_]+)\s*\);`)

	inService := false

	for scanner.Scan() {
		line := scanner.Text()

		// 查找 package
		if matches := packageRegex.FindStringSubmatch(line); matches != nil {
			protoPackage = matches[1]
			continue
		}

		// 查找 go_package
		if matches := goPackageRegex.FindStringSubmatch(line); matches != nil {
			goPackage = matches[1]

			continue
		}

		// 查找 service
		if matches := serviceRegex.FindStringSubmatch(line); matches != nil {
			serviceName = matches[1]
			inService = true
			continue
		}

		// 查找 rpc 方法
		if inService {
			if matches := rpcRegex.FindStringSubmatch(line); matches != nil {
				methodName := matches[1]
				requestType := matches[2]
				responseType := matches[3]
				fullMethodName := fmt.Sprintf("%s_%s_FullMethodName", toCamelCase(serviceName), toCamelCase(methodName))
				method := ServiceMethod{
					Name:           toCamelCase(methodName),
					RequestType:    toCamelCase(requestType),
					ResponseType:   toCamelCase(responseType),
					FullMethodName: fullMethodName,
				}
				methods = append(methods, method)
				continue
			}

			// 检查服务块是否结束
			if strings.Contains(line, "}") {
				inService = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Service{}, fmt.Errorf("扫描 proto 文件失败: %v", err)
	}

	if serviceName == "" {
		return Service{}, fmt.Errorf("在 proto 文件中未找到 service 定义")
	}

	if goPackage == "" {
		return Service{}, fmt.Errorf("在 proto 文件中未找到 go_package 选项")
	}

	sp := strings.Split(goPackage, "/")
	service := Service{
		ProtoPackage:   protoPackage,
		GoPackage:      goPackage,
		QuoteGoPackage: sp[len(sp)-1],
		Name:           serviceName,
		Methods:        methods,
	}

	return service, nil
}

// generateGoFile 使用模板生成 Go 文件
func generateGoFile(dir string, service Service) error {
	// 定义模板
	tmpl := `package rpccall

import (
	"{{.GoPackage}}"
	"google.golang.org/grpc"
	"github.com/openimsdk/protocol/rpccall"
)

func Init{{.ServiceNameCamel}}(conn *grpc.ClientConn) {
	{{- range .Methods }}
	{{ .Name }}.SetConn(conn)
	{{- end }}
}

var (
	{{- range .Methods }}
	{{ .Name }} = rpccall.NewRpcCaller[{{ $.QuoteGoPackage }}.{{ .RequestType }}, {{ $.QuoteGoPackage }}.{{ .ResponseType }}]({{ $.QuoteGoPackage }}.{{ .FullMethodName }})
	{{- end }}
)
`

	// 准备模板数据
	data := struct {
		ServiceNameCamel string
		GoPackage        string
		QuoteGoPackage   string
		Methods          []ServiceMethod
	}{
		ServiceNameCamel: toCamelCase(service.Name),
		GoPackage:        service.GoPackage,
		QuoteGoPackage:   service.QuoteGoPackage,
		Methods:          service.Methods,
	}

	// 解析模板
	t, err := template.New("goFile").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 生成内容
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	// 定义生成的文件路径
	goFileName := strings.ToLower(service.Name) + "_caller.go"
	goFilePath := filepath.Join(dir, goFileName)

	// 写入文件
	if err := os.WriteFile(goFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	log.Printf("生成 %s 成功", goFilePath)
	return nil
}

// toCamelCase 将字符串转换为大驼峰
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func upperFirst(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "_")
}
