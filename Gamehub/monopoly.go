package gamehub

import (
	db "Monopoly/DB"
	gameplay "Monopoly/Gameplay"
	gameroom "Monopoly/Gameroom"
	handler "Monopoly/Handler"
	models "Monopoly/Models"
	"Monopoly/router"
)





func Monopoly(db db.DbOperations, r router.Route, gr gameroom.Room) {

	rollDice := gameplay.CreateRollDice(db, gr)
	rollDiceCont := handler.CreateHubContoller(rollDice)
	r.HandleFunc(models.ROLLDICE, rollDiceCont.HandleHub)


	r.Run()

}