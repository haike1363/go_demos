package sql_tsdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const EmsDataFilter = "${EMS_DATA_FILTER}"
const EmsDataFilterHolder = "EMS_DATA_FILTER"
const EmsTPQuery = "EMS_TP:"
const EmsEmptyString = "EMS_EMPTY_STRING"
const EmsInstanceIdFilter = "${EMS_INSTANCE_ID}"
const EmsStartTimeFilter = "${EMS_START_TIME}"
const EmsEndTimeFilter = "${EMS_END_TIME}"
const EmsStartYmdFilter = "${EMS_START_YMD}"
const EmsEndYmdFilter = "${EMS_END_YMD}"

type Filter struct {
	Type     string `json:"type"`
	Tagk     string `json:"tagk"`
	TagkEn   string `json:"tagkEn"`
	TagkName string `json:"tagkName"`
	Filter   string `json:"filter"`
	GroupBy  bool   `json:"groupBy"`
}

func (thisObj *Filter) clone() *Filter {
	byteArray, _ := json.Marshal(thisObj)
	result := Filter{}
	_ = json.Unmarshal(byteArray, &result)
	return &result
}

type Query struct {
	Start        int64       `json:"start"`
	End          int64       `json:"end"`
	MsResolution bool        `json:"msResolution"`
	ShowQuery    bool        `json:"showQuery"`
	Queries      []*SubQuery `json:"queries"`
}

type SubQuery struct {
	Index          int                    `json:"index"`
	QueryId        string                 `json:"queryId"`
	Metric         string                 `json:"metric"`
	Aggregator     string                 `json:"aggregator"`
	Downsample     string                 `json:"downsample"`
	Filters        []*Filter              `json:"filters"`
	TimeOverOffset string                 `json:"time_over_offset"`
	TimeOverCount  int64                  `json:"time_over_count"`
	ValueFormatter string                 `json:"value_formatter"`
	TSDBMeta       *TSDBMeta              `json:"-"`
	Params         map[string]interface{} `json:"-"`
	OutStart       int64                  `json:"-"`
	OutEnd         int64                  `json:"-"`
	OutAlignDs     int64                  `json:"-"`
}

type SerialData struct {
	Metric        string                 `json:"metric"`
	Dps           map[string]interface{} `json:"dps"`
	Tags          map[string]string      `json:"tags"`
	AggregateTags []string               `json:"aggregate_tags"`
	Query         *SubQuery              `json:"query"`
	QueryId       string                 `json:"queryId"`
	LastTime      int64                  `json:"-"`
	LastValue     interface{}            `json:"-"`
	SumValue      interface{}            `json:"-"`
	MetricName    string                 `json:"metricName"`
}

type TSDBMeta struct {
	InstanceId   string
	SelectView    string
	TimeField     string
	TimeYMDField  string
	IsMillisecond bool
	SerialKeyList []string
	DBFetcher     func(sql string) (data []map[string]interface{}, err error) `json:"-"`
}

func QueryAsOpenTSDB(query *Query) (results []*SerialData, sampleTime int64, err error) {
	defer func() {
		if errPanic := recover(); errPanic != nil {
			results = nil
			err = fmt.Errorf("%v", errPanic)
		}
	}()
	if query.MsResolution {
		return nil, 0, fmt.Errorf("not support msResolution")
	}
	if query.Start > 9999999999 {
		query.Start = query.Start / 1000
	}
	if query.End > 9999999999 {
		query.End = query.End / 1000
	}
	if query.Start <= 0 {
		return nil, 0, fmt.Errorf("start <= 0")
	}
	if query.End == 0 {
		query.End = time.Now().Unix()
	}
	if query.Start > query.End {
		return nil, 0, fmt.Errorf("start > end")
	}
	sampleTime = 0
	results = make([]*SerialData, 0, 0)
	var wg sync.WaitGroup
	var lock sync.Mutex
	var errList []error
	for i, v := range query.Queries {
		if v.Index <= 0 {
			v.Index = i
		}
		wg.Add(1)
		go func(subQuery *SubQuery) {
			defer wg.Done()
			result, querySampleTime, err := doSubQuery(query, subQuery)
			lock.Lock()
			defer lock.Unlock()
			if err != nil {
				errList = append(errList, fmt.Errorf("do subQuery %v failed %v", subQuery.QueryId, err))
				return
			}
			results = append(results, result...)
			if len(result) > 0 && querySampleTime > sampleTime {
				sampleTime = querySampleTime
			}
		}(v)
	}
	wg.Wait()
	if len(errList) > 0 {
		return nil, 0, errList[0]
	}
	return results, sampleTime, nil
}

func doSubQuery(query *Query, subQuery *SubQuery) (results []*SerialData, dataLastTime int64, err error) {
	meta := subQuery.TSDBMeta
	serialKeyList := meta.SerialKeyList
	selectView := meta.SelectView
	timeField := meta.TimeField
	timeYMDField := meta.TimeYMDField
	if meta.IsMillisecond {
		return nil, 0, fmt.Errorf("not support millsecond time filed")
	}
	if len(serialKeyList) == 0 {
		return nil, 0, fmt.Errorf("SerialKeyList empty")
	}
	for _, v := range serialKeyList {
		if len(v) == 0 {
			return nil, 0, fmt.Errorf("SerialKeyList has empty serial key")
		}
	}
	if len(timeField) == 0 {
		return nil, 0, fmt.Errorf("TimeField is empty")
	}
	if len(selectView) == 0 {
		return nil, 0, fmt.Errorf("SelectView is empty")
	}
	dsSeconds, dsOp, fill, err := parseDownsample(subQuery.Downsample)
	if err != nil {
		return nil, 0, err
	}
	aggregator, err := parseAggregator(subQuery.Aggregator)
	if err != nil {
		return nil, 0, err
	}

	groupByTagKeys := make([]string, 0, 0)
	whereList := make([]string, 0, 0)

	// time filter
	timeRangeList := []TimeRange{
		{Start: query.Start, End: query.End},
	}
	if len(subQuery.TimeOverOffset) > 0 {
		ret, err := getTimeOverList(query.Start, query.End, subQuery.TimeOverOffset, subQuery.TimeOverCount)
		if err != nil {
			return nil, 0, err
		}
		timeRangeList = append(timeRangeList, ret...)
	}
	timeRangeStart := query.Start
	timeRangeEnd := query.End
	var timeConditionList []string
	for _, timeRange := range timeRangeList {
		if timeRange.Start < timeRangeStart {
			timeRangeStart = timeRange.Start
		}
		if timeRange.End > timeRangeEnd {
			timeRangeEnd = timeRange.End
		}
		var timeAndList []string
		timeAndList = append(timeAndList, fmt.Sprintf("%v<=%v", timeRange.Start, timeField))
		timeAndList = append(timeAndList, fmt.Sprintf("%v<=%v", timeField, timeRange.End))
		if !EmptyString(timeYMDField) {
			if timeYMDField != "ymd" {
				return nil, 0, fmt.Errorf("not support ymd filed %v", timeYMDField)
			}
			if timeRange.Start > 0 {
				startYMD := time.Unix(timeRange.Start, 0).Format("2006-01-02")
				timeAndList = append(timeAndList, fmt.Sprintf("%v>='%v'", timeYMDField, startYMD))
			}
			if timeRange.End > 0 {
				endYMD := time.Unix(timeRange.End, 0).Format("2006-01-02")
				timeAndList = append(timeAndList, fmt.Sprintf("%v<='%v'", timeYMDField, endYMD))
			}
		}
		timeConditionList = append(timeConditionList, fmt.Sprintf("(%v)", strings.Join(timeAndList, " AND ")))
	}
	timeCondition := strings.Join(timeConditionList, " OR ")
	whereList = append(whereList, fmt.Sprintf("(%v)", timeCondition))

	// tag filter
	filterList := processDupTag(subQuery)
	var fmtFilterList []*Filter
	var expFilterList []*Filter
	for _, filter := range filterList {
		if filter.GroupBy {
			groupByTagKeys = append(groupByTagKeys, filter.Tagk)
		}
		if len(filter.Filter) > 0 {
			parseFilterList := parseFilter(filter.Filter)
			if len(filter.Type) == 0 || filter.Type == "literal_or" {
				tagCondition := ""
				if len(parseFilterList) > 1 {
					tagCondition = fmt.Sprintf("(%v IN (%v))", filter.Tagk, strings.Join(parseFilterList, ", "))
				} else {
					tagCondition = fmt.Sprintf("(%v = %v)", filter.Tagk, parseFilterList[0])
				}
				whereList = append(whereList, tagCondition)
			} else if filter.Type == "iliteral_or" {
				tagCondition := ""
				if len(parseFilterList) > 1 {
					tagCondition = fmt.Sprintf("(%v NOT IN (%v))", filter.Tagk, strings.Join(parseFilterList, ", "))
				} else {
					tagCondition = fmt.Sprintf("(%v != %v)", filter.Tagk, parseFilterList[0])
				}
				whereList = append(whereList, tagCondition)
			} else if filter.Type == "wildcard" {
				if filter.Filter != "*" {
					whereList = append(whereList, fmt.Sprintf("(%v LIKE '%v')", filter.Tagk, strings.Replace(filter.Filter, "*", "%", -1)))
				}
			} else if filter.Type == "expression" {
				expFilterList = append(expFilterList, filter)
			} else if filter.Type == "weak_fmt_exp" || filter.Type == "fmt_exp" {
				fmtFilterList = append(fmtFilterList, filter)
			} else {
				return nil, 0, fmt.Errorf("not support filterType %v", filter.Type)
			}
		}
	}

	if len(fmtFilterList) > 0 {
		noNeedReplaceWhereString := strings.Join(whereList, " AND ")
		subQuery.Params[EmsDataFilterHolder] = noNeedReplaceWhereString
		for _, fmtFilter := range fmtFilterList {
			newFilter, missHolder, err := ReplaceHolder(fmtFilter.Filter, subQuery.Params)
			if err != nil {
				return nil, 0, err
			}
			if missHolder > 0 {
				if fmtFilter.Type == "weak_fmt_exp" {
					continue
				} else {
					return nil, 0, fmt.Errorf("tagk miss holder count %v Filter %v Params: %v", fmtFilter.Filter, missHolder, subQuery.Params)
				}
			}
			whereList = append(whereList, newFilter)
		}
	}
	// expFilterList 和 不会做用于fmtFilterList
	if len(expFilterList) > 0 {
		for _, v := range expFilterList {
			whereList = append(whereList, v.Filter)
		}
	}
	groupByKeyList := uniqueStringList(groupByTagKeys, serialKeyList)
	selectView = PushDownBaseCondition(selectView, timeRangeStart, timeRangeEnd, GetYMD(timeRangeStart), GetYMD(timeRangeEnd), subQuery.TSDBMeta.InstanceId)
	if strings.Contains(selectView, EmsDataFilter) {
		selectView = strings.Replace(selectView, EmsDataFilter, strings.Join(whereList, " AND "), -1)
	}
	selectViewWhere := fmt.Sprintf("%v WHERE %v", selectView, strings.Join(whereList, " AND "))

	subQuery.OutStart = timeRangeStart
	subQuery.OutEnd = timeRangeEnd
	subQuery.OutAlignDs = dsSeconds

	// get dataLastTime
	dataTimeSql := fmt.Sprintf("SELECT COALESCE(MAX(%v),-1) AS _ems_last_time, COALESCE(MIN(%v),-1) AS _ems_first_time FROM %v",
		timeField, timeField, selectViewWhere)
	dataTimePair, err := meta.DBFetcher(dataTimeSql)
	if err != nil {
		return nil, 0, err
	}
	// no data
	if len(dataTimePair) == 0 {
		return []*SerialData{}, 0, nil
	}
	v, found := dataTimePair[0]["_ems_last_time"]
	if !found {
		return nil, 0, fmt.Errorf("not found _ems_last_time")
	}
	dataLastTime, err = strconv.ParseInt(fmt.Sprint(v), 10, 64)
	if err != nil {
		return nil, 0, err
	}
	if dataLastTime == -1 {
		return []*SerialData{}, dataLastTime, nil
	}
	v, _ = dataTimePair[0]["_ems_first_time"]
	firstTime, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
	if err != nil {
		return nil, 0, err
	}
	// downsample
	selectList := make([]string, 0, 0)
	selectList = append(selectList, groupByKeyList...)

	downsampleSql := ""
	if dsOp == "none" {
		selectList = append(selectList, fmt.Sprintf("%v AS _ems_ts", timeField))
		selectList = append(selectList, fmt.Sprintf("%v AS _ems_val", subQuery.Metric))
		downsampleSql = fmt.Sprintf("SELECT %v FROM %v",
			strings.Join(selectList, ", "),
			selectViewWhere)
	} else if dsOp == "last" || dsOp == "first" {
		partOrder := "DESC"
		if dsOp == "first" {
			partOrder = "ASC"
		}
		tsAligned := ""
		var partitionKeyList []string
		partitionKeyList = append(partitionKeyList, serialKeyList...)
		if dsSeconds > 0 {
			tsAligned = GenAlignedTimeField(timeField, dsSeconds)
			partitionKeyList = append(partitionKeyList, tsAligned)
		} else {
			if dsOp == "last" {
				tsAligned = fmt.Sprint(dataLastTime)
			} else {
				tsAligned = fmt.Sprint(firstTime)
			}
		}
		selectViewWhere = fmt.Sprintf("(SELECT *,%v AS _ems_ts FROM (SELECT *,ROW_NUMBER() OVER (PARTITION BY %v ORDER BY %v %v) AS _ems_row_number FROM %v) _ems_row_number_t WHERE _ems_row_number=1) _ems_last_t",
			tsAligned, strings.Join(partitionKeyList, ", "), timeField, partOrder, selectViewWhere)
		selectList = append(selectList, "_ems_ts")
		selectList = append(selectList, fmt.Sprintf("%v AS _ems_last_ts", timeField))
		selectList = append(selectList, fmt.Sprintf("%v AS _ems_val", subQuery.Metric))
		downsampleSql = fmt.Sprintf("SELECT %v FROM %v\n",
			strings.Join(selectList, ", "),
			selectViewWhere)
	} else if isTPercentDownsample(dsOp) {
		percentValue, err := parseTPercentDownsample(dsOp)
		if err != nil {
			return nil, 0, err
		}
		percentValue = percentValue / 100
		partOrder := "ASC"
		tsAligned := ""
		var partitionKeyList []string
		partitionKeyList = append(partitionKeyList, serialKeyList...)
		if dsSeconds > 0 {
			tsAligned = GenAlignedTimeField(timeField, dsSeconds)
			partitionKeyList = append(partitionKeyList, tsAligned)
		} else {
			tsAligned = fmt.Sprint(dataLastTime)
		}
		partitionKeyString := strings.Join(partitionKeyList, ", ")
		selectViewWhere = fmt.Sprintf("(SELECT *,%v AS _ems_ts FROM (\nSELECT *,ROW_NUMBER() OVER (PARTITION BY %v ORDER BY %v %v) AS _ems_row_number,COUNT() OVER (PARTITION BY %v) AS _ems_count FROM %v\n) _ems_row_number_t WHERE _ems_row_number=CAST(CEIL(_ems_count*%v) AS INT)) _ems_precent_t",
			tsAligned, partitionKeyString, subQuery.Metric, partOrder, partitionKeyString, selectViewWhere, percentValue)
		selectList = append(selectList, "_ems_ts")
		selectList = append(selectList, fmt.Sprintf("%v AS _ems_val", subQuery.Metric))
		downsampleSql = fmt.Sprintf("SELECT %v FROM %v\n",
			strings.Join(selectList, ", "),
			selectViewWhere)
	} else {
		tsAligned := ""
		if dsSeconds > 0 {
			tsAligned = GenAlignedTimeField(timeField, dsSeconds)
		} else {
			tsAligned = fmt.Sprint(firstTime)
		}
		selectList = append(selectList, fmt.Sprintf("%v AS _ems_ts", tsAligned))
		dsEmsValStr := ""
		// select metric key
		if dsOp == "uniq" {
			dsEmsValStr = fmt.Sprintf("COUNT(DISTINCT(%v)) AS _ems_val", subQuery.Metric)
		} else {
			dsEmsValStr = fmt.Sprintf("%v(%v) AS _ems_val", dsOp, subQuery.Metric)
		}
		selectList = append(selectList, dsEmsValStr)
		downsampleSql = fmt.Sprintf("SELECT %v FROM %v\n GROUP BY %v",
			strings.Join(selectList, ", "),
			selectViewWhere,
			strings.Join(append([]string{"_ems_ts"}, groupByKeyList...), ", "))
	}

	sql := ""
	selectList = []string{}
	selectList = append(selectList, groupByTagKeys...)
	if aggregator == "none" {
		sql = downsampleSql
	} else {
		selectList = append(selectList, "_ems_ts")
		if aggregator == "uniq" {
			selectList = append(selectList, fmt.Sprintf("COUNT(DISTINCT(_ems_val)) AS _ems_val"))
		} else {
			selectList = append(selectList, fmt.Sprintf("%v(_ems_val) AS _ems_val", aggregator))
		}
		sql = fmt.Sprintf("SELECT %v FROM (\n%v\n) _ems_ds_table\n GROUP BY %v",
			strings.Join(selectList, ", "), downsampleSql, strings.Join(append([]string{"_ems_ts"}, groupByTagKeys...), ", "))
	}

	resultsMap := make(map[string]*SerialData)
	data, err := meta.DBFetcher(sql)
	if err != nil {
		return nil, 0, err
	}
	for _, line := range data {
		var emsTS int64 = -1
		var emsVal float64 = 0
		var tagsMap = map[string]string{}
		for k, v := range line {
			if k == "_ems_ts" {
				if emsTS == -1 {
					tsFloat, err := strconv.ParseFloat(fmt.Sprint(v), 64)
					if err != nil {
						return nil, 0, err
					}
					emsTS = int64(tsFloat)
				}
			} else if k == "_ems_val" {
				if len(subQuery.ValueFormatter) == 0 {
					emsVal, err = strconv.ParseFloat(fmt.Sprint(v), 64)
				} else {
					emsVal, err = strconv.ParseFloat(fmt.Sprintf(subQuery.ValueFormatter, v), 64)
				}
				if err != nil {
					return nil, 0, err
				}
			} else if k == "_ems_last_ts" {
				tsFloat, err := strconv.ParseFloat(fmt.Sprint(v), 64)
				if err != nil {
					return nil, 0, err
				}
				emsTS = int64(tsFloat)
			} else {
				tagsMap[k] = fmt.Sprint(v)
			}
		}
		var serialKey = genSerialKey(tagsMap)
		serialData, exist := resultsMap[serialKey]
		if !exist {
			var showSubQuery *SubQuery = nil
			if query.ShowQuery {
				showSubQuery = subQuery
			}
			serialData = &SerialData{
				Query:     showSubQuery,
				QueryId:   subQuery.QueryId,
				Metric:    subQuery.Metric,
				Tags:      tagsMap,
				Dps:       map[string]interface{}{},
				LastTime:  emsTS,
				LastValue: emsVal,
				SumValue:  emsVal,
			}
			resultsMap[serialKey] = serialData
		} else {
			if serialData.LastTime < emsTS {
				serialData.LastValue = emsVal
				serialData.LastTime = emsTS
			}
			if serialData.SumValue == nil {
				serialData.SumValue = emsVal
			} else {
				serialData.SumValue = serialData.SumValue.(float64) + emsVal
			}
		}
		serialData.Dps[fmt.Sprint(emsTS)] = emsVal
	}
	results = make([]*SerialData, 0, 0)
	for _, v := range resultsMap {
		results = append(results, v)
	}
	if dsSeconds > 0 && fill != nil {
		if fmt.Sprint(fill) == "null" {
			Fill(results, timeRangeStart, timeRangeEnd, dsSeconds, nil)
		} else {
			Fill(results, timeRangeStart, timeRangeEnd, dsSeconds, fill)
		}
	}
	results = replaceTagk(results, subQuery)
	results = orderSerials(results)
	return results, dataLastTime, nil
}

func replaceTagk(results []*SerialData, subQuery *SubQuery) []*SerialData {
	keyMap := map[string]string{}
	for _, filterItem := range subQuery.Filters {
		if len(filterItem.Tagk) > 0 && len(filterItem.TagkName) > 0 {
			keyMap[filterItem.Tagk] = filterItem.TagkName
		}
	}
	for _, resultItem := range results {
		newTags := map[string]string{}
		for oldTagk, tagv := range resultItem.Tags {
			if newTagk, found := keyMap[oldTagk]; found {
				newTags[newTagk] = tagv
			} else {
				newTags[oldTagk] = tagv
			}
		}
		resultItem.Tags = newTags
	}
	return results
}

func orderSerials(results []*SerialData) []*SerialData {
	type Item struct {
		Serial   *SerialData
		TagvList []string
	}
	var orderList []Item
	for _, v := range results {
		orderList = append(orderList, Item{
			Serial:   v,
			TagvList: getTagValueList(v.Tags),
		})
	}
	sort.Slice(orderList, func(i, j int) bool {
		return less(orderList[i].TagvList, orderList[j].TagvList)
	})
	for i := 0; i < len(results); i++ {
		results[i] = orderList[i].Serial
	}
	return results
}

func less(l, r []string) bool {
	for i := 0; i < len(l) && i < len(r); i++ {
		if l[i] < r[i] {
			return true
		}
		if l[i] > r[i] {
			return false
		}
	}
	return len(l) < len(r)
}

type TimeRange struct {
	Start int64
	End   int64
}

func getTimeOverList(start, end int64, offset string, count int64) ([]TimeRange, error) {
	timeUnit, err := GetTimeUnitValue(offset)
	if err != nil {
		return nil, err
	}
	var result []TimeRange
	for i := int64(1); i <= count; i++ {
		result = append(result, TimeRange{
			Start: start - timeUnit*i,
			End:   end - timeUnit*i,
		})
	}
	return result, nil
}

func GenAlignedTimeField(timeField string, dsSeconds int64) string {
	if dsSeconds >= 8*3600 {
		return fmt.Sprintf("(%v + 8*3600 - MOD(%v + 8*3600, %v) - 8*3600)", timeField, timeField, dsSeconds)
	}
	return fmt.Sprintf("(%v - MOD(%v, %v))", timeField, timeField, dsSeconds)
}

func genSerialKey(tagsMap map[string]string) string {
	var tagList []string
	for k, v := range tagsMap {
		tagList = append(tagList, fmt.Sprintf("%v=%v", k, v))
	}
	sort.Strings(tagList)
	return strings.Join(tagList, ",")
}

func getTagValueList(tagsMap map[string]string) []string {
	var tagkList []string
	for k := range tagsMap {
		tagkList = append(tagkList, k)
	}
	sort.Strings(tagkList)
	var tagvList []string
	for _, k := range tagkList {
		tagvList = append(tagvList, tagsMap[k])
	}
	sort.Strings(tagvList)
	return tagvList
}

func parseAggregator(aggregator string) (aggOp string, err error) {
	if aggregator == "avg" {
		aggOp = "AVG"
	} else if aggregator == "max" {
		aggOp = "MAX"
	} else if aggregator == "sum" {
		aggOp = "SUM"
	} else if aggregator == "min" {
		aggOp = "MIN"
	} else if aggregator == "count" {
		aggOp = "COUNT"
	} else if aggregator == "uniq" {
		aggOp = "uniq"
	} else if len(aggregator) == 0 || aggregator == "none" {
		aggOp = "none"
	} else {
		err = fmt.Errorf("not support aggregator %v", aggregator)
		return
	}
	return aggOp, nil
}

func isTPercentDownsample(ds string) bool {
	return strings.HasPrefix(ds, "tp")
}

func parseTPercentDownsample(ds string) (float64, error) {
	if !isTPercentDownsample(ds) {
		return 0, fmt.Errorf("bad PercentDownsample")
	}
	ret, err := strconv.ParseFloat(ds[2:], 64)
	if err != nil {
		return 0, nil
	}
	for ret > 100 {
		ret /= 10
	}
	return ret, nil
}

func parseDownsample(downsample string) (seconds int64, op string, fill interface{}, err error) {
	if len(downsample) == 0 || downsample == "none" {
		return 0, "none", nil, nil
	}
	seconds = 0
	op = ""
	fill = nil
	err = nil
	splits := strings.Split(downsample, "-")
	if len(splits) < 2 {
		err = fmt.Errorf("%v has less than 2 parts after split by '-'", downsample)
		return
	}
	if splits[0] == "0all" {
		seconds = 0
	} else {
		seconds, err = GetTimeUnitValue(splits[0])
		if err != nil {
			return 0, "", nil, err
		}
	}
	if splits[1] == "avg" {
		op = "AVG"
	} else if splits[1] == "max" {
		op = "MAX"
	} else if splits[1] == "sum" {
		op = "SUM"
	} else if splits[1] == "min" {
		op = "MIN"
	} else if splits[1] == "count" {
		op = "COUNT"
	} else if splits[1] == "uniq" {
		op = "uniq"
	} else if splits[1] == "last" {
		op = "last"
	} else if splits[1] == "first" {
		op = "first"
	} else if isTPercentDownsample(splits[1]) {
		op = splits[1]
	} else if splits[1] == "none" {
		op = "none"
	} else {
		err = fmt.Errorf("%v not support downsample operator %v", downsample, splits[1])
		return
	}
	if len(splits) > 2 {
		if fmt.Sprint(splits[2]) == "null" {
			fill = "null"
		} else if fmt.Sprint(splits[2]) == "none" {
			fill = nil
		} else {
			fillValue, err := strconv.ParseFloat(fmt.Sprint(splits[2]), 64)
			if err != nil {
				return 0, "", nil, err
			}
			fill = fillValue
		}
	}
	return seconds, op, fill, nil
}

func GetTimeUnitValue(s string) (int64, error) {
	timeUnit := s[len(s)-1]
	timeNum, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil || timeNum == 0 {
		err = fmt.Errorf("%v parse timeNum err %v", s, err)
		return 0, err
	}
	if timeUnit == 's' {
		timeNum *= 1
	} else if timeUnit == 'm' {
		timeNum *= 60
	} else if timeUnit == 'h' {
		timeNum *= 3600
	} else if timeUnit == 'd' {
		timeNum *= 24 * 3600
	} else if timeUnit == 'w' {
		timeNum *= 7 * 24 * 3600
	} else {
		err = fmt.Errorf("%v not support time unit %v", s, timeUnit)
		return 0, err
	}
	return timeNum, nil
}

func uniqueStringList(l1, l2 []string) []string {
	resultMap := make(map[string]string)
	for _, i := range l1 {
		resultMap[i] = i
	}
	for _, i := range l2 {
		resultMap[i] = i
	}
	result := make([]string, 0, 0)
	for k := range resultMap {
		result = append(result, k)
	}
	return result
}

func ReplaceHolder(s string, params map[string]interface{}) (string, int, error) {
	holderList, err := parsePlaceHolder(s)
	if err != nil {
		return "", 0, err
	}
	missHolder := len(holderList)
	for _, v := range holderList {
		if params != nil {
			val, found := params[v]
			if found {
				s = strings.Replace(s, fmt.Sprintf("${%v}", v), fmt.Sprint(val), -1)
				missHolder--
			}
		}
	}
	return s, missHolder, nil
}

func parseFilter(s string) (result []string) {
	var filterList []string
	if strings.Contains(s, "|") {
		filterList = strings.Split(s, "|")
	} else {
		filterList = append(filterList, s)
	}

	for _, f := range filterList {
		if f == EmsEmptyString {
			result = append(result, "''")
		} else {
			result = append(result, fmt.Sprintf("'%v'", f))
		}
	}
	return result
}

func parsePlaceHolder(s string) ([]string, error) {
	var result []string
	match := bytes.Buffer{}
	state := 0
	index := 0
	for i, c := range s {
		if state == 0 {
			if c == '$' {
				state = 1
				index = i
			}
		} else if state == 1 {
			if c == '{' {
				state = 2
			} else {
				state = 0
				match = bytes.Buffer{}
			}
		} else if state == 2 {
			if c == '}' {
				state = 0
				if match.Len() > 0 {
					result = append(result, match.String())
				}
				match = bytes.Buffer{}
			} else {
				match.WriteRune(c)
			}
		}
	}
	if state != 0 {
		return nil, fmt.Errorf("bad parsePlaceHolder failed %v index %v miss match %v", s, index, match.String())
	}
	return result, nil
}

func processDupTag(subQuery *SubQuery) []*Filter {
	literalOrFilterMap := map[string]*Filter{}
	var filterList []*Filter
	for _, filter := range subQuery.Filters {
		if len(filter.Filter) > 0 && (len(filter.Type) == 0 || filter.Type == "literal_or") && len(filter.Tagk) > 0 {
			v, found := literalOrFilterMap[filter.Tagk]
			if found {
				v.Filter += "|" + filter.Filter
			} else {
				literalOrFilterMap[filter.Tagk] = filter.clone()
			}
		} else {
			filterList = append(filterList, filter)
		}
	}
	for _, v := range literalOrFilterMap {
		filterList = append(filterList, v)
	}
	return filterList
}

func Fill(serials []*SerialData, start, end, alignDs int64, value interface{}) []*SerialData {
	tsMap := map[string]interface{}{}
	tsNoNeedFillMap := map[int64]interface{}{}
	for _, serial := range serials {
		for ts := range serial.Dps {
			tsMap[ts] = ""
			if alignDs > 0 {
				tsInt64, _ := strconv.ParseInt(ts, 10, 64)
				tsInt64 = AlignByDsSeconds(tsInt64, alignDs)
				tsNoNeedFillMap[tsInt64] = ""
			}
		}
	}
	if start > 0 && end > 0 && alignDs > 0 {
		start = AlignByDsSeconds(start, alignDs)
		end = AlignByDsSeconds(end, alignDs)
		for curr := start; curr <= end; curr += alignDs {
			_, found := tsNoNeedFillMap[curr]
			if !found {
				ts := fmt.Sprint(curr)
				tsMap[ts] = ts
			}
		}
	}
	for _, serial := range serials {
		for ts := range tsMap {
			_, found := serial.Dps[ts]
			if !found {
				serial.Dps[ts] = value
			}
		}
	}
	return serials
}

func AlignByDsSeconds(now int64, ds int64) int64 {
	if ds <= 0 {
		return now
	}
	if ds >= 8*3600 {
		now += 8 * 3600
		return now - now%ds - 8*3600
	}
	return now - now%ds
}

func EmptyString(s string) bool {
	if len(s) == 0 || strings.EqualFold(s, "null") || strings.EqualFold(s, "none") {
		return true
	}
	return false
}

func GetYMD(ts int64) string {
	if ts > int64(9999999999)*1000 {
		ts /= 1000000
	} else if ts > 9999999999 {
		ts /= 1000
	}
	return time.Unix(ts, 0).Format("2006-01-02")
}

func PushDownBaseCondition(selectView string, startTime int64, endTime int64, startYmd string, endYmd string, instanceId string) string {
	selectView = strings.ReplaceAll(selectView, EmsInstanceIdFilter, instanceId)
	selectView = strings.ReplaceAll(selectView, EmsStartTimeFilter, fmt.Sprint(startTime))
	selectView = strings.ReplaceAll(selectView, EmsEndTimeFilter, fmt.Sprint(endTime))
	selectView = strings.ReplaceAll(selectView, EmsStartYmdFilter, startYmd)
	selectView = strings.ReplaceAll(selectView, EmsEndYmdFilter, endYmd)
	return selectView
}
