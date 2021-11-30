cd /D "%~dp0"
SET  mycwd=%CD%\manifest.json
echo %mycwd%
REG ADD "HKCU\Software\Google\Chrome\NativeMessagingHosts\mynotesaver" /ve /t REG_SZ /d %mycwd% /f