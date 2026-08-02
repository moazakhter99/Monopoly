# API Info


## Basic API 

1. Create game
- Creates a new game with ID and MatchId
- Add the Player who created the game as the host and the first player
- Response with match Id that can be used by other players to join the game

2. Join Game
- Player who want to join enter the code with their respective info
- Player are matched using the match IDs
- Players info is added to their respective games
- Response with successfull addition to the game

3. Start Game
- Creates the board, dice, Money for the all the players
- Each player get the all the game info to start the game.
- If there is a old game it send the the oldn game info  
- Else it creaats data for a new game
- Respnse witn all the init gamee info for frontend to build the game


## REST based pooling API

1. Dice Roll
- Request for dice roll  
- Rolls the dice
- Response with dice value

2. Player Move
- Calculate the block the Player lands on
- Get the info of the block
- Response with the block the player land on and the actions it can perform

3. Block Buy
- Request to buy the block
- Response with succcess buy and list of cards owned

4. Property Request
- Request to buy Property
- Response with the possible blocks to buy properties on

5. Property Buy
- Request with the list of blocks to buy property
- Response with Sucess

6. Block Sell (GET)
- Request to Sell blocks
- List Of Property to sell

7. Block Sell (POST)
- Request with list of block to be sold
- Response with Success

8. Block Rent
- Request for rent
- Calculate Rent and the player send it to
- Response with Success and rent given

9. Block Action (GET)
- Lands on special block
- Response with list of Action that can be done

10. Block Action (POST)
- Action to be taken
- Perform the action
- Response with the Success of Action

11. Game Update
- Pooling For updates
- Response with game updates


## WebSocket based Events

1. Roll Dice
- Get Req for Roll dice
- Get Rand Value for Dice
- Resp For Dice Roll
- Create Req for move player position

2. Move Position
- Get Req for Move position
- Get the current location of the Player
- Update the player new position on board
- Check the block state for player move
- Resp for Move Pos

3. Buy Property
- Get Req for Buy property
    - If the Player lands on City or Utility Block
    - If Player has enough cash
- Get Player cash
- Update player cash 
- Update Block State and Ownership
- Colour...
- Resp for Buy Property

4. Action Card
- Get Requet for Action Card
- Switch / Case For different Action Type
- Have Updated position & cash depending on action type
- Jail & GoToJail Sends Jail Info with updated position
- Resp for Action Card or Jail

5. Jail
- Get Req for Jail

## Info for Future

#### Things to test
1. Did not test Action Card and Jail Properly
2. Action card is getting initiated twice for single diceroll

### Things to Complete
Could Create Branches for each functionallity to work (Bored in one fuction)

1. Use of Colour for calculating rent
2. Other functionallity like
- Selling Property
- Bankruptcy
- Property Exchange (Future)
- Mortgage (Future)
- Modes (Future)
3. Validation
4. State Management
5. Initialization of Game
- New game
- Old Game
6. Mobile App Support (Future)
- Pooling
- App Backend
7. AI Player (Future)
- Logical
- AI agent
8. Restructuring gameplay (Interface)
- Validation
- ProccessPlay
- Statemanagement
- Response
9. Multiple Games
- Scalable way of running multiple games simultaneously
10. Pub/Sub Model (Future)
- Inbuilt
- Redis
11. Update DB info Correctly



# New Structure
Creating ws Contoller from main function
calling ws Controller in Handle Func
Every user will call this ws and create a ws connection

Inside the Contoller 