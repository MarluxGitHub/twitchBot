# TwitchBot

Ein einfacher Twitch-Chat-Bot in Go.

## Features

- Verbindet sich mit dem Twitch IRC-Chat
- Antwortet automatisch auf Chat-Nachrichten
- Konfigurierbar über eine JSON-Datei

## Projektstruktur

```
cmd/bot/           # Einstiegspunkt (main.go)
config/            # Konfigurationsdateien
pkg/twitch/model/  # Konfigurationsmodelle
pkg/twitch/service/# Services (Twitch, ConfigLoader)
```

## Konfiguration

Lege deine Zugangsdaten in `config/config.json` ab:

```json
{
    "twitch": {
        "username": "DeinBotName",
        "oauth": "oauth:dein_token",
        "channel": "#dein_channel"
    }
}
```

> Beispiel: Siehe [config/config.json.dist](config/config.json.dist)

## Starten

1. Abhängigkeiten installieren:
    ```sh
    go mod tidy
    ```
2. Bot starten:
    ```sh
    go run ./cmd/bot
    ```

## Lizenz

MIT