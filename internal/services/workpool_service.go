package services

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/poomiiz/go-backend/internal/utils"
)

// Job โครงสร้างข้อมูลงานใน Firestore
type Job struct {
	Name        string    `firestore:"name"`
	Payload     string    `firestore:"payload"`     // เก็บ JSON encoded หรือข้อความที่ dispatch
	ScheduledAt time.Time `firestore:"scheduledAt"` // เวลาเรียกทำ
	Status      string    `firestore:"status"`      // "pending", "processing", "done", "failed"
	CreatedAt   time.Time `firestore:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt"`
}

type WorkpoolService struct {
	col *firestore.CollectionRef
}

func NewWorkpoolService() *WorkpoolService {
	return &WorkpoolService{
		col: utils.Client.Collection("jobs"),
	}
}

// ScheduleJob: ใส่ job เข้า Firestore
func (s *WorkpoolService) ScheduleJob(ctx context.Context, name, payload string, runAt time.Time) (string, error) {
	now := time.Now()
	job := Job{
		Name:        name,
		Payload:     payload,
		ScheduledAt: runAt,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	docRef, _, err := s.col.Add(ctx, job)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// ProcessDueJobs: ดึง job ที่ scheduledAt <= now และ status=="pending" มา execute
func (s *WorkpoolService) ProcessDueJobs(ctx context.Context) error {
	now := time.Now()
	q := s.col.Where("scheduledAt", "<=", now).Where("status", "==", "pending")
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	for _, doc := range docs {
		var job Job
		if err := doc.DataTo(&job); err != nil {
			continue // อ่านไม่ออกข้ามไป
		}

		// mark ว่ากำลังประมวลผล
		_, _ = doc.Ref.Update(ctx, []firestore.Update{
			{Path: "status", Value: "processing"},
			{Path: "updatedAt", Value: time.Now()},
		})

		// เรียกใช้ executeJob เพื่อ dispatch จริง
		errExec := s.executeJob(ctx, doc.Ref.ID, job.Name, job.Payload)
		newStatus := "done"
		if errExec != nil {
			newStatus = "failed"
		}

		// อัปเดตสถานะสุดท้าย
		_, _ = doc.Ref.Update(ctx, []firestore.Update{
			{Path: "status", Value: newStatus},
			{Path: "updatedAt", Value: time.Now()},
		})
	}

	return nil
}

// executeJob: ใส่ logic เรียก Service อื่นตามชื่อ job.Name
func (s *WorkpoolService) executeJob(ctx context.Context, jobID, name, payload string) error {
	// TODO: หมายเหตุ: ต้องเขียน dispatch logic ตามชื่อ job เช่น
	// if name == "send_reminder" {
	//   // แปลง payload เป็น struct ที่ต้องการ แล้วเรียก NotificationService.SendLineMessage/SendTelegramAlert
	// }
	// สำหรับตอนนี้ return nil (dummy)
	return nil
}
