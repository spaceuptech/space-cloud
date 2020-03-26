package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

type alphabetic []string

func (list alphabetic) Len() int { return len(list) }

func (list alphabetic) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

func (list alphabetic) Less(i, j int) bool {
	a, err := strconv.Atoi(strings.Split(list[i], "-")[0])
	if err != nil {
		logrus.Errorf("unable to convert string to int while sorting file name (i) (%s)", list[i])
		return false
	}
	b, err := strconv.Atoi(strings.Split(list[j], "-")[0])
	if err != nil {
		logrus.Errorf("unable to convert string to int while sorting file name (j) (%s)", list[j])
		return false
	}
	return a < b
}

// Apply reads the config file(s) from the provided file / directory and applies it to the server
func Apply(applyName string) error {
	if !strings.HasSuffix(applyName, ".yaml") {
		dirName := applyName
		if err := os.Chdir(dirName); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to switch to directory %s", dirName), err)
		}
		// list file of current directory. Note : we have changed the directory with the help of above function
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to fetch config files from %s", dirName), err)
		}

		account, err := utils.GetSelectedAccount()
		if err != nil {
			return utils.LogError("Unable to fetch account information", err)
		}
		login, err := utils.Login(account)
		if err != nil {
			return utils.LogError("Unable to login", err)
		}

		fileNames := alphabetic{}
		// filter directories
		for _, fileInfo := range files {
			if !fileInfo.IsDir() {
				fileNames = append(fileNames, fileInfo.Name())
			}
		}
		// sort file names alphanumerically
		sort.Stable(fileNames)

		for _, fileName := range fileNames {
			if strings.HasSuffix(fileName, ".yaml") {
				specs, err := utils.ReadSpecObjectsFromFile(fileName)
				if err != nil {
					return utils.LogError("Unable to read spec objects from file", err)
				}

				// Apply all spec
				for _, spec := range specs {
					if err := ApplySpec(login.Token, account, spec); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	account, err := utils.GetSelectedAccount()
	if err != nil {
		return utils.LogError("Unable to fetch account information", err)
	}
	login, err := utils.Login(account)
	if err != nil {
		return utils.LogError("Unable to login", err)
	}

	specs, err := utils.ReadSpecObjectsFromFile(applyName)
	if err != nil {
		return utils.LogError("Unable to read spec objects from file", err)
	}

	// Apply all spec
	for _, spec := range specs {
		if err := ApplySpec(login.Token, account, spec); err != nil {
			return err
		}
	}

	return nil
}

// ApplySpec takes a spec object and applies it
func ApplySpec(token string, account *model.Account, specObj *model.SpecObject) error {
	requestBody, err := json.Marshal(specObj.Spec)
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("error while applying service unable to marshal spec - %s", err.Error()), nil)
		return err
	}
	url, err := adjustPath(fmt.Sprintf("%s%s", account.ServerURL, specObj.API), specObj.Meta)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("error while applying service unable to send http request - %s", err.Error()), nil)
		return err
	}

	v := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	utils.CloseTheCloser(req.Body)
	if resp.StatusCode != 200 {
		_ = utils.LogError(fmt.Sprintf("error while applying service got http status code %s - %s", resp.Status, v["error"]), nil)
		return fmt.Errorf("%v", v["error"])
	}
	utils.LogInfo(fmt.Sprintf("Successfully applied %s", specObj.Type))
	return nil
}

func adjustPath(path string, meta map[string]string) (string, error) {
	newPath := path
	for {
		pre := strings.IndexRune(newPath, '{')
		if pre < 0 {
			return newPath, nil
		}
		post := strings.IndexRune(newPath, '}')

		key := strings.TrimSuffix(strings.TrimPrefix(newPath[pre:post], "{"), "}")
		value, p := meta[key]
		if !p {
			return "", fmt.Errorf("provided key (%s) does not exist in metadata", key)
		}

		newPath = newPath[:pre] + value + newPath[post+1:]
	}
}
