package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/robfig/cron/v3"
	"plumelog/config"
	"plumelog/log"
	"plumelog/model"
	"plumelog/utils"
	"strings"
	"time"
)

var (
	Client *elastic.Client
)

const index_prefix = "plume_log_"

func init() {
	Client, _ = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(config.Conf.Url), elastic.SetBasicAuth(config.Conf.Elastic.Username, config.Conf.Elastic.Password))
	_, _, err := Client.Ping(config.Conf.Url).Do(context.Background())
	if err != nil {
		panic(err)
	}
	timeStr := time.Now().Format("20060102")
	createPlumelogIndex("plume_log_" + timeStr)
	createPlumelogServicesIndex()
	// 不用协程会卡住
	go indexJob()
}

// indexJob 添加定时任务，每天零点删除n天前的索引并创建明天的索引
func indexJob() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@daily", func() {
		log.Info.Println("执行es索引定时任务。。。")
		format := time.Now().AddDate(0, 0, -config.Conf.KeepDays).Format("20060102")
		tomorrow := time.Now().AddDate(0, 0, 1).Format("20060102")
		log.Info.Println(format)
		createPlumelogIndex("plume_log_" + tomorrow)
		Client.DeleteIndex("plume_log_" + format).Do(context.Background())
		defer func() {
			err := recover()
			if err != nil {
				log.Error.Println(err)
			}
		}()
	})
	c.Start()
}

func createPlumelogServicesIndex() {
	do, err := Client.IndexExists("plume_log_services").Do(context.Background())
	if err != nil {
		panic(err)
	}
	if do {
		return
	}
	index := `{
		"mappings": {
			"properties": {
				"serviceName": {"type": "keyword"}
			}
		}
	}`
	_, err = Client.CreateIndex("plume_log_services").BodyJson(index).Do(context.Background())
	if err != nil {
		log.Error.Println(err)
	}
}

func createPlumelogIndex(index string) {
	do, err := Client.IndexExists(index).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if do {
		return
	}
	mappings := `{
				"mappings": {
					"properties": {
						"appName": {"type":"keyword"},
						"className": {"type":"keyword"},
						"method": {"type":"keyword"},
						"env": {"type":"keyword"},
						"logLevel": {"type":"keyword"},
						"serverName": {"type":"keyword"},
						"traceId": {"type":"keyword"},
						"url": {"type":"keyword"},
						"costTime": {"type":"long"},
						"dtTime": {"type":"date", "format":"strict_date_optional_time||epoch_millis"},
						"seq": {"type":"long"}
					}
				},
				"settings": {
					"number_of_shards": 1,
					"number_of_replicas": 1,
					"refresh_interval": "30s",
					"index":{
								  "max_result_window": 100000
						  }
				}
			}`
	_, err = Client.CreateIndex(index).BodyJson(mappings).Do(context.Background())
	if err != nil {
		log.Error.Println(err)
	}
}

func QueryPage(req *model.PlumelogInfoPageReq) (*model.PlmelogInfoPageResp, error) {
	resp := &model.PlmelogInfoPageResp{}
	plumelogInfos := make([]*model.PlumelogInfo, 0, req.Size)
	var sort bool
	if req.TraceId != "" {
		sort = true
	}
	service := getSearchService(req)
	if req.ServiceName != "" {
		service.Sort("seq", sort)
	}
	result, err := service.From(req.From).Size(req.Size).Sort("dtTime", sort).Pretty(true).Do(context.Background())
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}
	for _, hit := range result.Hits.Hits {
		plumelogInfo := &model.PlumelogInfo{}
		json.Unmarshal(hit.Source, plumelogInfo)
		plumelogInfos = append(plumelogInfos, plumelogInfo)
	}
	resp.Total = result.TotalHits()
	resp.CurrentPage = req.From
	resp.Size = req.Size
	resp.Records = plumelogInfos
	return resp, err
}

func GetAppNames() []string {
	result, err := Client.Search("plume_log_services").Do(context.Background())
	if err != nil {
		log.Error.Println(err)
		return nil
	}
	appNames := make([]string, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		marshalJSON, _ := hit.Source.MarshalJSON()
		m := make(map[string]string)
		json.Unmarshal(marshalJSON, &m)
		serviceName := m["serviceName"]
		if serviceName != "" {
			appNames = append(appNames, serviceName)
		}
	}
	return appNames
}

func QueryErrors(req *model.ErrorsReq) ([]*model.ErrorsResp, error) {
	resp := make([]*model.ErrorsResp, 0)
	index := getIndex(req.BeginDate, req.EndDate)
	query := elastic.NewBoolQuery().Filter(elastic.NewTermQuery("logLevel", "ERROR")).Filter(elastic.NewRangeQuery("dtTime").
		Gte(req.BeginDate).Lte(req.EndDate))
	if req.Env != "" {
		query.Filter(elastic.NewTermQuery("env", req.Env))
	}
	aggregation := elastic.NewTermsAggregation().Field("appName")
	if req.AppName != "" {
		query.Filter(elastic.NewTermQuery("appName", req.AppName))
		aggregation = aggregation.Field("className")
	}
	if req.ClassName != "" {
		query.Filter(elastic.NewTermQuery("className", req.ClassName))
		aggregation = aggregation.Field("method")
	}
	aggregationResp := model.AggregationResp{}
	do, err := Client.Search(index...).Query(query).Aggregation("errors", aggregation).Size(0).Do(context.Background())
	json.Unmarshal(do.Aggregations["errors"], &aggregationResp)
	for _, bucket := range aggregationResp.Buckets {
		errorsResp := model.ErrorsResp{
			Key:   bucket.Key,
			Count: bucket.DocCount,
		}
		resp = append(resp, &errorsResp)
	}
	return resp, err
}

func getIndex(beginDate, endDate int64) []string {
	days := utils.BetweenDays(beginDate, endDate) + 1
	index := make([]string, 0, days)
	if utils.SameDay(beginDate, endDate) {
		index = append(index, index_prefix+utils.ParseDay(beginDate, "20060102"))
	} else {
		for i := 0; i < days; i++ {
			index = append(index, index_prefix+utils.NextDay(beginDate, i, "20060102"))
		}
	}
	return index
}

func DownloadLog(req *model.PlumelogInfoPageReq) ([]string, error) {
	result, err := getSearchService(req).From(0).Size(50000).Sort("dtTime", false).Pretty(true).Do(context.Background())
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}
	res := make([]string, 0)
	for _, hit := range result.Hits.Hits {
		plumelogInfo := &model.PlumelogInfo{}
		json.Unmarshal(hit.Source, plumelogInfo)
		formart := fmt.Sprintf("%s %s [%s] %s --- %s   : %s \n", plumelogInfo.DateTime, plumelogInfo.LogLevel, plumelogInfo.AppName, plumelogInfo.ThreadName, plumelogInfo.ClassName, plumelogInfo.Content)
		res = append(res, formart)
	}
	return res, err
}

func getSearchService(req *model.PlumelogInfoPageReq) *elastic.SearchService {
	index := getIndex(req.BeginDate, req.EndDate)
	query := elastic.NewBoolQuery().Filter(elastic.NewRangeQuery("dtTime").Gte(req.BeginDate).Lte(req.EndDate))
	if req.TraceId != "" {
		query.Filter(elastic.NewTermQuery("traceId", strings.TrimSpace(req.TraceId)))
	}
	if req.AppName != "" {
		query.Filter(elastic.NewTermQuery("appName", strings.TrimSpace(req.AppName)))
	}
	if req.Method != "" {
		query.Filter(elastic.NewTermQuery("method", strings.TrimSpace(req.Method)))
	}
	if req.Env != "" {
		query.Filter(elastic.NewTermQuery("env", strings.TrimSpace(req.Env)))
	}
	if req.LogLevel != "" {
		query.Filter(elastic.NewTermQuery("logLevel", req.LogLevel))
	}
	if req.Url != "" {
		query.Filter(elastic.NewTermQuery("url", strings.TrimSpace(req.Url)))
	}
	if req.ServiceName != "" {
		query.Must(elastic.NewMatchQuery("service", strings.TrimSpace(req.ServiceName)))
	}
	if req.Content != "" {
		query.Must(elastic.NewMatchQuery("content", strings.TrimSpace(req.Content)))
	}
	if len(req.CostTime) > 0 {
		rangeQuery := elastic.NewRangeQuery("costTime")
		for _, timeRange := range req.CostTime {
			switch timeRange.Express {
			case "gt":
				rangeQuery.Gt(timeRange.Value)
			case "gte":
				rangeQuery.Gte(timeRange.Value)
			case "lt":
				rangeQuery.Lt(timeRange.Value)
			case "lte":
				rangeQuery.Lte(timeRange.Value)
			}
		}
		query.Filter(rangeQuery)
	}
	return Client.Search(index...).Query(query)
}
