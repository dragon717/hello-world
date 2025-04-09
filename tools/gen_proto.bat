@echo off
setlocal

:: Set paths
set WORKSPACE_ROOT=%~dp0..
set PROTO_PATH=%WORKSPACE_ROOT%\common\proto
set GO_OUT_PATH=%WORKSPACE_ROOT%\ai-mcp-server\proto

:: Create output directory if not exists
if not exist "%GO_OUT_PATH%" mkdir "%GO_OUT_PATH%"

:: Install specific versions of protoc generators compatible with Go 1.21
echo Installing compatible protoc generators...
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

:: Clean existing generated files
echo Cleaning existing generated files...
if exist "%GO_OUT_PATH%\ai" rd /s /q "%GO_OUT_PATH%\ai"
mkdir "%GO_OUT_PATH%\ai"

:: Generate protobuf Go code
echo Generating protobuf Go code...
protoc --proto_path=%PROTO_PATH% ^
       --go_out=%GO_OUT_PATH% ^
       --go_opt=paths=source_relative ^
       --go-grpc_out=%GO_OUT_PATH% ^
       --go-grpc_opt=paths=source_relative ^
       %PROTO_PATH%\ai_service.proto

if %ERRORLEVEL% NEQ 0 (
    echo Error generating protobuf code
    exit /b 1
)

echo Proto generation completed successfully
exit /b 0 