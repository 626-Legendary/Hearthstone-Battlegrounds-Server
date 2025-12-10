// blizzard/blizzard.go
package blizzard

import (
	"bgs-server/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// 用来接收 OAuth 返回的 access_token
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// 用来接收战棋英雄列表
type CardsResponse struct {
	CardCount int `json:"cardCount"`
	PageCount int `json:"pageCount"`
	Page      int `json:"page"`
	Cards     []struct {
		ID            int    `json:"id"`
		Collectible   int    `json:"collectible"`
		Slug          string `json:"slug"`
		ClassID       int    `json:"classId"`
		MultiClassIds []int  `json:"multiClassIds"`
		CardTypeID    int    `json:"cardTypeId"`
		CardSetID     int    `json:"cardSetId"`
		RarityID      int    `json:"rarityId"`
		ArtistName    string `json:"artistName"`
		Health        int    `json:"health"`
		ManaCost      int    `json:"manaCost"`
		Armor         int    `json:"armor"`

		Name       map[string]string `json:"name"`
		Text       map[string]string `json:"text"`
		FlavorText map[string]string `json:"flavorText"`
		Image      map[string]string `json:"image"`
		ImageGold  map[string]string `json:"imageGold"`

		CropImage string `json:"cropImage"`
		ChildIDs  []int  `json:"childIds"`

		IsZilliaxFunctionalModule bool `json:"isZilliaxFunctionalModule"`
		IsZilliaxCosmeticModule   bool `json:"isZilliaxCosmeticModule"`

		Battlegrounds struct {
			Hero        bool `json:"hero"`
			HeroPowerID int  `json:"heroPowerId"`
			Quest       bool `json:"quest"`
			Reward      bool `json:"reward"`
			DuosOnly    bool `json:"duosOnly"`
			SolosOnly   bool `json:"solosOnly"`
			CompanionID int  `json:"companionId"`

			Image     map[string]string `json:"image"`
			ImageGold map[string]string `json:"imageGold"`
		} `json:"battlegrounds"`
	} `json:"cards"`
}

// GetAccessToken 向 https://oauth.battle.net/token 发送 POST 请求，获取 access_token
func GetAccessToken(clientID, clientSecret string) (*TokenResponse, error) {
	endpoint := os.Getenv("OAuth_URL")

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Basic Auth: base64(clientID:clientSecret)
	basic := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+basic)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token 接口返回状态码 %d: %s", resp.StatusCode, string(body))
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	return &tr, nil
}

// GetBattlegroundHero 自动翻页，返回所有战棋英雄
func GetBattlegroundHero(accessToken string) ([]models.Heroes, error) {
	region := "us" // us / eu / kr / tw
	pageSize := 50
	page := 1

	allHeroes := []models.Heroes{}

	for {
		url := fmt.Sprintf(
			"https://%s.api.blizzard.com/hearthstone/cards?gameMode=battlegrounds&&type=hero&pageSize=%d&page=%d",
			region, pageSize, page,
		)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		// 注意：在 for 里 defer 会堆积，这里直接 Close 更安全
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("cards 接口返回状态码 %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var cr CardsResponse
		if err := json.Unmarshal(bodyBytes, &cr); err != nil {
			return nil, err
		}

		fmt.Printf("✅ 英雄数据请求成功 page=%d / pageCount=%d, 本页 card=%d, 总计 card=%d\n",
			cr.Page, cr.PageCount, len(cr.Cards), cr.CardCount)

		if len(cr.Cards) == 0 {
			break
		}

		for _, c := range cr.Cards {
			hero := models.Heroes{
				HSID: c.ID,
				// 现在 API 只有一个 name（locale = zh_CN），先当成中文名
				NameEN:      c.Name["en_US"], // 以后可以再补英文
				NameZH:      c.Name["zh_CN"], // 当前语言是 zh_CN
				Armor:       c.Armor,
				HeroPowerID: c.Battlegrounds.HeroPowerID,
				CompanionID: c.Battlegrounds.CompanionID,
				ImageEN:     c.Battlegrounds.Image["en_US"], // 以后可以用 en_US 再查
				ImageZH:     c.Battlegrounds.Image["zh_CN"], // 当前语言图片
				IsDuo:       c.Battlegrounds.DuosOnly,
				IsSolo:      c.Battlegrounds.SolosOnly,
			}
			allHeroes = append(allHeroes, hero)
		}

		if cr.Page >= cr.PageCount {
			break
		}
		page++
	}

	fmt.Printf("✅ 自动翻页完成，总共收集到英雄数量: %d\n", len(allHeroes))
	return allHeroes, nil
}
