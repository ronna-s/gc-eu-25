# Go Design Principles for Great Go Development

## Before we begin
- CLONE ME!!! [github.com/ronna-s/gc-eu-25](https://github.com/ronna-s/gc-eu-25)
- If you haven't already responded, please let me know how much experience you have with Go [here](https://app.sli.do/event/3eTt4zjeY4JEp7moQPpqVc).
- Install Go [here](https://go.dev/dl/)


## The Go type system and interfaces

### Task: build P&P Game
P&P™ stands for Platforms and Programmers.

P&P is a game in which a band of developers tries to take on the ultimate villain: PRODUCTION.

Today you are going to build parts of the game engine, add characters (developers), and define how they interact with PRODUCTION.

Game starts:

The band of developers initially consists of a test character named after your PM, Sir Tan Lee Knot.
Sir Tan Lee Knot has the following skills:
- Pays wages.

The test band of developers initially consists of a test character named after your PM, Sir Tan Lee Knot. You can find it in `cmd/pnp/main.go`.

- The band starts with 10 gold coins, which is used on every turn by the PM to pay the wages, it can be used to buy objects such as coffee or banana (you'll understand why in a sec), pizza, etc._
- PRODUCTION is "calm".

Run the game by executing `go run cmd/pnp/main.go` to see what you have in action.

#### New requirement 0:
- The main loop should skip dead players.
- The game should end if all players are dead (fired or quit).

#### New requirement 1:
- Every band has a minion - implemented in `pkg/pnp/minion.go`. Minions love PRODUCTION, the ultimate villain, and will do anything to serve it. 
- A minion has only one skill - to cause bugs. A minion cannot learn new skills. A minion can be distracted using a banana.
- Add a minion to every game in the game constructor.
- Since a minion cannot die (it's a minion), we need to figure out a way to plug the minion into the game. The minion must not have an `Alive()  bool` method.
- **Note: do not add the `Alive() bool` method to the minion.**
- When your minion is plugged into the game, you will notice that when it's the minion's turn - the title is messed up. Fix it.

#### New requirements 2:
- If all the players are minions, the game needs to end.
- So, we need to tell if a minion is a minion.

#### You're free now to define your own game:
- You can add a gopher who can add features with a chance of hurting production, or fix bugs to make production happy.
- Feature generate gold coins.
- A player dies at random if production reaches legacy.
- Add a satisfaction system to your game.
- A player dies (quits) if their satisfaction reaches 0.
- Add as many more players as you'd like (for instance, a manager who can buy Pizza for everyone - it's supported in the engine already).
- Add leveling up of the players where they have even more skills.
- Whatever strikes your fancy.
- Decide when the game is won - for instance, when you reach 100 gold coins.

