@echo off
REM 确保脚本在tools目录下执行
cd /d %~dp0

REM 1. 配置Protoc路径（如果未加入系统PATH）
set PROTOC=protoc-30.2-win64\bin\protoc.exe

REM 2. 定义输入/输出路径（相对当前目录）
set PROTO_DIR=.
set GO_OUT_DIR=../common/gen/go
set GO_MODULE=github.com/muniao/hello-world

REM 3. 批量编译当前目录下所有.proto文件
for %%f in (%PROTO_DIR%\*.proto) do (
    echo 正在编译: %%~nxf
    %PROTOC% ^
        --proto_path=%PROTO_DIR% ^
        --go_out=%GO_OUT_DIR% ^
        --go_opt=module=%GO_MODULE% ^
        %%f
    if errorlevel 1 (
        echo [错误] 编译失败: %%~nxf
        pause
        exit /b 1
    )
)

REM 4. 完成提示
echo -----------------------------------
echo 所有.proto文件编译完成！
pause