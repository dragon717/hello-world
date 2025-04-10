@echo off
setlocal enabledelayedexpansion

:: 参数检查
if "%~1"=="" (
    echo Error: Commit message is required!
    echo Usage: %0 "commit message"
    exit /b 1
)

:: 配置路径
set SOURCE="C:\Users\huangrui\go\src\Test\*"
set TARGET="D:\hello-world\server\core"

:: 删除目标目录所有文件（慎用！确保路径正确）
echo 正在清空目标目录...
cd /d %TARGET%
if errorlevel 1 (
    echo 错误：无法进入目标目录
    exit /b 1
)
del /q * >nul 2>&1
for /d %%d in (*) do rmdir /s /q "%%d"

:: 复制文件（保留目录结构）
echo 正在复制文件...
xcopy %SOURCE% %TARGET% /s /e /h /y /q
if errorlevel 1 (
    echo 错误：文件复制失败
    exit /b 1
)

:: 执行Git操作
echo 正在执行Git操作...
cd /d %TARGET%
git add --all
if errorlevel 1 (
    echo 错误：git add失败
    exit /b 1
)

git commit -m "%~1"
if errorlevel 1 (
    echo 错误：git commit失败
    exit /b 1
)

git push
if errorlevel 1 (
    echo 错误：git push失败
    exit /b 1
)

echo 操作成功完成！
endlocal