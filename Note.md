安装依赖
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/joho/godotenv
go get gorm.io/datatypes

API
1.随从 Minions
https://us.api.blizzard.com/hearthstone/cards?bgCardType=minion&gameMode=battlegrounds
Query parameters
type：minion
gameMode:battlegrounds

2.英雄 Heroes
https://us.api.blizzard.com/hearthstone/cards?type=hero&gameMode=battlegrounds
Query parameters
type：hero
gameMode:battlegrounds

3. 法术 Spell
https://us.api.blizzard.com/hearthstone/cards?bgCardType=spell&gameMode=battlegrounds
Query parameters
type：spell
gameMode:battlegrounds

3.任务 Quests
https://us.api.blizzard.com/hearthstone/cards?type=spell&gameMode=battlegrounds&sort=quest
Query parameters
type：spell
gameMode:battlegrounds
sort:quest

4.畸变 Anomalies
https://us.api.blizzard.com/hearthstone/cards?bgCardType=anomaly&gameMode=battlegrounds
Query parameters
bgCardType：anomaly
gameMode:battlegrounds

5.奖励 Rewards
https://us.api.blizzard.com/hearthstone/cards?bgCardType=reward&gameMode=battlegrounds
Query parameters
bgCardType：anomaly
gameMode:battlegrounds

6.获取饰品：
https://us.api.blizzard.com/hearthstone/cards
Query parameters
type：trinket
gameMode:battlegrounds

大饰品：
https://us.api.blizzard.com/hearthstone/cards
Query parameters
type：trinket
gameMode:battlegrounds
spellSchool:greater_trinkets

小饰品
https://us.api.blizzard.com/hearthstone/cards
Query parameters
type：trinket
gameMode:battlegrounds
spellSchool:lesser_trinkets


