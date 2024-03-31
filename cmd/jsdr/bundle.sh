cd $GOPATH/jsdr/internal/gosdrgui
fyne bundle -o bundled.go $GOPATH/jsdr/cmd/go_sdr/images/start.svg
fyne bundle -o bundled.go -append $GOPATH/jsdr/cmd/go_sdr/images/stop.svg
if [[ "$OSTYPE" == "darwin"* ]]; then
sed -i '' 's/main/jsdrgui/' bundled.go
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
sed 's/main/jsdrgui/' bundled.go
fi
cd $GOPATH/jsdr/cmd/jsdr