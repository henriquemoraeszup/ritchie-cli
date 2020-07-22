package credential

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const AddNew = "Add a new"

type Settings struct {
	file stream.FileWriteReadExister
}

func NewSettings(file stream.FileWriteReadExister) Settings {
	return Settings{file: file}
}

func (s Settings) ReadCredentialsFields(path string) (Fields, error) {
	fields := Fields{}

	if s.file.Exists(path) {
		cBytes, _ := s.file.Read(path)
		err := json.Unmarshal(cBytes, &fields)
		if err != nil {
			return fields, err
		}
	}

	return fields, nil
}

func (s Settings) ReadCredentialsValue() []ListCredData {
	var creds []ListCredData
	var cred ListCredData
	var detail Detail
	ctx := ctxArr()

	for _, c := range ctx {
		providers := providerByCtx(c)
		for _, p := range providers {
			cBytes, _ := s.file.Read(CredentialsPath() + c + "/" + p)
			_ = json.Unmarshal(cBytes, &detail)
			for k, v := range detail.Credential {
				cred.Provider = detail.Service
				cred.Context = c
				cred.Value = v
				cred.Name = k
				creds = append(creds, cred)
			}
		}
	}
	return creds
}

func providerByCtx(ctx string) []string {
	var providers []string
	files, err := ioutil.ReadDir(CredentialsPath() + "/" + ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		providers = append(providers, f.Name())
	}
	return providers
}

func ctxArr() []string {
	var ctx []string
	files, err := ioutil.ReadDir(CredentialsPath())
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			ctx = append(ctx, f.Name())
		}
	}
	return ctx
}

func (s Settings) WriteCredentialsFields(fields Fields, path string) error {
	fieldsData, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	err = s.file.Write(path, fieldsData)
	if err != nil {
		return err
	}
	return nil
}

// WriteDefault is a non override version of WriteCredentialsFields
// used to create providers.json if user dont have it
func (s Settings) WriteDefaultCredentialsFields(path string) error {
	if !s.file.Exists(path) {
		err := s.WriteCredentialsFields(NewDefaultCredentials(), path)
		return err
	}
	return nil
}

func NewDefaultCredentials() Fields {
	username := Field{
		Name: "username",
		Type: "text",
	}

	token := Field{
		Name: "token",
		Type: "password",
	}

	accessKey := Field{
		Name: "accessKey",
		Type: "text",
	}

	secretAccessKey := Field{
		Name: "secretAccessKey",
		Type: "password",
	}

	base64config := Field{
		Name: "base64config",
		Type: "text",
	}

	dc := Fields{
		AddNew:       []Field{},
		"github":     []Field{username, token},
		"gitlab":     []Field{username, token},
		"aws":        []Field{accessKey, secretAccessKey},
		"jenkins":    []Field{username, token},
		"kubeconfig": []Field{base64config},
	}

	return dc
}

func ProviderPath() string {
	homeDir, _ := os.UserHomeDir()
	providerDir := fmt.Sprintf("%s/.rit/credentials/providers.json", homeDir)
	return providerDir
}

func CredentialsPath() string {
	homeDir, _ := os.UserHomeDir()
	credentialPath := fmt.Sprintf("%s/.rit/credentials/", homeDir)
	return credentialPath
}

func NewProviderArr(fields Fields) []string {
	var providerArr []string
	for k := range fields {
		if k != AddNew {
			providerArr = append(providerArr, k)
		}
	}
	providerArr = append(providerArr, AddNew)
	return providerArr
}
