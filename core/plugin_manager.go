package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	waLog "go.mau.fi/whatsmeow/util/log"
)

var appLog waLog.Logger = AppLog()

func AddPlugin(FileName string, Url string) error {

	pluginExists := _CheckIfPluginExists(FileName)

	if pluginExists {
		return errors.New("pl_exists")
	}

	body, downloadError := _DownloadPlugin(Url)
	if downloadError != nil {
		return downloadError
	}
	writeError := _WritePlugin(FileName, *body)
	if writeError != nil {
		return writeError
	}

	databaseError := _AddToDatabase(FileName, Url)

	if databaseError != nil {
		return databaseError
	}

	return nil
}

func _CheckIfPluginExists(FileName string) bool {
	db := _CheckIfPluginExistsInDatabase(FileName)
	fileSystem := _CheckIfPluginExistsInFilesystem(FileName)

	if !db && !fileSystem {
		return false
	}

	if db && !fileSystem {
		// Db exists, but not in filesystem. delete from db.
		_DeletePluginFromDatabase(FileName)
		return false
	}

	if fileSystem && !db {
		// exists as file, but not in db. delete from filesystem.
		filePath := path.Join("pl", FileName)
		_DeletePluginFromFilesystem(filePath)
		return false
	}

	return true

}

func _CheckIfPluginExistsInDatabase(FileName string) bool {
	db := GetDatabase()

	var plugin Plugin
	err := db.Where("name = ?", FileName).First(&plugin)

	if err.Error != nil {
		return false
	}

	return true
}

func _CheckIfPluginExistsInFilesystem(FileName string) bool {
	filePath := path.Join("pl", FileName)
	_, err := os.Stat(filePath)

	if err != nil {
		return false
	}

	return true
}

func _AddToDatabase(fileName string, url string) error {
	database := GetDatabase()
	var plugin Plugin

	err := database.Where(Plugin{Name: fileName, Url: url}).Attrs(Plugin{
		Name: fileName,
		Url:  url,
	}).FirstOrCreate(&plugin)

	if err.Error != nil {
		appLog.Errorf("PluginManager:AddToDatabase - Error adding plugin to database ", fileName, url)
		return err.Error
	}

	return nil
}

func _DownloadPlugin(url string) (*io.ReadCloser, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)

	if err != nil {
		appLog.Errorf("PluginManager:DownloadPlugin - Error while downloading ", url)
		return nil, err
	}

	return &resp.Body, nil
}

func _WritePlugin(fileName string, body io.ReadCloser) error {
	out, err := os.Create(fmt.Sprintf("./pl/%s", fileName))

	if err != nil {
		appLog.Errorf("PluginManager:WritePlugin - Error while trying to create ", fileName)
		return err
	}

	_, copyErr := io.Copy(out, body)

	if copyErr != nil {
		appLog.Errorf("PluginManager:WritePluginCopy - Error while trying to create ", fileName)
		return err
	}

	defer out.Close()
	defer body.Close()

	return nil
}

func DeletePlugin(path string, fileName string) error {
	fsError := _DeletePluginFromFilesystem(path)

	if fsError != nil {
		return fsError
	}

	dbDeleteErr := _DeletePluginFromDatabase(fileName)

	if dbDeleteErr != nil {
		return dbDeleteErr
	}

	return nil
}

func _DeletePluginFromFilesystem(FilePath string) error {
	err := os.Remove(FilePath)

	if err != nil {
		appLog.Errorf("PluginManager:DeletePluginFromFilesystem - Could not delete", err)
		return err
	}

	return nil
}

func _DeletePluginFromDatabase(fileName string) error {
	db := GetDatabase()

	err := db.Unscoped().Where("name = ?", fileName).Delete(Plugin{})

	if err.Error != nil {
		appLog.Errorf("PluginManager:DeletePluginFromDatabase - Could not delete", fileName)

		return err.Error
	}

	return nil

}
