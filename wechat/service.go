package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/internal/util"
	"github.com/qiangmzsx/wechat-mcp/logger"
	"github.com/silenceper/wechat/v2"
	wechatcache "github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	wechatconfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"go.uber.org/zap"
)

// Service 微信服务
type Service struct {
	cfg *config.Config
	wc  *wechat.Wechat
}

// NewService 创建微信服务
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,

		wc: wechat.NewWechat(),
	}
}

// getOfficialAccount 获取公众号实例
func (s *Service) getOfficialAccount() *officialaccount.OfficialAccount {
	memory := wechatcache.NewMemory()
	wechatCfg := &wechatconfig.Config{
		AppID:     s.cfg.WechatAppID,
		AppSecret: s.cfg.WechatAppSecret,
		Cache:     memory,
	}
	return s.wc.GetOfficialAccount(wechatCfg)
}

// UploadMaterialResult 上传素材结果
type UploadMaterialResult struct {
	MediaID   string `json:"media_id"`
	WechatURL string `json:"wechat_url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// UploadMaterial 上传素材到微信
func (s *Service) UploadMaterial(filePath string) (*UploadMaterialResult, error) {
	startTime := time.Now()
	oa := s.getOfficialAccount()
	mat := oa.GetMaterial()

	// 处理文件路径
	localPath, err := DownloadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("download file: %w", err)
	}

	mediaID, url, err := mat.AddMaterial(material.MediaTypeImage, localPath)
	if err != nil {
		logger.Error("upload material failed",
			zap.String("path", filePath),
			zap.Error(err))
		return nil, fmt.Errorf("upload material: %w", err)
	}

	duration := time.Since(startTime)
	logger.Info("material uploaded",
		zap.String("path", filePath),
		zap.String("media_id", util.MaskID(mediaID)),
		zap.Duration("duration", duration))

	return &UploadMaterialResult{
		MediaID:   mediaID,
		WechatURL: url,
	}, nil
}

// CreateDraftResult 创建草稿结果
type CreateDraftResult struct {
	MediaID  string `json:"media_id"`
	DraftURL string `json:"draft_url,omitempty"`
}

// CreateDraft 创建草稿
func (s *Service) CreateDraft(articles []*draft.Article) (*CreateDraftResult, error) {
	startTime := time.Now()
	oa := s.getOfficialAccount()
	dm := oa.GetDraft()
	jsonData, err := json.Marshal(articles)

	logger.Info("CreateDraft params", zap.String("articles", string(jsonData)))

	mediaID, err := dm.AddDraft(articles)
	if err != nil {
		logger.Error("create draft failed", zap.Error(err))
		return nil, fmt.Errorf("create draft: %w", err)
	}

	duration := time.Since(startTime)
	logger.Info("draft created",
		zap.String("media_id", util.MaskID(mediaID)),
		zap.Duration("duration", duration))

	return &CreateDraftResult{
		MediaID:  mediaID,
		DraftURL: fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/appmsg?t=media/appmsg_edit_v2&action=edit&createType=0&token="),
	}, nil
}

// UploadMaterialFromBytes 从字节数据上传素材
func (s *Service) UploadMaterialFromBytes(data []byte, filename string) (*UploadMaterialResult, error) {
	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, "md2wechat_"+filename)
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	defer os.Remove(tmpPath)

	return s.UploadMaterial(tmpPath)
}

// AccessTokenResult 获取 access_token 结果
type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken 获取 access_token
func (s *Service) GetAccessToken() (*AccessTokenResult, error) {
	oa := s.getOfficialAccount()
	accessToken, err := oa.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	return &AccessTokenResult{
		AccessToken: accessToken,
		ExpiresIn:   7200,
	}, nil
}

// UploadMaterialWithRetry 带重试的上传
func (s *Service) UploadMaterialWithRetry(filePath string, maxRetries int) (*UploadMaterialResult, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		result, err := s.UploadMaterial(filePath)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(time.Second)
		}
	}
	return nil, lastErr
}

// DownloadFile 下载文件到临时目录，或返回本地文件路径
func DownloadFile(urlOrPath string) (string, error) {
	if !strings.HasPrefix(urlOrPath, "http://") && !strings.HasPrefix(urlOrPath, "https://") {
		if _, err := os.Stat(urlOrPath); err == nil {
			return urlOrPath, nil
		}
		return "", fmt.Errorf("local file not found: %s", urlOrPath)
	}

	url := urlOrPath
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	tmpDir := os.TempDir()
	ext := ".jpg"
	if parsedURL, err := neturl.Parse(url); err == nil {
		if pathExt := filepath.Ext(parsedURL.Path); pathExt != "" {
			ext = pathExt
		}
	}
	tmpPath := filepath.Join(tmpDir, "md2wechat_download_"+ext)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("write file: %w", err)
	}

	return tmpPath, nil
}

// JSONMarshal 自定义 JSON 序列化
func JSONMarshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// NewspicImageItem 小绿书图片项
type NewspicImageItem struct {
	ImageMediaID string `json:"image_media_id"`
}

// NewspicImageInfo 小绿书图片信息
type NewspicImageInfo struct {
	ImageList []NewspicImageItem `json:"image_list"`
}

// NewspicArticle 小绿书文章
type NewspicArticle struct {
	Title              string           `json:"title"`
	Content            string           `json:"content"`
	ArticleType        string           `json:"article_type"`
	ImageInfo          NewspicImageInfo `json:"image_info"`
	NeedOpenComment    int              `json:"need_open_comment,omitempty"`
	OnlyFansCanComment int              `json:"only_fans_can_comment,omitempty"`
}

// NewspicDraftRequest 小绿书草稿请求
type NewspicDraftRequest struct {
	Articles []NewspicArticle `json:"articles"`
}

// NewspicDraftResponse 微信 API 响应
type NewspicDraftResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MediaID string `json:"media_id"`
}

// CreateNewspicDraft 创建小绿书草稿
func (s *Service) CreateNewspicDraft(articles []NewspicArticle) (*CreateDraftResult, error) {
	startTime := time.Now()

	oa := s.getOfficialAccount()
	accessToken, err := oa.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	req := NewspicDraftRequest{Articles: articles}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/draft/add?access_token=%s", accessToken)

	httpResp, err := http.Post(apiURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("call wechat api: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var resp NewspicDraftResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if resp.ErrCode != 0 {
		logger.Error("create newspic draft failed",
			zap.Int("errcode", resp.ErrCode),
			zap.String("errmsg", resp.ErrMsg))
		return nil, fmt.Errorf("wechat api error: %d - %s", resp.ErrCode, resp.ErrMsg)
	}

	duration := time.Since(startTime)
	logger.Info("newspic draft created",
		zap.String("media_id", util.MaskID(resp.MediaID)),
		zap.Duration("duration", duration))

	return &CreateDraftResult{
		MediaID:  resp.MediaID,
		DraftURL: fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/appmsg?t=media/appmsg_edit_v2&action=edit&createType=0&token="),
	}, nil
}
