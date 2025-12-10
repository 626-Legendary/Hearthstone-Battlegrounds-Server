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

// 用来接收战棋英雄 / 饰品列表
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

		// 官方返回是 []int，这里直接按列表来处理即可
		GameModes []int `json:"gameModes"`

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

// 关键词 metadata 接口返回结构
type KeywordMeta struct {
	ID   int               `json:"id"`
	Slug string            `json:"slug"`
	Name map[string]string `json:"name"`

	Text      map[string]string `json:"text"` // 带 <b> 的富文本
	GameModes []int             `json:"gameModes"`
}

// GetAccessToken 向 https://oauth.battle.net/token 发送 POST 请求，获取 access_token
func GetAccessToken(clientID, clientSecret string) (*TokenResponse, error) {
	endpoint := os.Getenv("OAuth_URL")
	if endpoint == "" {
		endpoint = "https://oauth.battle.net/token"
	}

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

// 自动翻页，返回所有战棋英雄
func GetHeroCards(accessToken string) ([]models.Heroes, error) {
	region := "us" // us / eu / kr / tw
	pageSize := 50
	page := 1

	allHeroes := []models.Heroes{}

	for {
		url := fmt.Sprintf(
			"https://%s.api.blizzard.com/hearthstone/cards?gameMode=battlegrounds&type=hero&pageSize=%d&page=%d",
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
				HSID:        c.ID,
				NameEN:      c.Name["en_US"],
				NameZH:      c.Name["zh_CN"],
				Armor:       c.Armor,
				HeroPowerID: c.Battlegrounds.HeroPowerID,
				CompanionID: c.Battlegrounds.CompanionID,
				ImageEN:     c.Battlegrounds.Image["en_US"],
				ImageZH:     c.Battlegrounds.Image["zh_CN"],
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

// 大饰品
func GetGreaterTrinketsCards(accessToken string) ([]models.Trinkets, error) {
	/*
		https://us.api.blizzard.com/hearthstone/cards?type=trinket&gameMode=battlegrounds&spellSchool=greater_trinket
	*/
	region := "us" // us / eu / kr / tw
	pageSize := 50
	page := 1

	allTrinkets := []models.Trinkets{}

	for {
		url := fmt.Sprintf(
			"https://%s.api.blizzard.com/hearthstone/cards?type=trinket&gameMode=battlegrounds&spellSchool=greater_trinket&pageSize=%d&page=%d",
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

		fmt.Printf("✅ <大饰品数据> 请求成功 page=%d / pageCount=%d, 本页 card=%d, 总计 card=%d\n",
			cr.Page, cr.PageCount, len(cr.Cards), cr.CardCount)

		if len(cr.Cards) == 0 {
			break
		}

		for _, c := range cr.Cards {
			trinket := models.Trinkets{
				HSID:        c.ID,
				NameEN:      c.Name["en_US"],
				NameZH:      c.Name["zh_CN"],
				ManaCost:    c.ManaCost,
				TextEN:      c.Text["en_US"],
				TextZH:      c.Text["zh_CN"],
				ImageEN:     c.Battlegrounds.Image["en_US"],
				ImageZH:     c.Battlegrounds.Image["zh_CN"],
				TrinketType: 2, // 2 = greater
			}
			allTrinkets = append(allTrinkets, trinket)
		}

		if cr.Page >= cr.PageCount {
			break
		}
		page++
	}

	fmt.Printf("✅ 自动翻页完成，总共收集到大饰品数量: %d\n", len(allTrinkets))
	return allTrinkets, nil
}

// 小饰品
func GetLesserTrinketsCards(accessToken string) ([]models.Trinkets, error) {
	/*
		https://us.api.blizzard.com/hearthstone/cards?type=trinket&gameMode=battlegrounds&spellSchool=lesser_trinket
	*/
	region := "us" // us / eu / kr / tw
	pageSize := 50
	page := 1

	allTrinkets := []models.Trinkets{}

	for {
		url := fmt.Sprintf(
			"https://%s.api.blizzard.com/hearthstone/cards?type=trinket&gameMode=battlegrounds&spellSchool=lesser_trinket&pageSize=%d&page=%d",
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

		fmt.Printf("✅ <小饰品数据> 请求成功 page=%d / pageCount=%d, 本页 card=%d, 总计 card=%d\n",
			cr.Page, cr.PageCount, len(cr.Cards), cr.CardCount)

		if len(cr.Cards) == 0 {
			break
		}

		for _, c := range cr.Cards {
			trinket := models.Trinkets{
				HSID:        c.ID,
				NameEN:      c.Name["en_US"],
				NameZH:      c.Name["zh_CN"],
				ManaCost:    c.ManaCost,
				TextEN:      c.Text["en_US"],
				TextZH:      c.Text["zh_CN"],
				ImageEN:     c.Battlegrounds.Image["en_US"],
				ImageZH:     c.Battlegrounds.Image["zh_CN"],
				TrinketType: 1, // 1 = lesser
			}
			allTrinkets = append(allTrinkets, trinket)
		}

		if cr.Page >= cr.PageCount {
			break
		}
		page++
	}

	fmt.Printf("✅ 自动翻页完成，总共收集到小饰品数量: %d\n", len(allTrinkets))
	return allTrinkets, nil
}

// GetKeywords 拉取 keyword metadata，并筛选只在战棋模式(battlegrounds)中使用的关键词
func GetKeywords(accessToken string) ([]models.Keywords, error) {
	/*
		https://us.api.blizzard.com/hearthstone/metadata/keywords
	*/
	region := "us" // us / eu / kr / tw
	url := fmt.Sprintf("https://%s.api.blizzard.com/hearthstone/metadata/keywords", region)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keywords 接口返回状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var metas []KeywordMeta
	if err := json.Unmarshal(bodyBytes, &metas); err != nil {
		return nil, err
	}

	allKeywords := []models.Keywords{}
	const battlegroundsMode = 2 // gameModes 里包含 2 表示战棋

	for _, m := range metas {
		// 只保留战棋相关的关键词
		if !containsInt(m.GameModes, battlegroundsMode) {
			continue
		}

		kw := models.Keywords{
			HSID:   m.ID,
			NameEN: m.Name["en_US"],
			NameZH: m.Name["zh_CN"],

			TextEN: m.Text["en_US"],
			TextZH: m.Text["zh_CN"],
			// Minions 多对多关联先不用填，等后面用 keywordIds 反查再补
		}
		allKeywords = append(allKeywords, kw)
	}

	fmt.Printf("✅ Keyword 元数据拉取完成，战棋相关关键词数量: %d\n", len(allKeywords))
	return allKeywords, nil
}

// 小工具函数：判断切片里是否包含某个 int
func containsInt(list []int, target int) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// _______________________________________________________
// GetQuests 获取战棋任务牌（quest spells）
func GetQuests(accessToken string) ([]models.Quests, error) {
	/*
		https://us.api.blizzard.com/hearthstone/cards?type=spell&gameMode=battlegrounds&sort=quest
	*/
	region := "us" // us / eu / kr / tw
	url := fmt.Sprintf(
		"https://%s.api.blizzard.com/hearthstone/cards?type=spell&gameMode=battlegrounds&sort=quest",
		region,
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

	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quests 接口返回状态码 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 这里是单个 CardsResponse 对象，不是数组
	var cr CardsResponse
	if err := json.Unmarshal(bodyBytes, &cr); err != nil {
		return nil, err
	}

	allQuests := []models.Quests{}

	for _, c := range cr.Cards {
		// 如果你只想要真正标记为 quest 的牌，可以再加个保护：
		// if !c.Battlegrounds.Quest { continue }

		q := models.Quests{
			HSID: c.ID,

			NameEN: c.Name["en_US"],
			NameZH: c.Name["zh_CN"],

			TextEN: c.Text["en_US"],
			TextZH: c.Text["zh_CN"],

			ImageEN: c.Battlegrounds.Image["en_US"],
			ImageZH: c.Battlegrounds.Image["zh_CN"],
		}
		allQuests = append(allQuests, q)
	}

	fmt.Printf("✅ 任务数据拉取完成，战棋任务数量: %d\n", len(allQuests))
	return allQuests, nil
}
