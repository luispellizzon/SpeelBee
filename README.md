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
  - **Open/Closed Principle**: add new scoring schemes without modifying the game.
  - **Testability**: inject a deterministic scorer in tests.
  - **Composability**: chain behaviors (e.g., bonus-on-top-of-basic).

---

## 4) Factory — `games.Factory`

**Where**

- `internal/games/factory.go` → `type Factory struct { Dict dict.Repository; Scorer score.Scorer; Board IBoardProvider }`, method `New(kind string) (Game, error)`

**What / Why**

- Centralizes **how game instances are built** (wiring board provider, dictionary, and scorer).
- The caller asks for `New("singleplayer")`; the factory assembles the right combination (`NewPangramSingle(...)`).
- Benefits:
  - **Consistent construction** of complex objects.
  - One place to control dependencies (e.g., swap a scorer or repository for all new games).
  - **Extensibility**: add `"multiplayer"` later without changing callers, or any other type of game.

---

## 5) Singleton — `pangram.Board()`

**Where**

- `internal/pangram/board.go`
  - Global `GameBoard` with `sync.Once` and `Board() (GameBoard, error)`
  - `InitSource(s Source)` sets the **single** source of truth for today’s board
- `internal/pangram/provider.go` → `Provider.Board()` delegates to the singleton

**What / Why**

- The game board for “today” is **created once** and shared globally.
- `sync.Once` ensures **exactly-one** initialization even under concurrency.
- Benefits:
  - **Consistency**: all sessions see the same board for the day.
  - **Safety**: thread-safe lazy initialization; avoids races.
  - **Control**: source can be injected for tests or different environments.

---

## (Bonus) Lightweight Decorator — `games.pangramSingle`

**Where**

- `internal/games/pangram_singleplayer.go`

**What / Why**

- Wraps a core `Game` and **adjusts behavior** (e.g., alters `Name()` to tag “SINGLE PLAYER”, delegates other calls).
- Not requested in the prompt, but worth noting as a clean wrapper that extends behavior without modifying the original game.

---

## How They Work Together

- **Factory** builds games using a **Board** from the **Singleton**, a dictionary via the **Adapter** (optionally behind the **Proxy**), and a scoring **Strategy**.
- This yields a system that is **modular, testable, and easy to evolve**:
  - Swap the dictionary backend by replacing the adapter target.
  - Tune performance by inserting or removing the proxy.
  - Change game feel by switching scoring strategies.
  - Keep everyone consistent with a single daily board.

---

## Tradeoffs & Notes

- **Singletons** can hinder parallel experiments if you want different boards simultaneously; consider scoping via contexts in the future.
- **Factory** `switch kind` is simple; if kinds multiply, consider a registration map to avoid a long switch.
- **Proxy** is FIFO-eviction; for hot-word workloads, an LRU might improve hit rate.
- **Adapter** is currently file-based; abstract I/O (e.g., interfaces for readers) can ease migration to remote stores.
