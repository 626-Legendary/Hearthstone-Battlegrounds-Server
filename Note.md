安装依赖
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/joho/godotenv

bgs-backend/
  main.go
  database/
    database.go
  models/
    card.go
    keyword.go
  handlers/
    card_handler.go
  routes/
    routes.go

获取饰品：
https://us.api.blizzard.com/hearthstone/cards?locale=zh_CN&type=trinket&gameMode=battlegrounds

大饰品：
https://us.api.blizzard.com/hearthstone/cards?locale=zh_CN&type=trinket&gameMode=battlegrounds&spellSchool=greater_trinkets

小饰品
https://us.api.blizzard.com/hearthstone/cards?locale=zh_CN&type=trinket&gameMode=battlegrounds&spellSchool=lesser_trinket