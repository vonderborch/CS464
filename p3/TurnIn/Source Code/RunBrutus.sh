if dpkg-query -W golang-go; then
    echo "Has Golang!"
else
    sudo apt-get install golang-go
fi
clear
go run Officer.go Brutus

