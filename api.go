package main

import (
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	rePhoneNumber = regexp.MustCompile(`^1\d{10}$`)
)

func StartServer(addr string) error {
	r := gin.Default()

	r.GET("test", test)
	r.GET("records", QuerySendDetailRaw)
	r.GET("sendrecords", QuerySendRecord)
	r.POST("sms", SendSMS)

	r.Run(addr)
	return nil
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "world",
	})
	return
}

func setErrorCode(c *gin.Context, code int, message string) {
	c.JSON(200, gin.H{
		"Code":    code,
		"Message": message,
	})
}

func SendSMS(c *gin.Context) {
	ip := c.ClientIP()
	if !isIpAllow(ip) {
		setErrorCode(c, 40000, "当前IP("+ip+")不在白名单内")
		return
	}
	var data struct {
		PhoneNumber string
		User        string
		DeviceId    string
		State       string
	}
	if err := c.BindJSON(&data); err != nil {
		setErrorCode(c, 40001, "参数错误")
		return
	}
	if !rePhoneNumber.MatchString(data.PhoneNumber) {
		setErrorCode(c, 40002, "手机号码格式错误")
		return
	}
	if !isPhoneNumberAllow(data.PhoneNumber) {
		setErrorCode(c, 40002, "目标手机号码("+data.PhoneNumber+")不在白名单内")
		return
	}
	if data.State == "" {
		setErrorCode(c, 40003, "状态不能为空")
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		print(err)
		setErrorCode(c, 50001, "内部错误")
		return
	}
	if err := sendDeviceState(data.PhoneNumber, data.User, data.DeviceId,
		data.State, id.String()); err != nil {
		print(err)
		setErrorCode(c, 50002, "内部错误")
		return
	}
	setErrorCode(c, 20000, "发送成功")
}

func QuerySendDetailRaw(c *gin.Context) {
	phoneNumber := c.Query("phone")
	if phoneNumber == "" {
		setErrorCode(c, 40020, "必须指定要查询的手机号码")
		return
	}
	if !rePhoneNumber.MatchString(phoneNumber) {
		setErrorCode(c, 40002, "手机号码格式错误")
		return
	}
	result := querySmsDetail(phoneNumber)
	c.JSON(200, gin.H{
		"Code": 20000,
		"Data": result,
	})
}

func QuerySendRecord(c *gin.Context) {
	query := bson.D{}
	phoneNumber := c.Query("phone")
	if phoneNumber != "" {
		query = append(query, bson.E{"phonenumber", phoneNumber})
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pagesize", "20"))
	if err != nil || pageSize <= 0 || pageSize > 500 {
		setErrorCode(c, 40010, "pagesize参数不合法")
		return
	}
	pageIndex, err := strconv.Atoi(c.DefaultQuery("pageindex", "0"))
	if err != nil || pageIndex < 0 || pageIndex > 1000000 {
		setErrorCode(c, 40010, "pageindex参数不合法")
		return
	}
	if pageIndex*pageSize > 1000000 {
		setErrorCode(c, 40010, "超出允许的查询范围")
		return
	}
	opt := options.Find().SetLimit(int64(pageSize))
	if pageIndex != 0 {
		opt.SetSkip(int64(pageSize * pageIndex))
	}
	cursor, err := colSendRecord.Find(c, query, opt)
	if err != nil {
		dump(err)
		setErrorCode(c, 50003, "读取数据库错误")
		return
	}
	i := 0
	var records []SendRecord
	for cursor.Next(c) {
		i++
		var record SendRecord
		err := cursor.Decode(&record)
		if err != nil {
			dump(err)
			setErrorCode(c, 50004, "解析数据库返回值错误")
			return
		}
		records = append(records, record)
	}
	c.JSON(200, gin.H{
		"Code": 20000,
		"Data": records,
	})
}
