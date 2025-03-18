:: filepath: /Users/matttownsend/Documents/GitHub/frantic-cat/clear.bat
@echo off

:: Directories to clear
set directories=logs dumps backups database defaults

:: Base path
set base_path=.\data

:: Loop through each directory
for %%d in (%directories%) do (
    set target_dir=%base_path%\%%d
    
    :: Check if the directory exists
    if exist %target_dir% (
        :: Find and delete all files except .keep
        for /r %target_dir% %%f in (*) do (
            if not "%%~nxf"==".keep" del "%%f"
        )
        echo Cleared %target_dir%
    ) else (
        echo Directory %target_dir% does not exist
    )
)