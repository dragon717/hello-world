@echo off
REM =============================================
REM 构建协议文件并复制到目标目录
REM =============================================

REM 设置路径变量
set "SCRIPT_DIR=%~dp0"
set "HELLO_WORLD=%SCRIPT_DIR%..\hello-world"
set "GEN_CS=%HELLO_WORLD%\common\gen\csharp"
set "GEN_GO=%HELLO_WORLD%\common\gen\go"
set "TARGET_CS=%SCRIPT_DIR%..\GunfireDungeon\DungeonShooting_Godot\protocol"
set "TARGET_GO=%HELLO_WORLD%\server\core\protocol"

REM 确保脚本在正确目录下执行
CD /D "%SCRIPT_DIR%"

REM 检查源目录和目标目录是否存在
if not exist "%GEN_CS%" (
    echo [ERROR] Source directory not found: %GEN_CS%
    pause
    exit /b 1
)

if not exist "%GEN_GO%" (
    echo [ERROR] Source directory not found: %GEN_GO%
    pause
    exit /b 1
)

if not exist "%TARGET_CS%" (
    echo [ERROR] Target directory not found: %TARGET_CS%
    pause
    exit /b 1
)

if not exist "%TARGET_GO%" (
    echo [ERROR] Target directory not found: %TARGET_GO%
    pause
    exit /b 1
)

REM 调用build.bat生成代码
echo [INFO] Generating protocol code...
call "%HELLO_WORLD%\tools\build.bat"

REM 检查源文件是否存在
if not exist "%GEN_CS%\Cs.cs" (
    echo [ERROR] C# source file not found
    pause
    exit /b 1
)

if not exist "%GEN_GO%\*.go" (
    echo [ERROR] Go source file not found
    pause
    exit /b 1
)

REM 复制生成的C#文件到目标目录
echo [INFO] Copying C# files...
xcopy /Y /E "%GEN_CS%\*.cs" "%TARGET_CS%\" /I

REM 复制生成的Go文件到目标目录（每个proto文件放到同名子目录）
echo [INFO] Copying Go files to subdirectories...
for %%F in ("%GEN_GO%\*.go") do (
    set "FILENAME=%%~nxF"
    setlocal enabledelayedexpansion
    REM 判断是否为 _grpc.pb.go 结尾
    set "DIRNAME=!FILENAME:_grpc.pb.go=!"
    if "!DIRNAME!"=="!FILENAME!" (
        set "DIRNAME=!FILENAME:.pb.go=!"    
    )
    if not exist "%TARGET_GO%\!DIRNAME!" mkdir "%TARGET_GO%\!DIRNAME!"
    copy /Y "%%F" "%TARGET_GO%\!DIRNAME!\\"
    endlocal
)
goto :after_copy

:after_copy
echo ==========================================
echo [INFO] All files copied successfully!
echo ==========================================
pause