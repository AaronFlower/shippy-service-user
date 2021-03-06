package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/aaronflower/shippy-service-user/proto/auth"
	micro "github.com/micro/go-micro"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo         Repository
	tokenServcie Authable
	Publisher    micro.Publisher
}

func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(user)
	log.Println(user.Password, req.Password)

	// Compare our given password against the hashed password stored in the database.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return err
	}

	token, err := s.tokenServcie.Encode(user)
	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	log.Println("Creating user: ", req)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Password = string(hashedPass)
	if err := s.repo.Create(req); err != nil {
		log.Println(err)
		return err
	}
	res.User = req

	// Let's publish the event once we create a new user.
	if err := s.Publisher.Publish(ctx, req); err != nil {
		log.Println("Publish failed:", err)
		return err
	}
	return nil
}

func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	// Decode token
	claims, err := s.tokenServcie.Decode(req.Token)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}
	res.Valid = true
	return nil
}
