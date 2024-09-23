package services

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/storage"
	"github.com/Ntrashh/crawlerctl/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
	"log"
	"os"
)

type GitService struct {
	GitStorage storage.GitStorage
}

// NewGitService  创建一个新的 EnvService 实例
func NewGitService(gitStorage storage.GitStorage) *GitService {
	return &GitService{
		GitStorage: gitStorage,
	}
}

func (g GitService) GetGitByProjectID(projectID int) *models.Git {
	git, err := g.GitStorage.GetGitByProjectID(projectID)
	if err != nil {
		return nil
	}
	return git
}

func (g GitService) CreateGit(projectId int, gitPath, username, password, projectPath string) (*models.Git, error) {

	byIdGitConfig := g.GetGitByProjectID(projectId)
	gitConfig := &models.Git{
		ProjectID:           projectId,
		GitPath:             gitPath,
		UserName:            username,
		Password:            password,
		LocalRepositoryPath: projectPath,
	}
	if byIdGitConfig != nil {
		return byIdGitConfig, nil
	}
	err := g.ValidateHTTPConfig(gitConfig)
	if err != nil {
		return nil, err
	}
	err = g.GitStorage.CreateGit(gitConfig)
	if err != nil {
		return nil, err
	}
	return gitConfig, nil
}

// ValidateHTTPConfig 使用 go-git 验证 HTTP Git 配置
func (g GitService) ValidateHTTPConfig(config *models.Git) error {
	tempDir := util.CreateTempDir()
	_, err := g.createRepository(tempDir, "", 1, config)
	if err != nil {
		return fmt.Errorf("验证 HTTP 配置失败: %w", err)
	}
	defer os.RemoveAll(tempDir)
	return nil
}

func (g GitService) createRepository(savePath, remoteName string, depth int, gitConfig *models.Git) (*git.Repository, error) {
	auth := &http.BasicAuth{
		Username: gitConfig.UserName,
		Password: gitConfig.Password,
	}
	if remoteName == "" {
		remoteName = "master"
	}
	repository, err := git.PlainClone(savePath, false, &git.CloneOptions{
		URL:           gitConfig.GitPath,
		Auth:          auth,
		Depth:         depth,
		ReferenceName: plumbing.NewBranchReferenceName(remoteName),
		SingleBranch:  false,
		Tags:          git.NoTags,
	})
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func (g GitService) RemoteBranches(projectId int) (interface{}, error) {
	gitConfig, err := g.GitStorage.GetGitByProjectID(projectId)
	if err != nil {
		return nil, err
	}
	tempDir := util.CreateTempDir()
	defer os.RemoveAll(tempDir)
	repository, err := g.createRepository(tempDir, "", 1, gitConfig)
	if err != nil {
		return nil, err
	}
	refs, err := repository.References()
	if err != nil {
		return nil, fmt.Errorf("无法获取引用: %w", err)
	}
	branches := make([]string, 0)
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsRemote() {
			branchName := ref.Name().Short() // 获取简短名称，如 "origin/main"
			branches = append(branches, branchName)
		}
		return nil
	})
	return branches, nil
}

func (g GitService) GetRemoteBranchCommits(projectId, limit int, branchName string) ([]map[string]interface{}, error) {
	gitConfig, err := g.GitStorage.GetGitByProjectID(projectId)
	if err != nil {
		return nil, err
	}
	tempDir := util.CreateTempDir()
	defer os.RemoveAll(tempDir)
	repository, err := g.createRepository(tempDir, branchName, limit, gitConfig)
	if err != nil {
		return nil, err
	}

	// 获取引用
	ref, err := repository.Reference(plumbing.NewBranchReferenceName(branchName), true)

	if err != nil {
		return nil, fmt.Errorf("无法获取引用: %w", err)
	}
	// 获取提交历史迭代器
	commitIter, err := repository.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("无法获取提交历史: %w", err)
	}

	commits := make([]map[string]interface{}, 0)
	err = commitIter.ForEach(func(c *object.Commit) error {
		commitMsg := map[string]interface{}{
			"hash":       c.Hash.String(),
			"author":     c.Committer.Name,
			"commitTime": c.Committer.When.Format("2006-01-02 15:04:05"),
			"message":    c.Message,
		}
		commits = append(commits, commitMsg)
		return nil
	})
	return commits, nil
}

func (g GitService) BranchPull(projectId int, branchName string) (*models.Git, error) {
	gitConfig, err := g.GitStorage.GetGitByProjectID(projectId)
	if err != nil {
		return nil, err
	}
	newUUID, err := uuid.NewUUID()
	tempDir := fmt.Sprintf("/tmp/%s", newUUID.String())
	defer os.RemoveAll(tempDir)
	err = os.Rename(gitConfig.LocalRepositoryPath, tempDir)
	if err != nil {
		log.Fatalf("移动仓库到目标路径时出错: %v", err)
	}
	_, err = g.createRepository(gitConfig.LocalRepositoryPath, branchName, 1, gitConfig)
	if err != nil {
		err = os.Rename(tempDir, gitConfig.LocalRepositoryPath)
		return nil, err
	}
	return gitConfig, nil
}
