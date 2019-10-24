package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"go.mongodb.org/mongo-driver/bson"
)

type sendSmsResponse struct {
	RequestId string `json:"RequestId" xml:"RequestId"`
	BizId     string `json:"BizId" xml:"BizId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

type smsSendDetail struct {
	PhoneNum    string `json:"PhoneNum" xml:"PhoneNum"`
	SendStatus  string `json:"SendStatus" xml:"SendStatus"`
	ErrCode     string `json:"ErrCode" xml:"ErrCode"`
	Content     string `json:"Content" xml:"Content"`
	SendDate    string `json:"SendDate" xml:"SendDate"`
	ReceiveDate string `json:"ReceiveDate" xml:"ReceiveDate"`
	ExtraInfo   string `json:"ExtraInfo" xml:"ExtraInfo"`
}

// extraInfo 不会包含在短信内容中，但是会记录在发送记录里面，通过querySmsDetail
// 接口可查询到这些额外的信息
func sendDeviceState(phoneNumber string, user string, deviceId string, state string, extraInfo string) error {
	if len(phoneNumber) != 11 || phoneNumber[0] != '1' {
		return errors.New("phoneNum must be 1xxxxxxxxxx")
	}

	data, err := json.Marshal(map[string]string{
		"user":    user,
		"devid":   deviceId,
		"message": state,
	})
	if err != nil {
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Domain = globalConfig.SmsDomain
	request.SignName = globalConfig.SmsSignName
	request.PhoneNumbers = phoneNumber
	request.TemplateCode = globalConfig.SmsTemplateCode
	request.TemplateParam = string(data)
	request.OutId = extraInfo

	client, err := dysmsapi.NewClientWithAccessKey(globalConfig.SmsZone,
		globalConfig.AliyunAccessKey, globalConfig.AliyunAccessSecret)
	if err != nil {
		return err
	}
	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	//   response := sendSmsResponse{
	//     RequestId: "dry-run",
	//     BizId:     "dry-run",
	//     Code:      "OK",
	//     Message:   "dry-run",
	//   }

	if response.Code != "OK" {
		dump(response)
		return errors.New(response.Code + ": " + response.Message)
	}

	record := SendRecord{
		Domain:          globalConfig.SmsDomain,
		SignName:        globalConfig.SmsSignName,
		PhoneNumber:     phoneNumber,
		TemplateCode:    globalConfig.SmsTemplateCode,
		TemplateParam:   string(data),
		ExtraInfo:       extraInfo,
		RequestId:       response.RequestId,
		BizId:           response.BizId,
		ResponseCode:    response.Code,
		ResponseMessage: response.Message,
		CreatedAt:       now(),
	}
	res, err := colSendRecord.InsertOne(ctxEmpty, record)
	if err != nil {
		return err
	}
	dump(res)

	return nil
}

func parseStatusCode(code int64) string {
	switch code {
	case 1:
		return "等待回执"
	case 2:
		return "发送失败"
	case 3:
		return "发送成功"
	default:
		return "未知状态"
	}
}

func querySmsDetailByDate(
	client *dysmsapi.Client, phoneNumber string, sendDate string) (
	result []smsSendDetail, e error) {
	pageSize := requests.NewInteger(5)
	currentPage := 1
	totalCount := 0

	for totalCount == 0 || len(result) < totalCount {
		request := dysmsapi.CreateQuerySendDetailsRequest()
		request.Domain = globalConfig.SmsDomain
		request.PhoneNumber = phoneNumber
		request.SendDate = sendDate
		request.PageSize = pageSize
		request.CurrentPage = requests.NewInteger(currentPage)

		response, err := client.QuerySendDetails(request)
		if err != nil {
			return result, err
		}
		if totalCount == 0 {
			totalCount, err = strconv.Atoi(response.TotalCount)
			if err != nil {
				return result, err
			}
			if totalCount == 0 {
				return result, nil
			}
		}
		for i := range response.SmsSendDetailDTOs.SmsSendDetailDTO {
			detail := &response.SmsSendDetailDTOs.SmsSendDetailDTO[i]
			d := smsSendDetail{
				PhoneNum:    detail.PhoneNum,
				SendStatus:  parseStatusCode(detail.SendStatus),
				ErrCode:     detail.ErrCode,
				Content:     detail.Content,
				SendDate:    detail.SendDate,
				ReceiveDate: detail.ReceiveDate,
				ExtraInfo:   detail.OutId,
			}
			result = append(result, d)
		}
		currentPage++
	}

	return result, nil
}

func findRecordByExtraInfo(records []SendRecord, info string) *SendRecord {
	for n := range records {
		if records[n].ExtraInfo == info {
			return &records[n]
		}
	}
	return nil
}

func updateSendRecordsFromAliyun(records []SendRecord) error {
	needUpdate := false
	now := time.Now()
	// 整理出需要查询的记录, 按手机号码和发送日期归类
	all := make(map[string]map[string]bool)
	for n := range records {
		rec := &records[n]
		if rec.SendStatus == parseStatusCode(2) || // 已更新过状态, 发送失败
			rec.SendStatus == parseStatusCode(3) || // 已更新过状态, 发送成功
			now.Sub(rec.CreatedAt).Hours() > 35*24 { // 阿里云只能查30天的记录, 这里放宽到35天
			continue
		}
		if all[rec.PhoneNumber] == nil {
			all[rec.PhoneNumber] = make(map[string]bool)
		}
		all[rec.PhoneNumber][rec.CreatedAt.Format("20060102")] = true
		needUpdate = true
	}
	if !needUpdate {
		return nil
	}
	// 连接阿里云服务器
	client, err := dysmsapi.NewClientWithAccessKey(globalConfig.SmsZone,
		globalConfig.AliyunAccessKey, globalConfig.AliyunAccessSecret)
	if err != nil {
		dump(err)
		return err
	}
	// 开始查询, 遍历每个手机号码的每个日期
	for phoneNum, dateList := range all {
		for s := range dateList {
			// 这里先查出这个手机号码在这一天的所有短信发送记录
			data, err := querySmsDetailByDate(client, phoneNum, s)
			if err != nil {
				dump(err)
				continue
			}
			// 然后对于每条记录, 根据 ExtraInfo 再回去 records 里面找,
			// 如果有对应的, 则更新该记录, 如果没对应的, 则忽略
			for n := range data {
				record := findRecordByExtraInfo(records, data[n].ExtraInfo)
				if record != nil {
					record.SendStatus = data[n].SendStatus
					record.ErrCode = data[n].ErrCode
					record.Content = data[n].Content
					// 阿里云存储的是UTC+8的时间, 但是返回的字符串没有带时区
					// 所以这里需要按UTC+8来解析, 结果才正确
					tz := time.FixedZone("UTC+8", 8*60*60)
					record.SendDate, _ = time.ParseInLocation(
						"2006-01-02 15:04:05", data[n].SendDate, tz)
					record.ReceiveDate, _ = time.ParseInLocation(
						"2006-01-02 15:04:05", data[n].ReceiveDate, tz)
					// 对于找到的每条记录, 更新数据库中对应的记录, 这里本来应该
					// 用 UpdateOne, 不过用 ReplaceOne 可以省得列出每个修改项,
					// 所以就偷懒了
					_, err := colSendRecord.ReplaceOne(ctxEmpty,
						bson.M{"_id": record.Id}, record)
					if err != nil {
						dump(err)
					}
				}
			}
		}
	}
	return nil
}

func querySmsDetail(phoneNumber string) (result []smsSendDetail) {
	client, err := dysmsapi.NewClientWithAccessKey(globalConfig.SmsZone,
		globalConfig.AliyunAccessKey, globalConfig.AliyunAccessSecret)
	if err != nil {
		dump(err)
		return nil
	}

	for i := 0; i < 35; i++ {
		sendDate := time.Now().AddDate(0, 0, -i).Format("20060102")
		data, err := querySmsDetailByDate(client, phoneNumber, sendDate)
		result = append(result, data...)
		if err != nil {
			dump(err)
			return result
		}
	}
	return result
}
