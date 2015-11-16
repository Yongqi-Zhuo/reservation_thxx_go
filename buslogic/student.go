package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"github.com/shudiwsh2009/reservation_thxx_go/domain"
	"github.com/shudiwsh2009/reservation_thxx_go/sms"
	"github.com/shudiwsh2009/reservation_thxx_go/util"
	"strings"
	"time"
)

type StudentLogic struct {
}

// 学生预约咨询
func (sl *StudentLogic) MakeReservationByStudent(reservationId string, name string, gender string,
	studentId string, school string, hometown string, mobile string, email string, experience string,
	problem string) (*domain.Reservation, error) {
	if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(name, "") {
		return nil, errors.New("姓名为空")
	} else if strings.EqualFold(gender, "") {
		return nil, errors.New("性别为空")
	} else if strings.EqualFold(studentId, "") {
		return nil, errors.New("学号为空")
	} else if strings.EqualFold(school, "") {
		return nil, errors.New("院系为空")
	} else if strings.EqualFold(hometown, "") {
		return nil, errors.New("生源地为空")
	} else if strings.EqualFold(mobile, "") {
		return nil, errors.New("手机号为空")
	} else if strings.EqualFold(email, "") {
		return nil, errors.New("邮箱为空")
	} else if strings.EqualFold(experience, "") {
		return nil, errors.New("咨询经历为空")
	} else if strings.EqualFold(problem, "") {
		return nil, errors.New("咨询问题为空")
	} else if !util.IsStudentId(studentId) {
		return nil, errors.New("学号不正确")
	} else if !util.IsMobile(mobile) {
		return nil, errors.New("手机号格式不正确")
	} else if !util.IsEmail(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.Before(time.Now().Local()) {
		return nil, errors.New("咨询已过期")
	} else if reservation.Status != domain.AVAILABLE {
		return nil, errors.New("咨询已被预约")
	}
	studentReservations, err := data.GetReservationsByStudentId(studentId)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	for _, r := range studentReservations {
		if r.Status == domain.RESERVATED && r.StartTime.After(time.Now().Local()) {
			return nil, errors.New("你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62792453取消已有预约。")
		}
	}
	reservation.StudentInfo = domain.StudentInfo{
		Name:       name,
		Gender:     gender,
		StudentId:  studentId,
		School:     school,
		Hometown:   hometown,
		Mobile:     mobile,
		Email:      email,
		Experience: experience,
		Problem:    problem,
	}
	reservation.Status = domain.RESERVATED
	err = data.UpsertReservation(reservation)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}

	// send success sms
	if checkReservation, err := data.GetReservationById(reservationId); err == nil &&
		checkReservation.Status == domain.RESERVATED && strings.EqualFold(checkReservation.StudentInfo.Mobile, mobile) {
		sms.SendSuccessSMS(checkReservation)
	}
	return reservation, nil
}

// 学生拉取反馈
func (sl *StudentLogic) GetFeedbackByStudent(reservationId string, studentId string) (*domain.Reservation, error) {
	if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(studentId, "") || !util.IsStudentId(studentId) {
		return nil, errors.New("学号不正确")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == domain.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.StudentInfo.StudentId, studentId) {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	return reservation, nil
}

// 学生反馈
func (sl *StudentLogic) SubmitFeedbackByStudent(reservationId string, name string, problem string, choices string,
	score string, feedback string, studentId string) (*domain.Reservation, error) {
	if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(name, "") {
		return nil, errors.New("姓名为空")
	} else if strings.EqualFold(problem, "") {
		return nil, errors.New("咨询问题为空")
	} else if strings.EqualFold(choices, "") {
		return nil, errors.New("选项为空")
	} else if strings.EqualFold(score, "") {
		return nil, errors.New("总评为空")
	} else if strings.EqualFold(feedback, "") {
		return nil, errors.New("反馈为空")
	} else if strings.EqualFold(studentId, "") || !util.IsStudentId(studentId) {
		return nil, errors.New("学号不正确")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == domain.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.StudentInfo.StudentId, studentId) {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	reservation.StudentFeedback = domain.StudentFeedback{
		Name:     name,
		Problem:  problem,
		Choices:  choices,
		Score:    score,
		Feedback: feedback,
	}
	if err = data.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}
