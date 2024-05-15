package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/semho/chat-microservices/auth/internal/model"
	descAccess "github.com/semho/chat-microservices/auth/pkg/access_v1"
	desc "github.com/semho/chat-microservices/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"os"
	"path/filepath"
)

const (
	address = "localhost:50051"
	userID  = 1
)

func getUserByID(ctx context.Context, client desc.AuthV1Client, userID int64) (*desc.UserResponse, error) {
	response, err := client.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func createUser(ctx context.Context, client desc.AuthV1Client) (*desc.CreateResponse, error) {
	userDetail := &desc.UserDetail{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  desc.Role_user.Enum(),
	}
	pass := gofakeit.Password(true, false, false, false, false, 32)
	passwordDetail := &desc.UserPassword{
		Password:        pass,
		PasswordConfirm: pass,
	}
	request := &desc.CreateRequest{
		Detail:   userDetail,
		Password: passwordDetail,
	}

	response, err := client.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func updateUserByID(ctx context.Context, client desc.AuthV1Client, userID int64) (*emptypb.Empty, error) {
	userUpdate := &desc.UpdateUserInfo{
		Name:  &wrapperspb.StringValue{Value: gofakeit.Name()},
		Email: &wrapperspb.StringValue{Value: gofakeit.Email()},
	}

	request := &desc.UpdateRequest{Id: userID, Info: userUpdate}
	response, err := client.Update(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func deleteUserByID(ctx context.Context, client desc.AuthV1Client, userID int64) (*emptypb.Empty, error) {
	request := &desc.DeleteRequest{Id: userID}
	response, err := client.Delete(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func credsService() credentials.TransportCredentials {
	//Получаем текущую директорию
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Ошибка при получении текущей директории:", err)
	}

	certFile := filepath.Join(dir, "tls", "service.pem")

	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	return creds
}

var accessToken = flag.String("a", "", "Access token")

func main() {
	//TODO: для access
	flag.Parse()
	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	//TODO: для TLS
	//conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credsService()))
	//без tls
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	//TODO: для auth
	//c := desc.NewAuthV1Client(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//
	//r, err := getUserByID(ctx, c, userID)

	//r, err := createUser(ctx, c)
	//r, err := updateUserByID(ctx, c, userID)
	//r, err := deleteUserByID(ctx, c, userID)

	//if err != nil {
	//	log.Fatalf("failed to get user by id: %v", err)
	//}

	//log.Printf(color.RedString("Answer: \n"), color.GreenString("%+v", r))

	//TODO: для access
	cl := descAccess.NewAccessV1Client(conn)

	_, err = cl.Check(
		ctx, &descAccess.CheckRequest{
			EndpointAddress: model.PathUserCreate,
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Access granted")

}
