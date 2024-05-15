package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i LoginRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i AccessRepository -o ./mocks/ -s "_minimock.go"
