# Sandboxer
<p align="center"><img src="resources/icon_transparent.png" width="200"/></p>

## Inspect objects using Trend Micro Vision One or Deep Discovery Analyzer sandbox

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

### Major Features
- Inspected objects: Files and URLs.
- Supported sandboxes: Vision One and Deep Discovery Analyzer. 
- Supported platforms: Windows and macOS.

###  Minor Features
- Auto notification and download when new version is released.
- Store analysis results for two months (this is configurable option).
- Show system notifications when malicious file is detected
- Show Vision One sandbox quota
- Support HTTP proxy server including basic and NTLM authentication

Sandboxer submissions window:

<p align="center"><img src="resources/captures/submissions.png" width="300px"/></p>

## Installation
Installation and usage video:

[![Installation and usage](https://img.youtube.com/vi/beRX6YXjC4k/0.jpg)](https://www.youtube.com/watch?v=beRX6YXjC4k)

Download [latest release](https://github.com/mpkondrashin/sandboxer/releases/latest) for your platform, unpack zip file, and run setup.exe (for Windows) or SandboxerInstaller (for macOS).

<p align="center"><kbd><img src="resources/captures/page_0.png" width="500px"/></kbd></p>
<p align="center">Choose your sanbox type: <a href="https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-sandbox-analysis">Vision One</a> or <a href="https://www.trendmicro.com/en_ca/business/products/network/advanced-threat-protection/analyzer.html">Deep Discovery Analyzer</a>. Then press "Next" button.</p>
<p align="center"><kbd><img src="resources/captures/page_5.png" width="500px"/></kbd></p>
<p align="center">If you selected Vision One on the first step, then enter Token. Learn more about <a href="https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-api-keys">API Keys</a> and <a href="https://docs.trendmicro.com/en-US/documentation/article/trend-vision-one-configuring-user-rol">Roles</a>. 
If correct Domain value is not detected automatically, choose it from dropdown list.</p> 
<p align="center"><kbd><img src="resources/captures/page_6.png" width="500px"/></kbd></p>
<p align="center">If you selected Deep Discovery Analyzer on the first step, provide its IP/DNS address and API Key. If you are using self-signed certificate, check TLS Errors Ignore.</p>
<p align="center"><kbd><img src="resources/captures/page_7.png" width="500px"/></kbd></p>
<p align="center">Remove checkbox is there is not need to run Sandboxer automatically. It will be launched automatically upon file submission.</p>
<p align="center"><kbd><img src="resources/captures/page_9.png" width="500px"/></kbd></p>
<p align="center">Wait for file copy process to finish.</p>
<p align="center"><kbd><img src="resources/captures/page_11.png" width="500px"/></kbd></p>
<p align="center">Press "Quit" button.</p>

## Usage

### To Submit File On macOS

Right click on file and choose Quick Actions -> Sandboxer
<p align="center"><img src="resources/captures/quick_actions.png" width="300px"/></p>

### To Submit Files On Windows

Right click on file and choose Send To -> Sandboxer. Note that for latest Windows, you will have to chosse first "Show more options".
<p align="center"><img src="resources/captures/send_to.png" width="300px"/></p>

### To Submit URL
Run Sandboxer, if it is not yet running, and pick from its system tray icon menu "Submit URL" item.
<p align="center"><img src="resources/captures/submit_url.png" width="300px"/></p>

### To Get Results
Pick from Sandboxer system tray icon menu "Submissions" item. Right click on the menu icon "⋮" and choose "Show Report" or "Investigation Package".

## Bugs

### Notifications
On Windows if notifications are disabled not by Sandboxer Options window, but using Notification Center, then it is not possible to turn them back on.

### New version update
When Sandboxer shows that new version is available, but upon downloading it shows error, it means that user has to wait several minutes and try again.

### Unregistration
If  unregistration is not performed from Options dialog and/or during uninstallation Deep Discovery Analyzer connection was not available, then Sandboxer will be kept on the submitters list and the only option will be to remove it using Web UI. After installing Sandboxer once more, it will just register one more submitter.

### Install error
If dufing files copy phase you encounter some error, try to return to previous install phase and try again.

### Install crash
To make sure that your installation completed, check sandboxer_setup_wizard.log log file genearted along the installer executable. Last line should contain the following "... INFO G0001 Close Logging ...". If not, try to run installer once more.



#### Boxer picture
Icon is taken from <a href="https://www.flaticon.com/free-icons/dog" title="dog icons">Dog icons created by Freepik - Flaticon</a>

