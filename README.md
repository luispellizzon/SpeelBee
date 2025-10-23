# Design Patterns Used in the Pangram Game Codebase

This document explains **which design patterns are used**, **where they appear in the codebase**, and **why they were the best option**. Examples reference the file and type names from the repository.

---

## 1) Adapter — `dict.JSONAdapter`

**Where**

- `internal/dict/repo.go` → `type JSONAdapter struct{ inner *jsonMap }` with `Has(word string) (bool, error)`

**What / Why**

- The app wants a **uniform dictionary interface**: `Repository{ Has(word string) (bool, error) }`.
- Real data is stored as a JSON map loaded from disk (`jsonMap`). The JSON format is not the interface the game consumes.
- `JSONAdapter` **adapts** the concrete `jsonMap` to the `Repository` interface so game logic can check words without knowing how the data is stored.
- Benefits:
  - Clean separation between **data format** and **domain logic**.
  - Easy to swap to a different repository later (e.g., .TXT file) without touching game code.

---

## 2) Proxy — `dict.CacheProxy`

**Where**

- `internal/dict/repo.go` → `type CacheProxy struct { repo Repository; cache map[string]bool; ... }`

**What / Why**

- `CacheProxy` wraps a `Repository` and **intercepts** `Has` calls.
- It **caches** results (bounded by `capacity`) and returns cached answers quickly, logging the source (`FROM CACHE` vs `FROM DATABASE`) to inform us where the word was retrieved from.
- Benefits:
  - **Performance**: repeated lookups are fast and the request do not need to be sent out to deeper layers in the codebase to access database parts.
  - **Fast retrieve**: word submitted by user is checked in the first layer of the game logic and if the word was not submitted yet, it checks the repository.

---

## 3) Strategy — `score.Scorer` (+ `BasicScorer`, `BonusScorer`)

**Where**

- `internal/score/scorer.go` → `type Scorer interface { Score(length int, pangram bool) int }`
- Implementations: `BasicScorer`, `BonusScorer`
- Used by: `internal/games/pangram.go` → `pangramGame` holds a `score.Scorer`

**What / Why**

- The game delegates **how points are calculated** to a pluggable `Scorer` strategy.
- `BasicScorer` encodes the baseline rule; `BonusScorer` **decorates** another scorer and adds pangram bonuses.
- Benefits:
  - **Open/Closed Principle**: add new scoring schemes without modifying the game, extending for different type of bonuses.

---

## 4) Factory — `games.Factory`

**Where**

- `internal/games/factory.go` → `type Factory struct { Dict dict.Repository; Scorer score.Scorer; Board IBoardProvider }`, method `New(kind string) (Game, error)`

**What / Why**

- Centralizes **how game instances are built** (wiring board provider, dictionary, and scorer).
- The caller asks for `New("singleplayer")`; the factory code assembles the right kind of game (`NewPangramSingle(...)`).
- Benefits:
  - **Consistent construction** of complex objects.
  - One place to control dependencies (e.g., swap a scorer or repository for all new games).
  - **Extensibility**: I can add `"multiplayer"` later without changing callers.

---

## 5) Singleton — `pangram.Board()`

**Where**

- `internal/pangram/board.go`
  - Global `GameBoard` with `sync.Once` and `Board() (GameBoard, error)`
  - `InitSource(s Source)` sets the **single** source of truth for today’s board
- `internal/pangram/provider.go` → `Provider.Board()` delegates to the singleton

**What / Why**

- The game board for “today” is **created once** and shared globally for all games created. When the user join the server and create a new game, this game will fetch information from the board.
- `sync.Once` ensures **exactly-one** initialization.
- Benefits:
  - **Consistency**: all sessions see the same board for the day.
  - **Safety**: thread-safe lazy initialization; avoids races.

---

## 6) Singleton — `mgr (Manager)`

**Where**

- `internal/manager/manager.go`
  - Global `Manager` with `sync.Once`
  - Manager is a singleton that acts as a proxy to create new games and save each game in a map so users can also rejoin their game with their game_id using manager's `Get(id string)`.

**What / Why**

- The manager is global to the server and is responsible to manage new games and old games created. Requests passes through the manager class to either create a new game, or to search a game that was already created.
- `sync.Once` ensures **exactly-one** initialization.
- Benefits:
  - **Consistency**: The server sees only one manager that manages all games from the server.

---

## 7) Singleton — `logger Logger`

**Where**

- `internal/logger/logger.go`
  - Global `Logger` with `sync.Once`
  - Logger is a singleton that logs information about the server. Only one instance is created and can be used across the codebase to log logs with different types: INFO/ERROR

**What / Why**

- The logger is global to the server and is responsible to log logs across the project.
- `sync.Once` ensures **exactly-one** initialization.
- Benefits:
  - **Consistency**: The server sees only one logger instance, and its purpose is to log different types of logs according to their type (INFO/ERROR). It is extend to accept any type of interface, but mainly Errors.

---

## 8) SinglePlayer Decorator wrapper — `games.pangramSingle`

**Where**

- `internal/games/pangram_singleplayer.go`

**What / Why**

- Wraps a core `Game` and **adjusts behavior** (e.g., alters `Name()` to tag “SINGLE PLAYER”, delegates other calls).
- Not requested in the prompt, but worth noting as a clean wrapper that extends behavior without modifying the original game.
- The wrapper can implement others interface to meet requirements for the type of the game. For example, for future multiplayer game, I can create a pangramMultiplayer wrapper and integrate a interface called WithPlayers to gather players on the same game session.
- Benefits:
  - **Extensibility**: each type of game can have added interfaces specific to their purpose.

---

## How The Game is Built

- **Factory** builds games using a **Board** from the **Singleton**, a dictionary via the **Adapter** (optionally behind the **Proxy**), and a scoring **Strategy**.
- This yields a system that is \*_modular_:
  - Swap the dictionary backend by replacing the adapter target.
  - Tune word submission performance by inserting or removing the proxy.
  - Change game points by switching scoring strategies.
  - Keep everyone consistent with a single daily board.

---
