:: Settings
@echo off

:: Build
cd ..\cmd
go build -o ..\bin\chip-8.exe .
cd ..

:: Run
.\bin\chip-8.exe