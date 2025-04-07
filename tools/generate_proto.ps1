# Set the path to protoc compiler
$PROTOC = "$PSScriptRoot\protoc-30.2-win64\bin\protoc.exe"

# Set the root directory of the project (one level up from tools)
$PROJECT_ROOT = (Get-Item $PSScriptRoot).Parent.FullName

# Function to generate proto files
function Generate-Proto {
    param (
        [string]$ProtoFile
    )
    
    Write-Host "Generating code for: $ProtoFile"
    
    # Run protoc with go and go-grpc plugins
    & $PROTOC `
        --proto_path="$PROJECT_ROOT" `
        --proto_path="$PSScriptRoot\protoc-30.2-win64\include" `
        --go_out="$PROJECT_ROOT" `
        --go_opt=module=github.com/hello-world `
        --go-grpc_out="$PROJECT_ROOT" `
        --go-grpc_opt=module=github.com/hello-world `
        "$ProtoFile"
        
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to generate code for $ProtoFile"
        exit 1
    }
}

# Find all .proto files in the project's proto directory
Write-Host "Searching for .proto files..."
$protoFiles = Get-ChildItem -Path "$PROJECT_ROOT\proto" -Filter "*.proto" -Recurse

if ($protoFiles.Count -eq 0) {
    Write-Warning "No .proto files found!"
    exit 0
}

# Generate code for each proto file
foreach ($file in $protoFiles) {
    Generate-Proto -ProtoFile $file.FullName
}

Write-Host "Proto generation completed successfully!" -ForegroundColor Green 