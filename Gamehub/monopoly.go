package gamehub

import (
	db "Monopoly/DB"
	gameplay "Monopoly/Gameplay"
	handler "Monopoly/Handler"
	models "Monopoly/Models"
	"Monopoly/router"
)





func Monopoly(db db.DbOperations, r router.Route) {

	rollDice := gameplay.CreateRollDice(db)
	rollDiceCont := handler.CreateHubContoller(rollDice)
	r.HandleFunc(models.ROLLDICE, rollDiceCont.HandleHub)


	r.Run()

}