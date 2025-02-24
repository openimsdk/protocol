//go:build mage
// +build mage

package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Default = InstallDepend

var Aliases = map[string]any{
	"go":     GenGo,
	"java":   GenJava,
	"kotlin": GenKotlin,
	"csharp": GenCSharp,
	"js":     GenJavaScript,
	"ts":     GenTypeScript,
	"swift":  GenSwift,

	"dep": InstallDepend,

	"m:go":     Meeting.GenGo,
	"m:java":   Meeting.GenJava,
	"m:kotlin": Meeting.GenKotlin,
	"m:csharp": Meeting.GenCSharp,
	"m:js":     Meeting.GenJavaScript,
	"m:ts":     Meeting.GenTypeScript,
	"m:swift":  Meeting.GenSwift,
}

// Langeuage target
// Define output directories for each target language
const (
	GO     = "go"
	JAVA   = "java"
	CSharp = "csharp"
	Kotlin = "kotlin"
	JS     = "js"
	TS     = "ts"
	RS     = "rust"
	SWIFT  = "swift"
)

var protoModules = []string{
	"auth",
	"conversation",
	"errinfo",
	"group",
	"jssdk",
	"msg",
	"msggateway",
	"push",
	"relation",
	"rtc",
	"sdkws",
	"third",
	"user",
	"wrapperspb",
}

// install proto plugin
func InstallDepend() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)

	// log.Println("installing protoc-gen-go and protoc-gen-go-grpc")
	log.Println("installing protobuf dependencies in Go.")

	cmds := [][]string{
		{"install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"},
		{"install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command("go", cmdArgs...)

		// log.Println("running command:", "go", cmdArgs)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("command %v error: %v", cmdArgs, err)
			return err
		}
	}

	return nil
}

func GenDocs() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.Println("Generating documentation from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	docsOutDir := filepath.Join(".", "docs")

	for _, module := range protoModules {
		if err := os.MkdirAll(filepath.Join(docsOutDir, module), 0755); err != nil {
			return err
		}

		args := []string{
			// "--doc_out=" + filepath.Join(docsOutDir, module),
			"--doc_out=" + filepath.Join(docsOutDir),
			"--doc_opt=markdown," + strings.Join([]string{module, "md"}, "."),
			filepath.Join(module, module) + ".proto",
		}
		// log.Println(protoc, args)

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating documentation for module %s: %v\n", module, err)
			continue
		}
	}

	return nil
}

// Generate code for all languages (Go, Java, C#, JS, TS) from protobuf files.
func AllProtobuf() error {
	if err := GenGo(); err != nil {
		return err
	}
	if err := GenJava(); err != nil {
		return err
	}
	if err := GenCSharp(); err != nil {
		return err
	}
	if err := GenJavaScript(); err != nil {
		return err
	}
	if err := GenTypeScript(); err != nil {
		return err
	}
	return nil
}

// Generate Go code from protobuf files.
func GenGo() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Go code from proto files")

	// goOutDir := filepath.Join(protoDir, GO)

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		args := []string{
			// "--proto_path=" + filepath.Join(".", module),
			"--go_out=" + filepath.Join(".", module),
			"--go-grpc_out=" + filepath.Join(".", module),
			"--go_opt=module=github.com/openimsdk/protocol/" + strings.Join([]string{module}, "/"),
			"--go-grpc_opt=module=github.com/openimsdk/protocol/" + strings.Join([]string{module}, "/"),
			filepath.Join(module, module) + ".proto",
		}
		// log.Println("protoc args", args)

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Go code for module %s: %v\n", module, err)
			continue
		}
	}

	if err := removeOmitemptyTags(); err != nil {
		log.Println("Remove Omitempty is Error", err)
		return err
	} else {
		log.Println("Remove Omitempty is Success")
	}

	return nil
}

// Generate Java code from protobuf files.
func GenJava() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Java code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		javaOutDir := filepath.Join(".", module, JAVA)

		if err := os.MkdirAll(javaOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--java_out=lite:" + javaOutDir,
			filepath.Join(module, module) + ".proto",
		}
		log.Println(javaOutDir)

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Java code for module %s: %v\n", module, err)
			continue
		}
	}

	return nil
}

// Generate Kotlin code from protobuf files.
func GenKotlin() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Kotlin code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		kotlinOutDir := filepath.Join(".", module, Kotlin)

		if err := os.MkdirAll(kotlinOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--kotlin_out=" + kotlinOutDir,
			filepath.Join(module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Kotlin code for module %s: %v\n", module, err)
			continue
		}
	}

	return nil
}

// Generate C# code from protobuf files.
func GenCSharp() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating C# code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range protoModules {
		csharpOutDir := filepath.Join(".", module, CSharp)

		if err := os.MkdirAll(csharpOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--csharp_out=" + csharpOutDir,
			filepath.Join(module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating C# code for module %s: %v\n", module, err)
			continue
		}
	}

	return nil
}

func GenJavaScript() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating JavaScript code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	jsDir := filepath.Join(".", "pb", "js")
	args := []string{
		"--js_out=import_style=commonjs,binary:" + jsDir,
	}

	if err := os.MkdirAll(jsDir, 0755); err != nil {
		return err
	}

	for _, module := range protoModules {
		jsOutDir := filepath.Join(".", module, JS)

		if err := os.MkdirAll(jsOutDir, 0755); err != nil {
			return err
		}

		args = append(args,
			filepath.Join(module, module)+".proto")
	}

	cmd := exec.Command(protoc, args...)
	connectStd(cmd)

	if err := cmd.Run(); err != nil {
		// log.Printf("Error generating JS code for module %s: %v\n", module, err)
		log.Panicf("Error generating JS code: %v\n", err)
		// continue
	}

	return nil
}

// Need to install `ts-proto`.
// Generate TypeScript code from protobuf files.
func GenTypeScript() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.Println("Generating TypeScript code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	tsProto := filepath.Join(".", "node_modules", ".bin", "protoc-gen-ts_proto")

	if runtime.GOOS == "windows" {
		tsProto = filepath.Join(".", "node_modules", ".bin", "protoc-gen-ts_proto.cmd")
	}

	if _, err := os.Stat(tsProto); err != nil {
		log.Println("tsProto Not Found. Error: ", err, " tsProto Path: ", tsProto)
		return err
	}

	for _, module := range protoModules {
		// tsOutDir := filepath.Join(".", module, TS)
		tsOutDir := filepath.Join("pb", TS)

		if err := os.MkdirAll(tsOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--plugin=protoc-gen-ts_proto=" + tsProto,
			"--ts_proto_opt=messages=true,outputJsonMethods=false,outputPartialMethods=false,outputClientImpl=false,outputEncodeMethods=false,useOptionals=messages",
			"--ts_proto_out=" + tsOutDir,
			filepath.Join(module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating TypeScript code for module %s: %v\n", module, err)
			continue
		}
	}

	return nil
}

// Generate Swift code from protobuf files.
func GenSwift() error {
	// Configure logging
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Swift code from proto files")

	// Find protoc and Swift plugin paths
	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	// Iterate over proto modules to generate Swift code
	for _, module := range protoModules {
		swiftOutDir := filepath.Join(".", module, SWIFT)

		modulePath := filepath.Join(module, module+".proto")

		// Ensure the output directory for the module exists
		if err := os.MkdirAll(swiftOutDir, 0755); err != nil {
			return err
		}

		// Prepare protoc command
		args := []string{
			"--swift_out=" + swiftOutDir,
			"--swift_opt=Visibility=" + "Public",
			modulePath,
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd) // Connect command's output to standard output for logging

		// Run the command and handle errors
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Swift code for module %s: %v\n", module, err)
			continue
		}

		log.Printf("Successfully generated Swift code for module %s\n", module)
	}

	return nil
}

// Generate Harmony JavaScript code from protobuf files.
// Note: please install pbjs and pbts command first
// Reference Link: https://ohpm.openharmony.cn/#/cn/detail/@ohos%2Fprotobufjs
func GenHarmonyTS() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	log.Println("Generating Harmony TypeScript code from proto files")

	// Generate js

	outJSFile := "proto.js"
	args := []string{
		"-t", "static-module",
		"-w", "es6",
		"-o", outJSFile}

	for _, module := range protoModules {
		protoFile := filepath.Join(module, module) + ".proto"
		args = append(args, protoFile)
	}

	jscmd := exec.Command("pbjs", args...)
	jscmd.Env = os.Environ()
	connectStd(jscmd)

	log.Println("Running harmony js command", jscmd.String())
	if err := jscmd.Run(); err != nil {
		log.Printf("Error generating Harmony JS code: %v\n", err)
	}

	// Generate ts definition
	outTSDefFile := "proto.d.ts"
	tscmd := exec.Command("pbts",
		outJSFile,
		"-o", outTSDefFile,
	)

	tscmd.Env = os.Environ()
	connectStd(tscmd)

	log.Println("Running harmony ts command", tscmd.String())
	if err := tscmd.Run(); err != nil {
		log.Printf("Error generating Harmony TS code: %v\n", err)
	}

	// Modify the generated files
	// 1

	replaceStr := func(filePath, oldStr, newStr string) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Panic("failed to read file: %w", err)
		}

		originalContent := string(content)
		modifiedContent := strings.Replace(originalContent, oldStr, newStr, 1) // 只替换一次

		if originalContent == modifiedContent {
			return
		}
		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			log.Panic("failed to write file: %w", err)
		}
	}
	replaceStr(outJSFile, "import * as $protobuf from \"protobufjs/minimal\";", "import { index } from \"@ohos/protobufjs\"; \nconst $protobuf = index; \n import Long from 'long';\n$protobuf.util.Long=Long \n$protobuf.configure()")
	replaceStr(outTSDefFile, "import * as $protobuf from \"protobufjs\";\nimport Long = require(\"long\");", "import * as $protobuf from \"@ohos/protobufjs\"\nimport Long from 'long';")

	return nil
}

// ------------------
func connectStd(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func getWorkDirToolPath(name string) string {
	toolPath := ""
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("Error", err)
		return toolPath
	}
	toolsPath := filepath.Join(workDir, "tools")
	filepath.Walk(toolsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) == name {
			toolPath = path
		}
		return nil
	})

	return toolPath
}

func getToolPath(name string) (string, error) {
	// Get in work dir.
	toolPath := getWorkDirToolPath(name)
	if toolPath != "" {
		return toolPath, nil
	}

	// Get in env path.
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	// check under gopath
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	p := filepath.Join(gopath, "bin", name)

	if _, err := os.Stat(p); err != nil {
		return "", err
	}

	return p, nil
}

func removeOmitemptyTags() error {
	// protoGoDir := filepath.Join(protoDir, GO) // "./proto/go"

	re := regexp.MustCompile(`,\s*omitempty`)

	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("access path error:", err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".pb.go") {
			input, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("ReadFile error. Path: %s, Error %v", path, err)
				return err
			}

			output := re.ReplaceAllString(string(input), "")

			// check replace is happened
			if string(input) != output {
				err = os.WriteFile(path, []byte(output), info.Mode())
				if err != nil {
					fmt.Printf("Error writing file: %s, error: %v\n", path, err)
					return err
				}
				// fmt.Println("Modified file:", path)
			}
		}

		return nil
	})
}
