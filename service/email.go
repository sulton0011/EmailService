package service

import (
	"context"
	"fmt"
	"log"

	gpb "github.com/golang/protobuf/ptypes/empty"
	l "github.com/home_work/TaskUserService/EmailService/pkg/logger"
	"github.com/home_work/TaskUserService/EmailService/config"
	pb "github.com/home_work/TaskUserService/EmailService/genproto/email_service"
	"github.com/home_work/TaskUserService/EmailService/storage"
	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	gomail "gopkg.in/gomail.v2"
)

// SendService ...
type SendService struct {
	storage storage.I
	Conf config.Config
	DB   *sqlx.DB
	logger l.Logger
}

// NewSendService ...
func NewSendService(db *sqlx.DB, conf config.Config) *SendService {
	return &SendService{storage: storage.NewStoragePg(db), Conf: conf}
}

//Send ...
func (s *SendService) Send(ctx context.Context, req *pb.Email) (*gpb.Empty, error) {
	id, err := s.storage.SendEmail().CreatEmailText(
		uuid.New().String(), req.Subject, req.Body,
	)
	if err != nil {
		s.logger.Error("Error creat emailText", l.Error(err))
		return nil, err
	}
	
	for _, val := range req.Recipients {
		status := true

		if err := s.sendEmail(req.Subject, req.Body, val); err != nil {
			s.logger.Error("No message was sent to" + val, l.Error(err))
			status = false
		}

		err = s.storage.SendEmail().CreatEmail(
			uuid.New().String(),
			id,
			val,
			status,
		)
		if err != nil {
			s.logger.Error("Error creat posgresql" + val, l.Error(err))

		}
	}

	return &gpb.Empty{}, nil
}

func (s *SendService) sendEmail(subject, body, email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Conf.EmailFromHeader)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email to
	d := gomail.NewDialer(s.Conf.SMTPHost, s.Conf.SMTPPort, s.Conf.SMTPUser, s.Conf.SMTPUserPass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	log.Print("Sent")
	fmt.Println("Sent")
	return nil
}