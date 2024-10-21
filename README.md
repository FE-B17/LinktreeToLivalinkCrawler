# Linktree Crawler API

Dieses Projekt stellt eine API bereit, die das Crawlen von Linktree-Profilen ermöglicht. Die API extrahiert Informationen wie Links, Social-Media-Icons, den Titel des Profils und das Profilbild. Die Ergebnisse werden als JSON zurückgegeben und das Profilbild kann heruntergeladen und auf dem Server gespeichert werden.

## Features
- Crawlt Linktree-Profile und extrahiert Links und Social-Media-Icons.
- Lädt das Profilbild herunter und speichert es auf dem Server.
- Gibt alle gecrawlten Daten als JSON über eine API zurück.

## Verwendete Technologien
- **Go**: Programmiersprache für das Backend und die API.
- **Colly**: Web-Scraping-Framework für das Crawlen von Websites.
- **net/http**: Standardbibliothek in Go zum Erstellen eines HTTP-Servers.
- **encoding/json**: Um Datenstrukturen in JSON zu kodieren und zu dekodieren.

## Installation und Ausführung

### Voraussetzungen
- [Go](https://golang.org/doc/install) (ab Version 1.18)
- Internetverbindung (zum Crawlen von Linktree-Profilen)

### Schritte zur Ausführung
1. **Clone dieses Repository**:
   ```bash
   git clone https://github.com/username/LinktreeCrawler.git
   cd LinktreeCrawler
