package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const PATH_IN_HOMEDIR = ".lt-service"

const SCRIPTS_URL = "https://static.yyjlincoln.com/scripts/"

func GetLocalServiceDirectory() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		// Also look for that file in the home directory
		LocalService := fmt.Sprintf("%v%c%v", homeDir, os.PathSeparator, PATH_IN_HOMEDIR)
		return LocalService
	}
	return ""
}

func GetLocalServiceCacheDirectory() string {
	LocalService := GetLocalServiceDirectory()
	if LocalService != "" {
		return fmt.Sprintf("%v%c%v", LocalService, os.PathSeparator, "cache")
	}
	return ""
}

func InitLocalServiceDirectories() {
	var LocalService = GetLocalServiceDirectory()
	var LocalCache = GetLocalServiceCacheDirectory()

	// If Local Service does not exist, create it
	if _, err := os.Stat(LocalService); os.IsNotExist(err) {
		os.Mkdir(LocalService, 0755)
	}
	// If Local Cache does not exist, create it
	if _, err := os.Stat(LocalCache); os.IsNotExist(err) {
		os.Mkdir(LocalCache, 0755)
	}
}

func DownloadFileFromCloud(fileName string) error {
	cacheDir := GetLocalServiceCacheDirectory()
	if cacheDir == "" {
		return fmt.Errorf("cache directory not found")
	}

	resp, err := http.Get(SCRIPTS_URL + fileName)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%v%c%v", cacheDir, os.PathSeparator, fileName))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func DownloadCommandFromCloud(command string) error {
	InitLocalServiceDirectories()
	// Try to download binary
	if DownloadFileFromCloud(command) == nil {
		return nil
	}

	// Try and download shell script
	if DownloadFileFromCloud(command+".sh") == nil {
		return nil
	}

	// Try and download python script
	if DownloadFileFromCloud(command+".py") == nil {
		return nil
	}

	return fmt.Errorf("command not found")
}

func ExecuteCommandFromCloud(command string, args []string) (bool, int, error) {
	InitLocalServiceDirectories()
	if DownloadCommandFromCloud(command) == nil {
		return SearchAndExecute(command, args, true)
	}

	return false, 0, fmt.Errorf("command not found")
}
