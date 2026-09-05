[English](README.md) | [Русский](README.ru.md)
# PortKiller 🔪

CLI-утилита для управления процессами через занятые порты. Кроссплатформенная, один бинарник без зависимостей.

## Возможности

- 📋 Просмотр всех занятых TCP-портов с указанием процесса
- 💀 Убийство процесса по номеру порта
- 🔍 Убийство процессов по имени (поддержка подстрок и регистра)

## Установка

### Способ 1: Через Go (рекомендуется)

Требуется **Go 1.22+**:

```bash
go install github.com/Awak3r/PortKiller@latest
```

Добавь Go-бинарники в PATH (одноразово):

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc   # или ~/.bashrc
source ~/.zshrc
```

### Способ 2: Готовые бинарники (без Go)

1. Перейди в [Releases](https://github.com/Awak3r/PortKiller/releases)
2. Скачай архив под свою ОС:
   - Linux: `portkiller-linux-amd64.tar.gz`
   - macOS: `portkiller-darwin-arm64.tar.gz` (Apple Silicon) или `-amd64` (Intel)
   - Windows: `portkiller-windows-amd64.zip`
3. Распакуй и положи в директорию из PATH:

```bash
# Linux/macOS
sudo mv portkiller /usr/local/bin/
sudo chmod +x /usr/local/bin/portkiller

# Windows
# Скопируй portkiller.exe в C:\Windows\ или добавь папку в PATH
```

### Способ 3: Сборка из исходников

```bash
git clone https://github.com/Awak3r/PortKiller.git
cd PortKiller
go build -o portkiller .
sudo mv portkiller /usr/local/bin/
```

## Использование

> ⚠️ **Важно:** утилита требует прав администратора (sudo) — она автоматически попросит пароль при запуске.

### Список занятых портов

```bash
PortKiller list
```

Выводит таблицу: `PORT | PID | NAME`

### Убийство по порту

```bash
PortKiller kill -port 5000
```

### Убийство по имени

```bash
PortKiller kill -name node
```

## Флаги

| Флаг | Описание |
|------|----------|
| `-port <число>` | Номер порта для убийства (1–65535) |
| `-name <строка>` | Имя процесса (подстрока, без учёта регистра) |
| `-h` / `--help` | Справка |
| `-v` / `--version` | Версия |

## Примеры

```bash
# Найти, кто занял порт 3000
PortKiller list | grep 3000

# Убить все Node-процессы
PortKiller kill -name node

# Убить сервер разработки без вопросов
PortKiller kill -port 3000 -f
```

