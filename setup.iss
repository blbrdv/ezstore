#define Name        "ezstore"
#define Version     "<<<PRODUCT_VERSION>>>"
#define FileVersion "<<<FILE_VERSION>>>"
#define Publisher   "blbrdv"
#define URL         "https://github.com/blbrdv/ezstore"
#define ExeName     "ezstore.exe"


;------------------------------------------------------------------------------
[Setup]
AppId={{27F6680A-E22B-4611-A5A0-089A828C5F96}
AppName={#Name}
AppVersion={#Version}
VersionInfoVersion={#FileVersion}
AppPublisher={#Publisher}
AppPublisherURL={#URL}
AppSupportURL={#URL}
AppUpdatesURL={#URL}

DefaultDirName={commonpf}\{#Name}
DefaultGroupName={#Name}

OutputDir=.\output
OutputBaseFileName=ezsetup

LicenseFile=LICENSE
SetupIconFile=dist\icon.ico

Compression=lzma
SolidCompression=yes
ChangesEnvironment=true


;------------------------------------------------------------------------------
[Files]
Source: "output\ezstore.exe"; DestDir: "{app}/bin"; Flags: ignoreversion
Source: "dist\README.txt"; DestDir: "{app}"; Flags: ignoreversion


;------------------------------------------------------------------------------
[Icons]
Name: "{group}\{#Name}"; Filename: "{app}\{#ExeName}"


;------------------------------------------------------------------------------
[Tasks]
Name: envPath; Description: "Add to PATH variable"


;------------------------------------------------------------------------------
#include "innosetup/environment.iss"
[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall
  then EnvAddPath(ExpandConstant('{app}') + '\bin');
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall
  then EnvRemovePath(ExpandConstant('{app}') + '\bin');
end;
