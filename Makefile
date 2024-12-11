.DEFAULT_GOAL := dev

tw:
	@npx tailwindcss -i input.css -o public/static/css/tw.css --watch

dev:
	go run cmd/server/main.go
