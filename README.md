# go-backend (Core)

## ภาพรวม
Repo นี้เป็น service หลักของระบบ MooMoon ฝั่ง Go โดยรวบรวม Logic หลัก (Services) และ Routing ย่อย (Endpoints) เพื่อให้ทำงานร่วมกับ AI service, ฐานข้อมูล, LINE/Telegram, etc.

## โครงสร้างโฟลเดอร์

- **config/**  
  เก็บไฟล์ตั้งค่า เช่น `config.go` หรือไฟล์ตัวอย่าง `.env.example`  

- **internal/services/**  
  - `user_service.go`: ดูแลการยืนยันตัวตน (auth), ตั้งค่า role, แก้โปรไฟล์, 2FA  
  - `coin_service.go`: ตรวจสอบ, เติม, หัก, โอนเหรียญ  
  - `package_service.go`: จัดการแพ็กเกจสมาชิก (Membership), สิทธิ์ต่างๆ  
  - `payment_service.go`: สร้าง Payment, Verify การจ่ายเงิน, คำนวณค่าคอมมิชชั่น  
  - `notification_service.go`: ส่งข้อความแจ้งเตือนผ่าน LINE/Telegram  
  - `ai_router_service.go`: เลือก AI model ตาม config แล้ว forward ไปยัง ai-service  
  - `workpool_service.go`: จัดการงานแบ็คกราวด์ (work queue), จับเวลากำหนดเสร็จ  
  - `review_service.go`: รับ/อนุมัติ/ลบ รีวิว, อนุญาตให้ผู้ใช้ยื่นอุทธรณ์  
  - `rank_service.go`: คำนวณ Ranking, ค่าคอมมิชชั่น, โบนัส

- **internal/repositories/**  
  ถ้าต้องแยก layer สำหรับเข้าถึงฐานข้อมูล (เช่น Postgres, Firestore) ให้วาง Struct และ Methods ที่สัมพันธ์กับ DB ที่นี่

- **internal/models/**  
  Struct หรือ DTO ที่ใช้สื่อสารระหว่าง service, handler, database

- **internal/utils/**  
  ฟังก์ชันช่วยเหลือ (เช่น Logger, Validator, Helper functions)

- **routes/**  
  - `user.go`: จัดการ login, auth, 2FA  
  - `coin.go`: เช็คเหรียญ, เติม, โอน  
  - `package.go`: ซื้อแพ็กเกจ, ตรวจสิทธิ์  
  - `booking.go`: จองคิว, เลือก slot, แจ้งเตือน  
  - `ai.go`: ส่งคำถามไป AI, ตีความไพ่  
  - `notification.go`: Broadcast แจ้งเตือน LINE  
  - `review.go`: จัดการ รีวิว (ให้, อนุมัติ, ลบ)  
  - `rank.go`: จัดอันดับ, คำนวณค่าคอมมิชชั่น

- **scripts/**  
  สคริปต์ช่วยเหลือเช่น  
  - `migrate.sh`: สร้างตาราง, seed data เบื้องต้น  
  - `build.sh`: build binary หรือ docker image

- **Dockerfile / docker-compose.yml**  
  กรณีต้องการรันพร้อม DB หรือ service อื่นในเครื่องเดียวกัน

## วิธีใช้งานเบื้องต้น

1. ตั้งค่า environment variables (ดูตัวอย่างใน `config/.env.example`)  
2. ติดตั้ง dependencies  
   ```bash
   go mod download
   ```  
3. รัน migrations (ถ้ามี)  
   ```bash
   ./scripts/migrate.sh
   ```  
4. รันเซอร์วิส  
   ```bash
   go run cmd/main.go
   ```  
   หรือ build แล้วรัน  
   ```bash
   go build -o moo-backend ./cmd
   ./moo-backend
   ```  
5. เลือกอ่านเพิ่มเติมในแต่ละไฟล์ service / route  

## Environment Variables (ตัวอย่าง)  
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=...
DB_PASS=...
DB_NAME=...
LINE_CHANNEL_TOKEN=...
TELEGRAM_BOT_TOKEN=...
AI_ROUTER_URL=http://localhost:8000
...
```