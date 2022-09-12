package utils

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

func Update() error {

	osType := runtime.GOOS
	arch := runtime.GOARCH

	v, err := http.Get("https://api.github.com/repos/Faagerholm/clockify-cli/releases/latest")
	if err != nil {
		return err
	}
	defer v.Body.Close()

	info := struct {
		TagName string  `json:"tag_name"`
		Assets  []asset `json:"assets"`
	}{}

	err = json.NewDecoder(v.Body).Decode(&info)
	if err != nil {
		return err
	}

	if strings.Compare(info.TagName, getCurrentVersion()) == 0 {
		fmt.Println("You are running the latest version of clockify-cli")
		return nil
	}

	var osAss []asset
	for _, asset := range info.Assets {
		if strings.Contains(asset.Name, osType) && strings.Contains(asset.Name, arch) {
			osAss = append(osAss, asset)
		}
	}

	if len(osAss) == 0 {
		fmt.Printf("No release found for %s\n", osType)
		return nil
	}

	home := getInstallDir()
	err = downloadFiles(osAss, home+"/bin")
	if err != nil {
		return err
	}
	// update version file
	err = ioutil.WriteFile(home+"/.version", []byte(info.TagName), 0755)
	if err != nil {
		return err
	}

	fmt.Printf("DONE! Updated to version %s\n", info.TagName)
	return nil
}

func getCurrentVersion() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %s", err)
		return ""
	}

	file, err := ioutil.ReadFile(home + "/.clockify-cli/.version")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(file))
}

func downloadFiles(assets []asset, dest string) error {
	// var md5 []byte
	for _, asset := range assets {
		if strings.Contains(asset.Name, "md5") {
			// handle md5
		} else if strings.Contains(asset.Name, "tar.gz") {
			resp, err := http.Get(asset.BrowserDownloadUrl)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			gzr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			defer gzr.Close()
			err = untar(gzr, dest)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func untar(reader io.Reader, dest string) error {
	tarReader := tar.NewReader(reader)
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := dest + header.Name
		switch header.Typeflag {
		case tar.TypeReg:
			// overwrite existing files
			if _, err := os.Stat(target); err == nil {
				err = os.Remove(target)
				if err != nil {
					return err
				}
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
			file.Close()
		}
	}
}

func getInstallDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if strings.Contains(ex, "go-build") {
		// we are in dev mode
		h, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return h
	}

	s := strings.Split(ex, "/")
	return strings.Join(s[:len(s)-2], "/")
}
