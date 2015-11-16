package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"strings"
	"time"
)

type AdminLogic struct {
}

// 管理员添加咨询
func (al *AdminLogic) AddReservationByAdmin(startTime string, endTime string, teacherUsername string,
	teacherFullname string, teacherMobile string, username string, userType models.UserType) (*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(startTime, "") {
		return nil, errors.New("开始时间为空")
	} else if strings.EqualFold(endTime, "") {
		return nil, errors.New("结束时间为空")
	} else if strings.EqualFold(teacherUsername, "") {
		return nil, errors.New("咨询师工号为空")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherMobile, "") {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	start, err := time.Parse(utils.TIME_PATTERN, startTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.Parse(utils.TIME_PATTERN, endTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := models.GetUserByUsername(teacherUsername)
	if err != nil {
		if teacher, err = models.AddFullUser(teacherUsername, TeacherDefaultPassword, teacherFullname,
			teacherMobile, models.TEACHER); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertUser(teacher); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation, err := models.AddReservation(start, end, teacher.Fullname, teacher.Username, teacher.Mobile)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员编辑咨询
func (al *AdminLogic) EditReservationByAdmin(reservationId string, startTime string, endTime string,
	teacherUsername string, teacherFullname string, teacherMobile string, username string,
	userType models.UserType) (*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(startTime, "") {
		return nil, errors.New("开始时间为空")
	} else if strings.EqualFold(endTime, "") {
		return nil, errors.New("结束时间为空")
	} else if strings.EqualFold(teacherUsername, "") {
		return nil, errors.New("咨询师工号为空")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherMobile, "") {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == models.RESERVATED {
		return nil, errors.New("不能编辑已被预约的咨询")
	}
	start, err := time.Parse(utils.TIME_PATTERN, startTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.Parse(utils.TIME_PATTERN, endTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := models.GetUserByUsername(teacherUsername)
	if err != nil {
		if teacher, err = models.AddFullUser(teacherUsername, TeacherDefaultPassword, teacherFullname,
			teacherMobile, models.TEACHER); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertUser(teacher); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherUsername = teacher.Username
	reservation.TeacherFullname = teacher.Fullname
	reservation.TeacherMobile = teacher.Mobile
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员删除咨询
func (al *AdminLogic) RemoveReservationsByAdmin(reservationIds []string, username string, userType models.UserType) error {
	if strings.EqualFold(username, "") {
		return errors.New("请先登录")
	} else if userType != models.TEACHER {
		return errors.New("权限不足")
	} else if reservationIds == nil {
		return errors.New("咨询Id列表为空")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return errors.New("管理员账户出错,请联系技术支持")
	}
	for _, reservationId := range reservationIds {
		if reservation, err := models.GetReservationById(reservationId); err == nil {
			reservation.Status = models.DELETED
			models.UpsertReservation(reservation)
		}
	}
	return nil
}

// 管理员取消预约
func (al *AdminLogic) CancelReservationsByAdmin(reservationIds []string, username string, userType models.UserType) error {
	if strings.EqualFold(username, "") {
		return errors.New("请先登录")
	} else if userType != models.TEACHER {
		return errors.New("权限不足")
	} else if reservationIds == nil {
		return errors.New("咨询Id列表为空")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return errors.New("管理员账户出错,请联系技术支持")
	}
	for _, reservationId := range reservationIds {
		reseravtion, err := models.GetReservationById(reservationId)
		if err != nil || reseravtion.Status == models.DELETED {
			continue
		}
		if reseravtion.Status == models.RESERVATED && reseravtion.StartTime.After(time.Now().Local()) {
			reseravtion.Status = models.AVAILABLE
			reseravtion.StudentInfo = models.StudentInfo{}
			reseravtion.StudentFeedback = models.StudentFeedback{}
			reseravtion.TeacherFeedback = models.TeacherFeedback{}
			models.UpsertReservation(reseravtion)
		}
	}
	return nil
}

// 管理员拉取反馈
func (al *AdminLogic) GetFeedbackByAdmin(reservationId string, username string, userType models.UserType) (*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	}
	return reservation, nil
}

// 管理员提交反馈
func (al *AdminLogic) SubmitFeedbackByAdmin(reservationId string, teacherFullname string, teacherId string,
	studentName string, problem string, solution string, adviceToCenter string, username string,
	userType models.UserType) (*models.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherId, "") {
		return nil, errors.New("咨询师工作证号为空")
	} else if strings.EqualFold(studentName, "") {
		return nil, errors.New("学生姓名为空")
	} else if strings.EqualFold(problem, "") {
		return nil, errors.New("咨询问题为空")
	} else if strings.EqualFold(solution, "") {
		return nil, errors.New("解决方法为空")
	} else if strings.EqualFold(adviceToCenter, "") {
		return nil, errors.New("工作建议为空")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(teacherId, reservation.TeacherUsername) {
		return nil, errors.New("咨询师工号不匹配")
	}
	if reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty() {
		utils.SendFeedbackSMS(reservation)
	}
	reservation.TeacherFeedback = models.TeacherFeedback{
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherId,
		StudentFullname: studentName,
		Problem:         problem,
		Solution:        solution,
		AdviceToCenter:  adviceToCenter,
	}
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员查看学生信息
func (al *AdminLogic) GetStudentInfoByAdmin(reservationId string, username string, userType models.UserType) (*models.StudentInfo, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,无法查看")
	}
	return &reservation.StudentInfo, nil
}

// 管理员导出咨询
func (al *AdminLogic) ExportReservationsByAdmin(reservationIds []string, username string, userType models.UserType) (string, error) {
	if strings.EqualFold(username, "") {
		return "", errors.New("请先登录")
	} else if userType != models.TEACHER {
		return "", errors.New("权限不足")
	} else if reservationIds == nil {
		return "", errors.New("咨询Id列表为空")
	}
	admin, err := models.GetUserByUsername(username)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	var reservations []*models.Reservation
	for _, reservationId := range reservationIds {
		reservation, err := models.GetReservationById(reservationId)
		if err != nil {
			continue
		}
		reservations = append(reservations, reservation)
	}
	filename := "export_" + time.Now().Local().Format(utils.TIME_PATTERN) + utils.ExcelSuffix
	if len(reservations) == 0 {
		return "", nil
	}
	if err = utils.ExportReservationsToExcel(reservations, filename); err != nil {
		return "", err
	}
	return utils.ExportPrefix + filename, nil
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (ul *UserLogic) SearchTeacher(fullname string, username string, mobile string, admin string, userType models.UserType) (*models.User, error) {
	if strings.EqualFold(admin, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	user, err := models.GetUserByUsername(admin)
	if err != nil || user.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	if !strings.EqualFold(fullname, "") {
		user, err := models.GetUserByFullname(fullname)
		if err == nil {
			return user, nil
		}
	}
	if !strings.EqualFold(username, "") {
		user, err := models.GetUserByUsername(username)
		if err == nil {
			return user, nil
		}
	}
	if !strings.EqualFold(mobile, "") {
		user, err := models.GetUserByMobile(mobile)
		if err == nil {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}
