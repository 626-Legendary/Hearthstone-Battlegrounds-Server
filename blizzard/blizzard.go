// blizzard/blizzard.go
package blizzard

import (
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

// 用来接收战棋英雄列表（这里只取部分字段演示）
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
		Armor         int    `json:"armor"` // 英雄护甲值

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

// FetchBattlegroundHero 用 Bearer token 拉取战棋英雄并打印
func FetchBattlegroundHero(accessToken string) error {
	region := "us" // 也可以换成 eu / kr / tw

	pageSize := 50 // 一页多少条
	page := 1      // 第几页

	url := fmt.Sprintf(
		"https://%s.api.blizzard.com/hearthstone/cards?gameMode=battlegrounds&&type=hero&pageSize=%d&page=%d",
		region, pageSize, page,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cards 接口返回状态码 %d: %s", resp.StatusCode, string(body))
	}

	var cr CardsResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return err
	}

	fmt.Printf("✅ 英雄数据请求成功 page=%d / pageCount=%d, 共 cardCount=%d\n", cr.Page, cr.PageCount, cr.CardCount)

	var heroPowerIdArr []int

	for i, c := range cr.Cards {
		cardID := c.ID
		heroName := c.Name["en_US"]
		heroNameZh := c.Name["zh_CN"]
		heroArmor := c.Armor

		heroPowerId := c.Battlegrounds.HeroPowerID
		companionId := c.Battlegrounds.CompanionID

		img := c.Battlegrounds.Image["en_US"]
		imgZh := c.Battlegrounds.Image["zh_CN"]

		duosOnly := c.Battlegrounds.DuosOnly
		solosOnly := c.Battlegrounds.SolosOnly

		heroPowerIdArr = append(heroPowerIdArr, heroPowerId)

		fmt.Printf("%d. ID=%d, 英雄名称=%s (%s), 护甲=%d, 英雄技能ID=%d, 伙伴ID=%d, 图片=%s (%s), 双人=%v, 单人=%v\n\n",
			i+1, cardID, heroName, heroNameZh, heroArmor, heroPowerId, companionId, img, imgZh, duosOnly, solosOnly)
	}

	fmt.Println(heroPowerIdArr)
	return nil
}
