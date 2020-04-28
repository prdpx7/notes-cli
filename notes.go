package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const notesFilePathPrefix = "daily_notes_"

// GistNotesMap - Offline mapping of local notes files with respective remote gist-ids
type GistNotesMap struct {
	GistID   string `json:"gist_id"`
	Filename string `json:"filename"`
}

func getOrCreateNotesDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dirPath := homeDir + "/.notes/"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dirPath
}
func getOrCreateNotesDataDir() string {
	dirPath := getOrCreateNotesDir()
	dataDirPath := dirPath + "data/"
	if _, err := os.Stat(dataDirPath); os.IsNotExist(err) {
		err := os.Mkdir(dataDirPath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dataDirPath
}

func getOrCreateNotesConfigDir() string {
	dirPath := getOrCreateNotesDir()
	configDirPath := dirPath + "config/"
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		err := os.Mkdir(configDirPath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return configDirPath
}

func getOrCreateLocalGistStore() string {
	configDirPath := getOrCreateNotesConfigDir()
	localGistStore := configDirPath + "gist_store.json"
	_, err := os.Open(localGistStore)
	if err != nil {
		os.Create(localGistStore)
	}
	return localGistStore
}

func getGithubPersonalToken() (token string) {
	token, exists := os.LookupEnv("GITHUB_PERSONAL_TOKEN")
	if exists == false {
		fmt.Println(`No GITHUB_PERSONAL_TOKEN found..!
You can get your personal token from https://github.com/settings/tokens/
Paste the following line in your ~/.bashrc or ~/.zshrc
export GITHUB_PERSONAL_TOKEN='asdaspersonaltoken123'`,
		)
		os.Exit(1)
	}
	return token
}

func getAllLocalNotesFiles() []string {
	dataDirPath := getOrCreateNotesDataDir()
	markdownsFiles, err := filepath.Glob(dataDirPath + "*.md")
	if err != nil {
		log.Fatal(err)
	}
	return markdownsFiles
}

func getMarkdownFilename(filePath string) string {
	tmp := strings.Split(filePath, "/")
	return tmp[len(tmp)-1]
}

func getOrCreateGist(token string, markdownFilePath string, GistID string) string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	markdownFileName := getMarkdownFilename(markdownFilePath)
	markdownFileContentByte, err := ioutil.ReadFile(markdownFilePath)
	markdownFileContent := string(markdownFileContentByte)
	if err != nil {
		log.Fatal(err)
	}
	tmpGistFile := github.GistFile{Filename: &markdownFileName, Content: &markdownFileContent}
	var tmpFilesObj = map[github.GistFilename]github.GistFile{
		github.GistFilename(markdownFileName): tmpGistFile,
	}
	var gistVisibilityToPublic = false
	var gistDescription = "Daily Notes via notes app"
	// create gist
	if GistID == "" {
		tmpGistObj := github.Gist{
			Files:       tmpFilesObj,
			Public:      &gistVisibilityToPublic,
			Description: &gistDescription,
		}
		gistResponse, _, err := client.Gists.Create(ctx, &tmpGistObj)
		if err != nil {
			log.Fatal(err)
		}
		return *gistResponse.ID
	}
	tmpGistObj := github.Gist{
		Files:       tmpFilesObj,
		Public:      &gistVisibilityToPublic,
		Description: &gistDescription, ID: &GistID,
	}
	gistResponse, _, err := client.Gists.Edit(ctx, GistID, &tmpGistObj)
	if err != nil {
		log.Fatal(err)
	}
	return *gistResponse.ID
}

func DoCloudSync() {
	fmt.Println("Syncing your gists....")
	loader := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	loader.Start()
	token := getGithubPersonalToken()
	localGistStorePath := getOrCreateLocalGistStore()
	gistStoreFile, _ := ioutil.ReadFile(localGistStorePath)
	data := []GistNotesMap{}
	_ = json.Unmarshal([]byte(gistStoreFile), &data)

	localMarkdownsFiles := getAllLocalNotesFiles()

	var filesToBeSynced []GistNotesMap
	for i := 0; i < len(localMarkdownsFiles); i++ {
		var fileFoundInStore = false
		for j := 0; j < len(data); j++ {
			if data[j].Filename == localMarkdownsFiles[i] {
				filesToBeSynced = append(filesToBeSynced, data[j])
				fileFoundInStore = true
				break
			}
		}
		if fileFoundInStore == false {
			filesToBeSynced = append(filesToBeSynced, GistNotesMap{"", localMarkdownsFiles[i]})
		}
	}
	for i := 0; i < len(filesToBeSynced); i++ {
		filesToBeSynced[i].GistID = getOrCreateGist(token, filesToBeSynced[i].Filename, filesToBeSynced[i].GistID)
	}

	syncedFileData, _ := json.MarshalIndent(filesToBeSynced, "", " ")

	_ = ioutil.WriteFile(localGistStorePath, syncedFileData, 0644)
	loader.Stop()
	fmt.Println("Done!")
}

func RunEditor(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func isVimEditor(editor string) bool {
	re := regexp.MustCompile("(g|n|neo)?vim")

	if re.Match([]byte(editor)) {
		return true
	}
	return false
}

func GetWorkingTextEditorWithFileBrowsingSupport() string {
	commands := [...]string{"vim", "vi", "nvim", "emacs"}
	prefixPath := [...]string{"/usr/bin/", "/usr/local/bin/"}
	for i := 0; i < len(commands); i++ {
		for j := 0; j < len(prefixPath); j++ {
			cmdFilePath := prefixPath[j] + commands[i]
			if _, err := os.Stat(cmdFilePath); !os.IsNotExist(err) {
				return cmdFilePath
			}
		}
	}
	fmt.Println("Unable to find terminal based filebrowser like vim, emacs, nvim, vi etc..!")
	log.Fatal("Please install any of the above software to browse notes in terminal")
	return ""
}

func GetEditorCommand(editor string, mode string) *exec.Cmd {
	// fmt.Println("editor",editor)
	currentTime := time.Now()
	// fmt.Println(currentTime.Format("2006-01-02"))
	var filename string
	dataDirPath := getOrCreateNotesDataDir()
	if mode == "write" {
		filename = dataDirPath + notesFilePathPrefix + currentTime.Format("2006_01_02") + ".md"
	} else {
		filename = dataDirPath
	}
	// fmt.Println(filename)

	if isVimEditor(editor) && mode == "write" {
		cmd := exec.Command(editor, "+normal Go", filename)
		return cmd
	}
	return exec.Command(editor, filename)
}
