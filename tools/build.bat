@echo off
REM =============================================
REM 批量编译.proto文件生成C#和Go代码
REM =============================================

REM 确保脚本在tools目录下执行
cd /d %~dp0

REM 1. 配置工具路径
set PROTOC=protoc-26.1-win64\bin\protoc.exe
set GRPC_PLUGIN=grpc_csharp_plugin.exe
set GO_PLUGIN=protoc-gen-go.exe
set GO_GRPC_PLUGIN=protoc-gen-go-grpc.exe

REM 2. 定义输入/输出路径
set PROTO_DIR=../common/out_protocol
set CS_OUT_DIR=../common/gen/csharp
set GO_OUT_DIR=../common/gen/go

REM 3. 创建输出目录（如果不存在）
if not exist "%CS_OUT_DIR%" mkdir "%CS_OUT_DIR%"
if not exist "%GO_OUT_DIR%" mkdir "%GO_OUT_DIR%"

REM 4. 检查插件是否存在
if not exist "%GRPC_PLUGIN%" (
    echo [错误] 找不到gRPC C#插件: %GRPC_PLUGIN%
    echo 请确保已安装gRPC C#插件
    pause
    exit /b 1
)

if not exist "%GO_PLUGIN%" (
    echo [错误] 找不到Go插件: %GO_PLUGIN%
    echo 请执行安装命令：go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    pause
    exit /b 1
)

if not exist "%GO_GRPC_PLUGIN%" (
    echo [错误] 找不到Go gRPC插件: %GO_GRPC_PLUGIN%
    echo 请执行安装命令：go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    pause
    exit /b 1
)

REM 5. 批量编译.proto文件
for %%f in (%PROTO_DIR%\*.proto) do (
    echo ==========================================
    echo 正在编译: %%~nxf
    
    REM 5.1 生成C#消息类
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --csharp_out=%CS_OUT_DIR% ^
        %%f
    
    if errorlevel 1 (
        echo [错误] C#消息类生成失败: %%~nxf
        pause
        exit /b 1
    )
    
    REM 5.2 生成C# gRPC服务类
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --grpc_out=%CS_OUT_DIR% ^
        --plugin=protoc-gen-grpc="%GRPC_PLUGIN%" ^
        %%f
    
    if errorlevel 1 (
        echo [警告] C# gRPC生成失败: %%~nxf（可能没有service定义）
    )
    
    REM 5.3 生成Go消息类
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --go_out=%GO_OUT_DIR% ^
        --go_opt=paths=source_relative ^
        %%f
    
    if errorlevel 1 (
        echo [错误] Go消息类生成失败: %%~nxf
        pause
        exit /b 1
    )
    
    REM 5.4 生成Go gRPC服务
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --go-grpc_out=%GO_OUT_DIR% ^
        --go-grpc_opt=paths=source_relative ^
        %%f
    
    if errorlevel 1 (
        echo [警告] Go gRPC生成失败: %%~nxf（可能没有service定义）
    )
    
    echo 完成编译: %%~nxf
)

REM 6. 完成提示
echo ==========================================
echo 所有.proto文件编译完成！
echo ==========================================