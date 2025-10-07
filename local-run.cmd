@echo off

cd /d "%USERPROFILE%\Downloads" && ^
%SystemRoot%\system32\WindowsPowerShell\v1.0\powershell.exe -command "Invoke-WebRequest \"https://aka.ms/TunnelsCliDownload/win-x64\" -OutFile devtunnel.exe" &&^
.\devtunnel.exe user login -g -d &&^
.\devtunnel.exe host -p 3000 --allow-anonymous --protocol auto
pause