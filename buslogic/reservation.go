package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"strings"
	"time"
)

type ReservationLogic struct {
}

// 学生查看前后一周内的所有咨询
func (rl *ReservationLogic) GetReservationsByStudent() ([]*models.Reservation, error) {
	from := time.Now().Local().AddDate(0, 0, 7)
	to := time.Now().Local().AddDate(0, 0, -7)
	reservations, err := models.GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(time.Now().Local()) {
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

// 咨询师查看负7天之后的所有咨询
func (rl *ReservationLogic) GetReservationsByTeacher(username string, userType models.UserType) ([]*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	teacher, err := models.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	from := time.Now().Local().AddDate(0, 0, -7)
	reservations, err := models.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(time.Now().Local()) {
			continue
		} else if strings.EqualFold(r.TeacherUsername, teacher.Username) {
			result = append(result, r)
		}
	}
	return result, nil
}

// 管理员查看负7天之后的所有咨询
func (rl *ReservationLogic) GetReservationsByAdmin(username string, userType models.UserType) ([]*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from := time.Now().Local().AddDate(0, 0, -7)
	reservations, err := models.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(time.Now().Local()) {
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

// 管理员查看指定日期后30天内的所有咨询
func (rl *ReservationLogic) GetReservationsMonthlyByAdmin(username string, userType string) ([]*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	from := time.Now().Local().AddDate(0, 0, -30)
	reservations, err := models.GetReservationsAfterTime(from)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	var result []*models.Reservation
	for _, r := range reservations {
		if r.Status == models.AVAILABLE && r.StartTime.Before(time.Now().Local()) {
			continue
		}
		result = append(result, r)
	}
	return result, nil
}
