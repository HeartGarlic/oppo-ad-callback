package oppo_ad_callback

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Oppo 广告回传
const (
	sendDataURL = "https://sapi.ads.oppomobile.com/v1/clue/sendData"
)

// TransFormTypeSubmit 转化类型编码（transformType）
const (
	TransFormTypeSubmit    = 1   // 表单提交（或线索提交）
	TransFormTypeCall      = 3   // 电话拨打
	TransFormTypeChat      = 4   // 在线咨询
	TransFormTypeBuy       = 5   // 商品购买
	TransFormTypeVisit     = 6   // 页面访问
	TransFormTypeCopy      = 7   // 微信复制
	TransFormTypePay       = 8   // 付费
	TransFormTypeOther     = 9   // 其他
	TransFormTypeH5Pay     = 10  // H5付费
	TransFormTypeH5Pre     = 11  // H5预授信
	TransFormTypeQuick     = 12  // 快应用付费
	TransFormTypeQuick2    = 13  // 快应用加桌
	TransFormTypeOne       = 14  // 一元领取
	TransFormTypeFree      = 15  // 免费领取
	TransFormTypeCircle    = 16  // 调起朋友圈
	TransFormTypeNewSubmit = 101 // 表单提交（新）
	TransFormTypeNewKey    = 102 // 表单关键行为（新）
	TransFormTypeNewChat   = 103 // 有效咨询（新）
	TransFormTypeNewCopy   = 104 // 微信关注（新）
	TransFormTypeNewBuy    = 105 // 网页购买（新）
	TransFormTypeNewCall   = 106 // 电话拨打（新）
)

// PageType 落地页类型编码（pageType）
const (
	PageType = 7 // H5_API
)

// ClueDataItem 类型编码（type）
const (
	ClueDataItemTypeText     = 1  // 文本框
	ClueDataItemTypeNumber   = 2  // 数字输入框
	ClueDataItemTypePhone    = 3  // 电话
	ClueDataItemTypeEmail    = 4  // 邮箱
	ClueDataItemTypeRadio    = 5  // 单选
	ClueDataItemTypeCheckbox = 6  // 多选
	ClueDataItemTypeName     = 7  // 姓名
	ClueDataItemTypeSex      = 8  // 性别
	ClueDataItemTypeCity     = 9  // 城市
	ClueDataItemTypeDate     = 10 // 日期
	ClueDataItemTypeSelect   = 11 // 下拉单选
	ClueDataItemTypeSelects  = 12 // 下拉多选
)

// OppoAdCallbackConfig Oppo 广告回传参数
type OppoAdCallbackConfig struct {
	OwnerId int64  `json:"ownerId"` // 用户账号ID 营销平台的广告主ID
	ApiId   string `json:"apiId"`   // API授权接入方唯一身份标识
	ApiKey  string `json:"apiKey"`  // ApiKey 授权接入方的密钥
}

// OppoAdCallback Oppo 广告回传
type OppoAdCallback struct {
	Config    *OppoAdCallbackConfig
	TimeStamp int64
}

// NewOppoAdCallback 创建一个新的 OppoAdCallback
func NewOppoAdCallback(oppoAdCallbackParams *OppoAdCallbackConfig) *OppoAdCallback {
	return &OppoAdCallback{
		Config:    oppoAdCallbackParams,
		TimeStamp: time.Now().Unix(),
	}
}

// SendDataParamsItem 发送数据参数项
type SendDataParamsItem struct {
	Column  int64    `json:"column,omitempty"`
	Type    string   `json:"type,omitempty"`
	IfNeed  bool     `json:"ifNeed,omitempty"`
	Desc    string   `json:"desc,omitempty"`
	Value   int      `json:"value,omitempty"`
	Options []string `json:"options,omitempty"`
}

// SendDataParams 发送数据参数
type SendDataParams struct {
	PageId        int64                `json:"pageId,omitempty"`
	OwnerId       int64                `json:"ownerId,omitempty"`
	Ip            string               `json:"ip,omitempty"`
	Tid           string               `json:"tid,omitempty"`
	LbId          string               `json:"lbid,omitempty"`
	Items         []SendDataParamsItem `json:"items,omitempty"`
	TransformType int                  `json:"transformType,omitempty"`
	PageType      int                  `json:"pageType,omitempty"`
}

// SendDataResponse 发送数据响应
type SendDataResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// SendData 发送数据
func (o *OppoAdCallback) SendData(sendDataParams SendDataParams) (sendDataResponse SendDataResponse, err error) {
	sendDataParams.OwnerId = o.Config.OwnerId
	httpClient := NewHttpClient(sendDataURL)
	header := http.Header{"Content-Type": []string{"application/json"}, "Authorization": []string{fmt.Sprintf("%s %s", "Bearer", o.generateToken())}}
	postJson, err := httpClient.PostJsonAndHeader(sendDataParams, header)
	if err != nil {
		return SendDataResponse{}, err
	}
	err = json.Unmarshal([]byte(postJson), &sendDataResponse)
	if err != nil {
		return SendDataResponse{}, err
	}
	return sendDataResponse, nil
}

// generateToken 生成token token=base64(owner_id+“,”+api_id+“,”+time_stamp+“,”+sign)
func (o *OppoAdCallback) generateToken() string {
	str := strconv.FormatInt(o.Config.OwnerId, 10) + "," + o.Config.ApiId + "," + strconv.FormatInt(o.TimeStamp, 10) + "," + o.generateSign()
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// generateSign 生成秘钥 sign=sha1(api_id+api_key+time_stamp)
func (o *OppoAdCallback) generateSign() string {
	// 拼接字符串
	str := o.Config.ApiId + o.Config.ApiKey + strconv.FormatInt(o.TimeStamp, 10)
	// 计算SHA1哈希值
	h := sha1.New()
	h.Write([]byte(str))
	sign := fmt.Sprintf("%x", h.Sum(nil))
	return sign
}

//
