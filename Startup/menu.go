package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type ShortcutInfo struct {
	IconName       string
	ShortcutName   string
	ShortcutFolder string
	ExecutableName string
}

func main() {
	// Get the path of the currently running file
	currentPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting current file path:", err)
		return
	}

	// Define the list of shortcut information
	shortcutInfoList := []ShortcutInfo{
		{
			IconName:       currentPath + "/clientneco.ico",
			ShortcutName:   currentPath + "Client",
			ShortcutFolder: currentPath + "/Output",
			ExecutableName: currentPath + "/clientneco.exe",
		},
		{
			IconName:       currentPath + "/serverpfp.ico",
			ShortcutName:   currentPath + "Server",
			ShortcutFolder: currentPath + "/Output",
			ExecutableName: currentPath + "/serverneco.exe",
		},
	}

	for _, shortcutInfo := range shortcutInfoList {
		err := createShortcut(currentPath, shortcutInfo)
		if err != nil {
			fmt.Printf("Error creating shortcut: %v\n", err)
			continue
		}
		fmt.Println("Shortcut created successfully!")
	}

	fmt.Println("All shortcuts created.")
}

func createShortcut(currentPath string, shortcutInfo ShortcutInfo) error {
	// Get the current file's directory path
	dir := filepath.Dir(currentPath)

	// Prepare the paths for the icon, shortcut, and executable files
	iconPath := filepath.Join(dir, shortcutInfo.IconName)
	shortcutPath := filepath.Join(dir, shortcutInfo.ShortcutFolder, shortcutInfo.ShortcutName)
	executablePath := filepath.Join(dir, shortcutInfo.ExecutableName)

	// Read the contents of the custom icon file
	iconBytes, err := ioutil.ReadFile(iconPath)
	if err != nil {
		return err
	}

	// Create the shortcut folder if it doesn't exist
	err = os.MkdirAll(filepath.Join(dir, shortcutInfo.ShortcutFolder), 0755)
	if err != nil {
		return err
	}

	// Create a temporary icon file for the shortcut
	tempIconPath := filepath.Join(dir, shortcutInfo.ShortcutFolder, "tempicon.ico")
	err = ioutil.WriteFile(tempIconPath, iconBytes, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tempIconPath) // Clean up the temporary icon file

	// Prepare the command to create the shortcut using PowerShell
	// Change the PowerShell command if you are on a non-Windows platform
	cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf(`$WS = New-Object -ComObject WScript.Shell; $SC = $WS.CreateShortcut('%s'); $SC.TargetPath = '%s'; $SC.IconLocation = '%s'; $SC.Save(); (New-Object -COM Shell.Application).NameSpace('%s').ParseName('%s').InvokeVerb('Open')`, shortcutPath, executablePath, tempIconPath, dir, shortcutInfo.ExecutableName))

	// Execute the PowerShell command
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
