# GoBookMarker

A native mobile bookmarking application written in Go using Gio UI framework.

## Features
- Native Android/iOS support
- Local storage with SQLite
- Optional cloud sync (Google Drive/OneDrive)
- Share integration for bookmarking links and images
- Tag management
- Customizable sync settings

## Project Structure
```
goBookMarker/
├── cmd/
│   └── mobile/          # Mobile app entry point
├── internal/
│   ├── app/            # Application core
│   ├── auth/           # Authentication services
│   ├── storage/        # Local storage implementation
│   ├── sync/           # Cloud sync services
│   ├── ui/             # UI components
│   └── models/         # Data models
└── pkg/
    ├── share/          # Share handling
    └── utils/          # Utility functions
```

## Development Setup
1. Install Go 1.21 or later
2. Install Gio dependencies
3. Install gomobile
4. Run `go mod tidy`
5. Build for Android: `gomobile build -target=android`
6. Build for iOS: `gomobile build -target=ios`
