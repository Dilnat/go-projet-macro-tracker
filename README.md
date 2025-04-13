## Installation

### 1. Cloner le projet

```bash
git clone git@github.com:Dilnat/go-projet-macro-tracker.git
cd go-projet-macro-tracker
```
### 2. Démarrer la base de données

```bash
docker-compose up -d
```

### 3. Créer la base de données

```bash
make reset-db
```

Mot de passe : `secret`

### 4.Lancer la CLI
```bash
make run
```

ou 
```bash 
go run main.go
```
