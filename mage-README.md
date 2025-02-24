# Generate Protocol Buffers with Mage

Only generate Go and TS pb file in this repo.

## Preconditions

- Go 1.18 or higher
- Install Protocol Buffer Compiler (protoc) in your PATH.
- Install mage bin in your PATH.

> Notice:  
> Usually, different versions of `protoc` do not have a significant impact, as they are compatible with each other.

## Install Go

- View the [Go install documentation](https://go.dev/doc/install) to install Go.
- Add `go/bin` to the PATH environment variable.

## Install Protocol Buffer Compiler (protoc)

view the [Proto Buffer Compiler Installation docs](https://grpc.io/docs/protoc-installation/).

- Download the latest release of the Protocol Buffers compiler zip file corresponding to your OS and Arch.
- Unzip the file under `$HOME/.local` or directory of your choice.
- Update the PATH environment variable to include the `protoc` executable.

## Install mage

Using Go Install with Go version >= 1.18:

```shell
go install github.com/magefile/mage@latest
```

<details>
<summary>With Go version < 1.18:</summary>

You can use `bootstrap_install_mage.bat` or `bootstrap_install_mage.sh` to install mage fast.

</details>

## Compiling Your Protocol Buffers

### Go Generated Code

- Execute `mage InstallDepend` to install Go dependencies.

- Execute `mage GenGo` to generate Go code.

- You can also view the [Go Usage Docs](https://grpc.io/docs/languages/go/quickstart/#prerequisites) for more information.

### TypeScript Generated Code

- Execute `npm install ts-proto` to the workdir.
- Execute `mage GenTypeScript` to generate TypeScript code.

## Modify the Protocol Buffers

### Writing Protocol Buffers

In our example, we have a simple `hello/hello.proto` file:

```proto
syntax = "proto3";

// define a request message
message HelloRequest {
  string name = 1;
  UserInfo user = 2;
}

// define a response message
message HelloResponse {
  string message = 1;
}

// define a parameter message
message UserInfo {
  string name = 1;
  int32 age = 2;
}

// define a service
service HelloService {
  // define a rpc method
  rpc SayHello (HelloRequest) returns (HelloResponse);
}
```

- Write your method request and response messages. Like `HelloRequest` and `HelloResponse`.
- Write your service method. Like `SayHello`.
- You can also define the parameter message, like `UserInfo`.
- add the module name to `protoModules` variable in `magefile.go`. In this example, you need append `"hello"` in `protoModules`. We recommend using the module name as the directory name and file name.
- Execute corresponding languge command to generate protobuf code. More to view [Compiling Your Protocol Buffers](#compiling-your-protocol-buffers).

## More:

- [Protobuf docs](https://protobuf.dev/)
- [gRPC docs](https://grpc.io/docs)
