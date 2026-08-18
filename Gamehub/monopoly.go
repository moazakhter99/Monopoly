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

	movePos := gameplay.CreateMovePos(db)
	movePosCont := handler.CreateHubContoller(movePos)
	r.HandleFunc(models.MOVEPOS, movePosCont.HandleHub)

	buyBlock := gameplay.CreateBuyBlock(db)
	buyBlockCont := handler.CreateHubContoller(buyBlock)
	r.HandleFunc(models.BUYBLOCK, buyBlockCont.HandleHub)

	actionCard := gameplay.CreateActionCard(db)
	actionCardCont := handler.CreateHubContoller(actionCard)
	r.HandleFunc(models.ACTIONCARD, actionCardCont.HandleHub)

	jail := gameplay.CreateJail(db)
	jailCont := handler.CreateHubContoller(jail)
	r.HandleFunc(models.JAIL, jailCont.HandleHub)

	changePlayer := gameplay.CreateChangePlayer(db)
	changePlayerCont := handler.CreateHubContoller(changePlayer)
	r.HandleFunc(models.CHANGEPLAYER, changePlayerCont.HandleHub)

	calculateRent := gameplay.CreateCalculateRent(db)
	calculateRentCont := handler.CreateHubContoller(calculateRent)
	r.HandleFunc(models.CALCULATERENT, calculateRentCont.HandleHub)

	payRent := gameplay.CreatePayRent(db)
	payRentCont := handler.CreateHubContoller(payRent)
	r.HandleFunc(models.PAYRENT, payRentCont.HandleHub)


	r.Run()

}