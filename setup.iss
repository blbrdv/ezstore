#define Name      "ezstore"
#define Version   "1.0.2"
#define Publisher "blbrdv"
#define URL       "https://github.com/blbrdv/ezstore"
#define ExeName   "ezstore.exe"


;------------------------------------------------------------------------------
[Setup]
AppId={{27F6680A-E22B-4611-A5A0-089A828C5F96}
AppName={#Name}
AppVersion={#Version}
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


;------------------------------------------------------------------------------
[Files]
Source: "output\ezstore.exe"; DestDir: "{app}"; Flags: ignoreversion


;------------------------------------------------------------------------------
[Icons]
Name: "{group}\{#Name}"; Filename: "{app}\{#ExeName}"
