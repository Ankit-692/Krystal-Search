# Krystal-Search

Krystal-Search is a lightning-fast, lightweight desktop utility application built with Go and Wails. It acts as an omnipresent command palette for your Linux system—allowing you to launch applications, locate files instantly, or execute quick terminal tasks from anywhere with a single keystroke.

---

## Features

* **Application Launcher:** Start typing the name of any installed app and hit `Enter` or just `click` any results to open it immediately.
* **File & Folder Search:** Prefix your query with `f:` (e.g., `f:Documents`) to bypass apps and search your system storage directly.
* **Inline Terminal Execution:** Prefix your query with `t:` (e.g., `t:mkdir Projects`) to instantly run terminal commands inside the app context.

## How It Works

* **Instant Access:** Press your custom shortcut (e.g., `Alt + Space` or `Clt + Alt + Space` etc.) anytime to bring Krystal-Search to the front. 
*(Note: You need to manually configure this—see [Step 2: Set Global Shortcut](#step-2-set-global-shortcut) below).*
* **Smart Auto-Hide:** The application window automatically hides itself if you press `Escape` or click anywhere outside the app window.

---

## Guided Tour & Screenshots


### 1. The Core App Launcher
<img width="1920" height="1080" alt="first_look" src="https://github.com/user-attachments/assets/3cf884b1-496c-4e88-b1da-7cc98d74be5d" />
<img width="1920" height="1080" alt="apps" src="https://github.com/user-attachments/assets/5352b657-cccc-469b-b519-2f8a0e2aa6aa" />


### 2. File Search Mode (`f:`)
<img width="1920" height="1080" alt="files_folder" src="https://github.com/user-attachments/assets/9d01e41d-687e-4240-b9be-cd61725fcdf8" />


### 3. Quick Terminal Commands (`t:`)
<img width="1920" height="1080" alt="terminal" src="https://github.com/user-attachments/assets/28023eec-f056-4fcd-8548-acf26792e4c2" />



---

## Installation (Ubuntu / Debian-based)

Getting started is simple. Download the pre-packaged binary directly from our repository releases.

### Step 1: Install the Package
Head over to the **[Releases](https://github.com/Ankit-692/Krystal-Search/releases)** page, download the latest `.deb` package, and install it using `apt` to handle all dependencies automatically:

```bash
sudo apt update
sudo apt install ./Krystal-Search_0.1.0_amd64.deb
```
## Step 2: Set Global Shortcut
Since Krystal-Search launches from anywhere, bind it to a keyboard shortcut (like `Alt + Space`):

1. Go to **Settings** -> **Keyboard Shortcuts**.
2. Scroll down and click **Custom Shortcuts** (or the `+` icon).
3. Set the configuration:
   * **Name:** `Krystal-Search`
   * **Command:** `Krystal-Search`
   * **Shortcut:** *(e.g., `Alt + Space`)*
4. Click **Add**.

---

## Compatibility & Wayland Warning

* **Supported:** Ubuntu and Ubuntu-based distros (Linux Mint, Pop!_OS) running **X11 / Xorg**.
* **Tech Stack:** Go, Wails, JS, and `wmctrl` (for window management).

### Running Wayland? 
Because this app relies on `wmctrl` (an X11 utility), **it will likely break on Wayland**. 

If you are on Wayland and want to help test, please install it and **[Open an Issue](https://github.com/Ankit-692/Krystal-Search/issues)** with what breaks so we can plan future native Wayland support!

---

## 📄 License
MIT License. See [LICENSE](LICENSE) for details.
