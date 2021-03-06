package parameter

import (
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"net/url"
	"strconv"
	"strings"
)

type Parameters struct {
	Page        string
	PageInt     int
	PageSize    string
	PageSizeInt int
	SortField   string
	Columns     []string
	SortType    string
	Animation   bool
	URLPath     string
	Fields      map[string]string
}

const (
	Page     = "__page"
	PageSize = "__pageSize"
	Sort     = "__sort"
	SortType = "__sort_type"
	Columns  = "__columns"
	Prefix   = "__prefix"
	Pjax     = "_pjax"

	sortTypeDesc = "desc"
	sortTypeAsc  = "asc"

	IsAll      = "is_all"
	PrimaryKey = "pk"

	True  = "true"
	False = "false"

	FilterRangeParamStartSuffix = "_start__goadmin"
	FilterRangeParamEndSuffix   = "_end__goadmin"
	FilterParamJoinInfix        = "_goadmin_join_"
	FilterParamOperatorSuffix   = "__goadmin_operator__"
	FilterParamCountInfix       = "__goadmin_index__"
)

var operators = map[string]string{
	"like": "like",
	"gr":   ">",
	"gq":   ">=",
	"eq":   "=",
	"ne":   "!=",
	"le":   "<",
	"lq":   "<=",
	"free": "free",
}

var keys = []string{Page, PageSize, Sort, Columns, Prefix, Pjax, form.NoAnimationKey}

func BaseParam() Parameters {
	return Parameters{Page: "1", PageSize: "1", Fields: make(map[string]string)}
}

func GetParam(u *url.URL, defaultPageSize int, p ...string) Parameters {
	values := u.Query()

	primaryKey := "id"
	defaultSortType := "desc"

	if len(p) > 0 {
		primaryKey = p[0]
		defaultSortType = p[1]
	}

	page := getDefault(values, Page, "1")
	pageSize := getDefault(values, PageSize, strconv.Itoa(defaultPageSize))
	sortField := getDefault(values, Sort, primaryKey)
	sortType := getDefault(values, SortType, defaultSortType)
	columns := getDefault(values, Columns, "")

	animation := true
	if values.Get(form.NoAnimationKey) == "true" {
		animation = false
	}

	fields := make(map[string]string)

	for key, value := range values {
		if !modules.InArray(keys, key) && value[0] != "" {
			if key == SortType {
				if value[0] != sortTypeDesc && value[0] != sortTypeAsc {
					fields[key] = sortTypeDesc
				}
			} else {
				if strings.Contains(key, FilterParamOperatorSuffix) &&
					values.Get(strings.Replace(key, FilterParamOperatorSuffix, "", -1)) == "" {
					continue
				}
				fields[key] = value[0]
			}
		}
	}

	columnsArr := make([]string, 0)
	if columns != "" {
		columnsArr = strings.Split(columns, ",")
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	return Parameters{
		Page:        page,
		PageSize:    pageSize,
		PageSizeInt: pageSizeInt,
		PageInt:     pageInt,
		URLPath:     u.Path,
		SortField:   sortField,
		SortType:    sortType,
		Fields:      fields,
		Animation:   animation,
		Columns:     columnsArr,
	}
}

func GetParamFromUrl(urlStr string, defaultPageSize int, defaultSortType, primaryKey string, fromList ...bool) Parameters {

	if len(fromList) > 0 && !fromList[0] {
		return BaseParam()
	}

	u, err := url.Parse(urlStr)

	if err != nil {
		return BaseParam()
	}

	return GetParam(u, defaultPageSize, primaryKey, defaultSortType)
}

func (param Parameters) WithPK(id ...string) Parameters {
	param.Fields["pk"] = strings.Join(id, ",")
	return param
}

func (param Parameters) PK() []string {
	return strings.Split(param.Fields[PrimaryKey], ",")
}

func (param Parameters) IsAll() bool {
	return param.Fields[IsAll] == True
}

func (param Parameters) WithIsAll(isAll bool) Parameters {
	if isAll {
		param.Fields[IsAll] = True
	} else {
		param.Fields[IsAll] = False
	}
	return param
}

func (param Parameters) GetFilterFieldValueStart(field string) string {
	return param.Fields[field] + FilterRangeParamStartSuffix
}

func (param Parameters) GetFilterFieldValueEnd(field string) string {
	return param.Fields[field] + FilterRangeParamEndSuffix
}

func (param Parameters) GetFieldValue(field string) string {
	return param.Fields[field]
}

func (param Parameters) GetFieldOperator(field, suffix string) string {
	if param.Fields[field+FilterParamOperatorSuffix+suffix] == "" {
		return "eq"
	}
	return param.Fields[field+FilterParamOperatorSuffix+suffix]
}

func (param Parameters) Join() string {
	p := param.GetFixedParamStr()
	p.Add(Page, param.Page)
	return p.Encode()
}

func (param Parameters) SetPage(page string) Parameters {
	param.Page = page
	return param
}

func (param Parameters) GetRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add(Page, param.Page)
	return "?" + p.Encode()
}

func (param Parameters) GetRouteParamStrWithoutPageSize() string {
	p := url.Values{}
	p.Add(Sort, param.SortField)
	p.Add(Page, param.Page)
	p.Add(SortType, param.SortType)
	if len(param.Columns) > 0 {
		p.Add(Columns, strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p.Add(key, value)
	}
	return "?" + p.Encode()
}

func (param Parameters) GetLastPageRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add(Page, strconv.Itoa(param.PageInt-1))
	return "?" + p.Encode()
}

func (param Parameters) GetNextPageRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add(Page, strconv.Itoa(param.PageInt+1))
	return "?" + p.Encode()
}

func (param Parameters) GetFixedParamStr() url.Values {
	p := url.Values{}
	p.Add(Sort, param.SortField)
	p.Add(PageSize, param.PageSize)
	p.Add(SortType, param.SortType)
	if len(param.Columns) > 0 {
		p.Add(Columns, strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		if key != constant.EditPKKey && key != constant.DetailPKKey {
			p.Add(key, value)
		}
	}
	return p
}

func (param Parameters) Statement(wheres, delimiter string, whereArgs []interface{}, columns, existKeys []string,
	filterProcess func(string, string, string) string, getJoinTable func(string) string) (string, []interface{}, []string) {
	var multiKey = make(map[string]uint8)
	for key, value := range param.Fields {

		keyIndexSuffix := ""

		keyArr := strings.Split(key, FilterParamCountInfix)

		if len(keyArr) > 1 {
			key = keyArr[0]
			keyIndexSuffix = FilterParamCountInfix + keyArr[1]
		}

		if keyIndexSuffix != "" {
			multiKey[key] = 0
		} else if _, exist := multiKey[key]; !exist && modules.InArray(existKeys, key) {
			continue
		}

		var op string
		if strings.Contains(key, FilterRangeParamEndSuffix) {
			key = strings.Replace(key, FilterRangeParamEndSuffix, "", -1)
			op = "<="
		} else if strings.Contains(key, FilterRangeParamStartSuffix) {
			key = strings.Replace(key, FilterRangeParamStartSuffix, "", -1)
			op = ">="
		} else if !strings.Contains(key, FilterParamOperatorSuffix) {
			op = operators[param.GetFieldOperator(key, keyIndexSuffix)]
		}

		if modules.InArray(columns, key) {
			wheres += modules.FilterField(key, delimiter) + " " + op + " ? and "
			if op == "like" && !strings.Contains(value, "%") {
				whereArgs = append(whereArgs, "%"+filterProcess(key, value, keyIndexSuffix)+"%")
			} else {
				whereArgs = append(whereArgs, value)
			}
		} else {
			keys := strings.Split(key, FilterParamJoinInfix)
			if len(keys) > 1 {
				if joinTable := getJoinTable(keys[1]); joinTable != "" {
					value := filterProcess(key, value, keyIndexSuffix)
					wheres += joinTable + "." + modules.FilterField(keys[1], delimiter) + " " + op + " ? and "
					if op == "like" && !strings.Contains(value, "%") {
						whereArgs = append(whereArgs, "%"+value+"%")
					} else {
						whereArgs = append(whereArgs, value)
					}
				}
			}
		}

		existKeys = append(existKeys, key)
	}

	if len(wheres) > 3 {
		wheres = wheres[:len(wheres)-4]
	}

	return wheres, whereArgs, existKeys
}

func getDefault(values url.Values, key, def string) string {
	value := values.Get(key)
	if value == "" {
		return def
	}
	return value
}
