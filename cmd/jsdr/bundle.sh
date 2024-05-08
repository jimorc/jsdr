cd $GOPATH/jsdr/internal/jsdrgui
fyne bundle -o bundled.go $GOPATH/jsdr/cmd/jsdr/images/black-wrench.svg
fyne bundle -o bundled.go -append $GOPATH/jsdr/cmd/jsdr/images/start.svg
fyne bundle -o bundled.go -append $GOPATH/jsdr/cmd/jsdr/images/stop.svg
fyne bundle -o bundled.go -append $GOPATH/jsdr/cmd/jsdr/images/logsettings.svg
fyne bundle -o bundled.go -append $GOPATH/jsdr/cmd/jsdr/images/radio.svg

if [[ "$OSTYPE" == "darwin"* ]]; then
sed -i '' 's/main/jsdrgui/' bundled.go
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
sed 's/main/jsdrgui/' bundled.go
fi
cd $GOPATH/jsdr/cmd/jsdr