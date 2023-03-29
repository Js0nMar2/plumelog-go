package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

var ConnMap map[string]*websocket.Conn

type PlumelogInfo struct {
	DateTime   string `json:"dateTime,omitempty"`
	TraceId    string `json:"traceId,omitempty"`
	Method     string `json:"method,omitempty"`
	AppName    string `json:"appName,omitempty"`
	ServerName string `json:"serverName,omitempty"`
	ClassName  string `json:"className,omitempty"`
	Env        string `json:"env,omitempty"`
	Content    string `json:"content,omitempty"`
	ThreadName string `json:"threadName,omitempty"`
	Url        string `json:"url,omitempty"`
	SpanId     string `json:"spanId,omitempty"`
	DtTime     int64  `json:"dtTime,omitempty"`
	CostTime   int64  `json:"costTime,omitempty"`
	LogLevel   string `json:"logLevel,omitempty"`
	Seq        int32  `json:"seq,omitempty"`
}

func (plumelogInfo *PlumelogInfo) String() string {
	b, err := json.Marshal(*plumelogInfo)
	if err != nil {
		return fmt.Sprintf("%+v", *plumelogInfo)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *plumelogInfo)
	}
	return out.String()
}

type RequestInfo struct {
	Path          string
	Ip            string
	RequestMethod string
	ContentType   string
	IsRequestBody bool
	Time          string
	Token         string
	Header        map[string]string
}

func (reqInfo *RequestInfo) String() string {
	b, err := json.Marshal(*reqInfo)
	if err != nil {
		return fmt.Sprintf("%+v", *reqInfo)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *reqInfo)
	}
	return out.String()
}

type UserInfo struct {
	Username string
	Password string
}

type PlumelogInfoPageReq struct {
	BeginDate   int64 `json:"beginDate" binding:"required"`
	EndDate     int64 `json:"endDate" binding:"required"`
	AppName     string
	Env         string
	LogLevel    string
	TraceId     string
	Url         string
	Content     string
	ServiceName string
	Method      string
	CostTime    []CostTimeRange
	Size        int
	From        int
}

type CostTimeRange struct {
	Express string
	Value   int
}

type PlmelogInfoPageResp struct {
	Total       int64           `json:"total"`
	CurrentPage int             `json:"currentPage"`
	Size        int             `json:"size"`
	Records     []*PlumelogInfo `json:"records"`
}

type ErrorsReq struct {
	BeginDate int64  `json:"beginDate" binding:"required"`
	EndDate   int64  `json:"endDate" binding:"required"`
	AppName   string `json:"appName"`
	ClassName string `json:"className"`
	Env       string `json:"env"`
}

type ErrorsResp struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type AggregationResp struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []*Buckets
}

type Buckets struct {
	Key      string
	DocCount int `json:"doc_count"`
}
