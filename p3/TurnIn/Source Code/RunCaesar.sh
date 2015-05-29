if dpkg-query -W golang-go; then
    echo "Has Golang!"
else
    sudo apt-get install golang-go
fi
rm cutfile.txt
clear
go run Caesar.go $1

