@echo off
REM 确保脚本在tools目录下执行
cd /d %~dp0

REM 1. 配置Protoc和gRPC插件路径
set PROTOC=protoc-30.2-win64\bin\protoc.exe
set GRPC_PLUGIN=grpc_csharp_plugin.exe

REM 2. 定义输入/输出路径
set PROTO_DIR=.
set CS_OUT_DIR=../common/gen/csharp

REM 3. 检查gRPC插件是否存在
if not exist "%GRPC_PLUGIN%" (
    echo [错误] 找不到gRPC插件: %GRPC_PLUGIN%
    echo 请确保已安装gRPC C#插件
    pause
    exit /b 1
)

REM 4. 批量编译.proto文件
for %%f in (%PROTO_DIR%\*.proto) do (
    echo 正在编译: %%~nxf
    
    REM 生成普通proto消息类
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --csharp_out=%CS_OUT_DIR% ^
        %%f
    
    if errorlevel 1 (
        echo [错误] 编译失败: %%~nxf
        pause
        exit /b 1
    )
    
    REM 生成gRPC服务类（仅当文件包含service定义）
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --csharp_out=%CS_OUT_DIR% ^
        --grpc_out=%CS_OUT_DIR% ^
        --plugin=protoc-gen-grpc="%GRPC_PLUGIN%" ^
        %%f
    
    if errorlevel 1 (
        echo [警告] 生成gRPC服务失败: %%~nxf（可能没有service定义或插件问题）
    )
)

echo -----------------------------------
echo 所有.proto文件编译完成！
pause