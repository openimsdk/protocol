//go:build mage
// +build mage

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
)

type Meeting mg.Namespace

// var meetingModules = []string{
// 	"admin",
// 	"meeting",
// 	"user",
// }

var meetingPath = filepath.Join(".", "openmeeting")

var meetingModules = []string{
	"admin",
	"meeting",
	"user",
}

func (Meeting) GenGo() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Go code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")

	for _, module := range meetingModules {
		meetingGoOutPath := filepath.Join(meetingPath, module, GO)

		if err := os.MkdirAll(meetingGoOutPath, 0755); err != nil {
			return err
		}

		args := []string{
			"--go_out=" + meetingGoOutPath,
			"--go-grpc_out=" + meetingGoOutPath,
			"--go_opt=module=github.com/openimsdk/protocol/openmeeting/" + strings.Join([]string{module}, "/"),
			"--go-grpc_opt=module=github.com/openimsdk/protocol/openmeeting/" + strings.Join([]string{module}, "/"),
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Go code for meeting module %s: %v\n", module, err)
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

func (Meeting) GenJava() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Java code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	var meetingPath = filepath.Join(".", "openmeeting")

	for _, module := range meetingModules {
		meetingJavaOutPath := filepath.Join(meetingPath, module, JAVA)

		if err := os.MkdirAll(meetingJavaOutPath, 0755); err != nil {
			return err
		}

		args := []string{
			"--java_out=lite:" + meetingJavaOutPath,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Java code for meeting module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

func (Meeting) GenKotlin() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating Kotlin code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range meetingModules {
		meetingKotlinOutPath := filepath.Join(meetingPath, module, Kotlin)

		if err := os.MkdirAll(meetingKotlinOutPath, 0755); err != nil {
			return err
		}

		args := []string{
			"--kotlin_out=lite:" + meetingKotlinOutPath,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Error generating Kotlin code for meeting module %s: %v\n", module, err)
			continue
		}
	}
	return nil
}

// Generate C# code from protobuf files.
func (Meeting) GenCSharp() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating C# code from meeting proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range meetingModules {
		meetingCsharpOutDir := filepath.Join(meetingPath, module, CSharp)

		if err := os.MkdirAll(meetingCsharpOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--csharp_out=" + meetingCsharpOutDir,
			filepath.Join(meetingPath, module, module) + ".proto",
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

func (Meeting) GenJavaScript() error {
	log.SetOutput(os.Stdout)
	// log.SetFlags(log.Lshortfile)
	log.Println("Generating JavaScript code from proto files")

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	for _, module := range meetingModules {
		meetingJsOutDir := filepath.Join(meetingPath, module, JS)

		if err := os.MkdirAll(meetingJsOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--js_out=import_style=commonjs,binary:" + meetingJsOutDir,
			filepath.Join(meetingPath, module, module) + ".proto",
		}

		cmd := exec.Command(protoc, args...)
		connectStd(cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("Error generating JS code for module %s: %v\n", module, err)
			continue
		}
	}

	return nil
}

// Generate TypeScript code from protobuf files.
func (Meeting) GenTypeScript() error {
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

	for _, module := range meetingModules {
		// meetingTsOutDir := filepath.Join(meetingPath, module, TS)
		meetingTsOutDir := filepath.Join(meetingPath, "pb", TS)

		if err := os.MkdirAll(meetingTsOutDir, 0755); err != nil {
			return err
		}

		args := []string{
			"--plugin=protoc-gen-ts_proto=" + tsProto,
			// "--proto_path=" + module,
			"--ts_proto_opt=messages=true,outputJsonMethods=false,outputPartialMethods=false,outputClientImpl=false,outputEncodeMethods=false,useOptionals=messages",
			"--ts_proto_out=" + meetingTsOutDir,
			filepath.Join(meetingPath, module, module) + ".proto",
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
func (Meeting) GenSwift() error {
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
	for _, module := range meetingModules {
		swiftOutDir := filepath.Join(".", module, SWIFT)

		modulePath := filepath.Join(module, module+".proto")
		// outputPath := filepath.Join(swiftOutDir, module)

		// Ensure the output directory for the module exists
		if err := os.MkdirAll(swiftOutDir, 0755); err != nil {
			return err
		}

		// Prepare protoc command
		args := []string{
			// "--proto_path=" + protoDir,
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

// ------------------
// func connectStd(cmd *exec.Cmd) {
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// }

// func getWorkDirToolPath(name string) string {
// 	toolPath := ""
// 	workDir, err := os.Getwd()
// 	if err != nil {
// 		log.Println("Error", err)
// 		return toolPath
// 	}
// 	toolsPath := filepath.Join(workDir, "tools")
// 	filepath.Walk(toolsPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) == name {
// 			toolPath = path
// 		}
// 		return nil
// 	})

// 	return toolPath
// }

// func getToolPath(name string) (string, error) {
// 	// Get in work dir.
// 	toolPath := getWorkDirToolPath(name)
// 	if toolPath != "" {
// 		return toolPath, nil
// 	}

// 	// Get in env path.
// 	if p, err := exec.LookPath(name); err == nil {
// 		return p, nil
// 	}

// 	// check under gopath
// 	gopath := os.Getenv("GOPATH")
// 	if gopath == "" {
// 		gopath = build.Default.GOPATH
// 	}
// 	p := filepath.Join(gopath, "bin", name)

// 	if _, err := os.Stat(p); err != nil {
// 		return "", err
// 	}

// 	return p, nil
// }

// func removeOmitemptyTags() error {
// 	// protoGoDir := filepath.Join(protoDir, GO) // "./proto/go"

// 	re := regexp.MustCompile(`,\s*omitempty`)

// 	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			fmt.Println("access path error:", err)
// 			return err
// 		}
// 		if !info.IsDir() && strings.HasSuffix(path, ".pb.go") {
// 			input, err := os.ReadFile(path)
// 			if err != nil {
// 				fmt.Println("ReadFile error. Path: %s, Error %v", path, err)
// 				return err
// 			}

// 			output := re.ReplaceAllString(string(input), "")

// 			// check replace is happened
// 			if string(input) != output {
// 				err = os.WriteFile(path, []byte(output), info.Mode())
// 				if err != nil {
// 					fmt.Printf("Error writing file: %s, error: %v\n", path, err)
// 					return err
// 				}
// 				// fmt.Println("Modified file:", path)
// 			}
// 		}

// 		return nil
// 	})
// }
