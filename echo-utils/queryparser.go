// TODO rewrite all
package utils

import (
	"regexp"
	"strconv"

	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
)

var (
	rFields = regexp.MustCompile(`fields\[(\w+)\]`)
	rIgnore = regexp.MustCompile(`ignore\[(\w+)\]`)
	rSort   = regexp.MustCompile(`sort\[(\w+)\]`)
	rIDs    = regexp.MustCompile(`ids\[(\w+)\]`)
)

type Data struct {
	Params Params `json:"params"`
}

type Params struct {
	Fields    map[string][]string `json:"fields"`
	Ignore    map[string][]string `json:"ignore"`
	Sort      []SortItem          `json:"sort"`
	Include   []string            `json:"include"`
	Page      Page                `json:"page"`
	Github    Github              `json:"github"`
	Aggregate Aggregate           `json:"aggregate"`

	Where map[string]Where `json:"where"`
	// use either logic here ?
	WhereCondition   Condition    `json:"whereCondition"`
	WhereExpressions []Expression `json:"whereExpressions"`

	Kpis    OneManyParams `json:"kpis"`
	TagsIDs TagsIDs       `json:"tagsIDs"`
}
type SortItem struct {
	Number   string `json:"table"`
	Field    string `json:"field"`
	Order    string `json:"order"`
	Priority string `json:"priority"`
}
type Page struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}
type Github struct {
	Location string `json:"location"`
}
type Aggregate struct {
	Count map[string]string `json:"count"`
	Max   map[string]string `json:"max"`
}

type Where struct {
	Operation   string       `json:"operation"` // and, or, not
	Expressions []Expression `json:"expressions"`
}
type Expression struct {
	Operation  string      `json:"operation"` // and, or, not
	Condition  Condition   `json:"condition"`
	Conditions []Condition `json:"conditions"`
}
type Condition struct {
	Table    string   `json:"table"`    // ex 'lists'
	Field    string   `json:"field"`    // ex 'id'
	Operator string   `json:"operator"` // ex '=', '<>', '>', '>=', '<', '<=', 'STARTS_WITH', 'CONTAINS', 'ENDS_WITH'
	Value    string   `json:"value"`    // 1
	Values   []string `json:"values"`   // [1,2,3]
}

// for KPIs in Params
type OneManyParams struct {
	Fields           []string     `json:"fields"`
	Sort             []SortItem   `json:"sort"`
	Page             Page         `json:"page"`
	Where            Where        `json:"where"`
	WhereCondition   Condition    `json:"whereCondition"`
	WhereExpressions []Expression `json:"whereExpressions"`
}

// for TagsIDs in Params
type TagsIDs struct {
	Tags      TaggedParams `json:"tags"`
	Lists     TaggedParams `json:"lists"`
	Companies TaggedParams `json:"companies"`
	Contacts  TaggedParams `json:"contacts"`
}
type TaggedParams struct {
	Fields []string   `json:"fields"`
	Sort   []SortItem `json:"sort"`
	Page   Page       `json:"page"`
	Where  Where      `json:"where"`
}

// the above so we can get SQL below
type SQL struct {
	Fields  []string `json:"fields"`
	Table   string   `json:"table"`
	Foreign string   `json:"string"`
	IDs     []int64  `json:"ids"`

	ListJoinSQL    string `json:"listJoinSQL"`    // for join
	TagsIdsJoinSQL string `json:"tagsIdsJoinSQL"` // for join
	KpiJoinSQL     string `json:"kpiJoinSQL"`     // for join
	FileJoinSQL    string `json:"fileJoinSQL"`    // for join

	WhereSQL string `json:"whereSQL"`

	ListWhereSQL    string `json:"listWhereSQL"`    // for join
	TagsIdsWhereSQL string `json:"tagsIdsWhereSQL"` // for join
	KpiWhereSQL     string `json:"kpiWhereSQL"`     // for join
	FileWhereSQL    string `json:"fileWhereSQL"`    // for join

	OrderSQL string `json:"orderSQL"`
	Offset   string `json:"offset"`
	Limit    string `json:"limit"`
}
type ParamIDs struct {
	IDs []string `json:"ids"`
}

func isOperator(operator string) bool {
	switch operator {
	case
		"=", "<>", ">", ">=", "<", "<=", "%X", "%%", "X%", "start_with", "end_with", "not in", "in", "not or", "or", "like", "ilike":
		return true
	}
	return false
}

func valuesNotEmpty(values []string) bool {
	// TODO do, errors/
	return true
}

func (condition *Condition) IsValid() bool {
	// TODO better error output
	if condition.Field != "" && (condition.Value != "" || (len(condition.Values) > 0 && valuesNotEmpty(condition.Values))) && isOperator(condition.Operator) {
		return true
	}
	return false
}

func (condition *Condition) Text() string {

	//fmt.Println(">>>>>>  condition *Condition Text()")
	//Pretty(*condition)

	//Table    string   `json:"table"`    // ex 'lists'
	//Field    string   `json:"field"`    // ex 'id'
	//Operator string   `json:"operator"` // ex '=', '<>', '>', '>=', '<', '<=', 'STARTS_WITH', 'CONTAINS', 'ENDS_WITH'
	//Value    string   `json:"value"`    // 1
	//Values   []string `json:"values"`   // [1,2,3]

	// TODO better error output
	if len(condition.Values) > 0 {
		// TODO int vs string
		return condition.Field + " " + condition.Operator + " (" + strings.Join(condition.Values, ",") + ")"
	}
	if condition.Operator == "like" {
		return fmt.Sprintf(" %s LIKE '%%%s%%' ", condition.Field, condition.Value)
	}
	if condition.Operator == "ilike" {
		return fmt.Sprintf(" %s ILIKE '%%%s%%' ", condition.Field, condition.Value)
	}
	if condition.Operator == "STARTS_WITH" {
		return fmt.Sprintf(" %s ILIKE '%s%%' ", condition.Field, condition.Value)
	}

	if condition.Operator == "=" {
		return fmt.Sprintf(" %s = '%s' ", condition.Field, condition.Value)
	}
	if condition.Operator == ">=" {
		return fmt.Sprintf(" %s >= '%s' ", condition.Field, condition.Value)
	}
	if condition.Operator == ">" {
		return fmt.Sprintf(" %s > '%s' ", condition.Field, condition.Value)
	}

	if condition.Operator == "<=" {
		return fmt.Sprintf(" %s <= '%s' ", condition.Field, condition.Value)
	}
	if condition.Operator == "<" {
		return fmt.Sprintf(" %s <= '%s' ", condition.Field, condition.Value)
	}

	if condition.Operator == "=!" {
		return fmt.Sprintf(" %s = %s ", condition.Field, condition.Value)
	}
	if condition.Operator == ">=!" {
		return fmt.Sprintf(" %s >= %s ", condition.Field, condition.Value)
	}
	if condition.Operator == ">!" {
		return fmt.Sprintf(" %s > %s ", condition.Field, condition.Value)
	}

	if condition.Operator == "<=!" {
		return fmt.Sprintf(" %s <= %s ", condition.Field, condition.Value)
	}
	if condition.Operator == "<!" {
		return fmt.Sprintf(" %s <= %s ", condition.Field, condition.Value)
	}

	return condition.Field + " " + condition.Operator + " " + condition.Value
}

func (expression *Expression) Text() string {

	//fmt.Println(">>> *Expression Text()")

	sql := " "

	// TODO fix check if either
	if expression.Condition.Field != "" {
		condition := expression.Condition

		//fmt.Println(">>>> expression.Condition.Field start:", expression.Condition.Field)
		//Pretty(expression.Condition)

		if condition.IsValid() {
			sql = " (" + condition.Text() + ") "
		}
		//fmt.Println(">>>> expression.Condition.Field end:")
	}

	if expression.Conditions != nil {

		//fmt.Println("expression.Conditions:")
		//Pretty(expression.Conditions)

		for key, condition := range expression.Conditions {
			fmt.Println("condition:", condition)
			if key == 0 {
				sql = sql + "(" + condition.Text() + ") "
			} else {
				sql = sql + expression.Operation + " (" + condition.Text() + ") "
			}
		}
	}

	return sql
}

type ConvertError struct {
	What  string
	Value string
}

func (e ConvertError) Error() string {
	return fmt.Sprintf("What:%v Value:%v", e.What, e.Value)
}

func (where *Where) ConvertWhereFields(table string, tags map[string]string) error {
	// TODO convert multiple tables, but throw error if illegal fields

	for key, _ := range where.Expressions {
		fmt.Println("ConvertWhereFields key:", key)

		// TODO fix check if either
		if where.Expressions[key].Condition.Field != "" {

			field := where.Expressions[key].Condition.Field
			fieldKey := strings.TrimPrefix(field, table+".")

			//fmt.Println("key:", key,"field:", field, " fieldKey:",fieldKey)

			if _, ok := tags[fieldKey]; !ok {
				return ConvertError{
					field,
					fieldKey,
				}
			}

			where.Expressions[key].Condition.Field = tags[fieldKey]

		}

		if where.Expressions[key].Conditions != nil {
			for key2, _ := range where.Expressions[key].Conditions {
				field := where.Expressions[key].Conditions[key2].Field
				fieldKey := strings.TrimPrefix(field, table+".")

				//fmt.Println("key:", key,"key2:", key2,"field:", field, " fieldKey:",fieldKey)
				if _, ok := tags[fieldKey]; !ok {
					return ConvertError{
						field,
						fieldKey,
					}
				}

				//fmt.Println("key:", key,"key2:", key2,"field:", field, " fieldKey:",fieldKey, " tags[fieldKey]:",tags[fieldKey])
				where.Expressions[key].Conditions[key2].Field = tags[fieldKey]
			}
		}

	}

	// TODO better erroring
	return nil
}

func (where *Condition) ConvertWhereCondition(table string, tags map[string]string) error {

	// TODO convert multiple tables, but throw error if illegal fields
	//for key, _ := range where.Expressions {
	//	field := where.Expressions[key].Condition.Field
	//	fieldKey := strings.TrimPrefix(field, table+".")
	//
	//	if _, ok := tags[fieldKey]; !ok {
	//		return ConvertError{
	//			field,
	//			fieldKey,
	//		}
	//	}
	//
	//	where.Expressions[key].Condition.Field = tags[fieldKey]
	//
	//	for key, _ := range where.Expressions[key].Conditions {
	//		field := where.Expressions[key].Conditions[key].Field
	//		fieldKey := strings.TrimPrefix(field, table+".")
	//
	//		if _, ok := tags[fieldKey]; !ok {
	//			return ConvertError{
	//				field,
	//				fieldKey,
	//			}
	//		}
	//
	//		where.Expressions[key].Conditions[key].Field = tags[fieldKey]
	//	}
	//}

	// TODO better erroring
	return nil
}

func (condition *Condition) ConvertWhereFields(table string, tags map[string]string) error {
	// TODO convert multiple tables, but throw error if illegal fields

	field := condition.Field
	fieldKey := strings.TrimPrefix(field, table+".")

	if _, ok := tags[fieldKey]; !ok {
		return ConvertError{
			field,
			fieldKey,
		}
	}
	condition.Field = tags[fieldKey]

	// TODO better erroring
	return nil
}

func (where *Where) Text() string {
	sql := " "
	fmt.Println("Where Text() where.Operation:", where.Operation)

	// TODO test or, not
	if where.Operation != "" && where.Expressions != nil && (where.Operation == "or" || where.Operation == "and" || where.Operation == "not") {

		if len(where.Expressions) > 1 {
			fmt.Println("Where len(where.Expressions) > 1")
			for key, expression := range where.Expressions {
				if key == 0 {
					sql = sql + " ( " + expression.Text() + " ) "
				} else {
					sql = sql + where.Operation + " ( " + expression.Text() + " ) "
				}
				fmt.Println("sql:", sql)
			}

		}
		if len(where.Expressions) == 1 {
			fmt.Println("Where len(where.Expressions) = 1")
			if where.Operation == "not" {
				sql = " not ( " + where.Expressions[0].Text() + ") "
			} else {
				sql = where.Expressions[0].Text()
			}
		}
	}
	return sql
}

func (s *SQL) Print() {
	fmt.Println("Fields:", s.Fields)
	fmt.Println("Table:", s.Table)
	fmt.Println("Foreign:", s.Foreign)
	fmt.Println("IDs:", s.IDs)

	fmt.Println("ListJoinSQL:", s.ListJoinSQL)
	fmt.Println("TagsIdsJoinSQL:", s.TagsIdsJoinSQL)
	fmt.Println("KpiJoinSQL:", s.KpiJoinSQL)
	fmt.Println("FileJoinSQL:", s.FileJoinSQL)

	fmt.Println("WhereSQL:", s.WhereSQL)

	fmt.Println("ListWhereSQL:", s.ListWhereSQL)
	fmt.Println("TagsIdsWhereSQL:", s.TagsIdsWhereSQL)
	fmt.Println("KpiWhereSQL:", s.KpiWhereSQL)
	fmt.Println("FileWhereSQL:", s.FileWhereSQL)

	fmt.Println("OrderSQL:", s.OrderSQL)
	fmt.Println("Offset:", s.Offset)
	fmt.Println("Limit:", s.Limit)

}

func ParseQuery(c echo.Context) *Params {

	if c.Request().Method == "POST" {
		data := new(Data)

		fmt.Println(c.Request().Body)

		if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
			fmt.Print(err)
			return nil
		}

		fmt.Printf("%+v\n", data)
		fmt.Println("=============")

		if data.Params.Page.Number <= 0 {
			data.Params.Page.Number = 1
		}
		if data.Params.Page.Size <= 0 {
			data.Params.Page.Size = 1000
		}

		return &data.Params
	}

	// TODO Parse Query
	// '=', '<>', '>', '>=', '<', '<=', 'STARTS_WITH', 'CONTAINS', 'ENDS_WITH'

	params := new(Params)
	params.Ignore = make(map[string][]string)
	params.Fields = make(map[string][]string)
	for p, a := range c.QueryParams() {
		switch true {
		case p == "include":
			params.Include = a
		case rFields.MatchString(p):
			match := rFields.FindStringSubmatch(p)
			t := match[1]
			params.Fields[t] = a
		case rIgnore.MatchString(p):
			match := rIgnore.FindStringSubmatch(p)
			t := match[1]
			params.Ignore[t] = a
			// TODO fix Sort became array
			//case rSort.MatchString(p):
			//	match := rSort.FindStringSubmatch(p)
			//	t := match[1]
			//	params.Sort[t] = a
		}
	}

	number, err := strconv.Atoi(c.QueryParam("page[number]"))
	if err != nil {
		params.Page.Number = 1
	}

	if number <= 0 {
		params.Page.Number = 1
	} else if number > 5000 {
		params.Page.Number = 5000
	} else {
		params.Page.Number = number
	}

	size, err := strconv.Atoi(c.QueryParam("page[size]"))
	if err != nil {
		params.Page.Size = 50
	}

	if size <= 0 {
		params.Page.Size = 50
	} else if size > 5000 {
		params.Page.Size = 5000
	} else {
		params.Page.Size = size
	}

	//params.WhereExpressions
	//params.WhereCondition

	// TODO add validation

	return params
}

func ParseQueryIDs(c echo.Context) *ParamIDs {

	params := new(ParamIDs)

	for p, a := range c.QueryParams() {
		switch true {
		case p == "ids":
			params.IDs = a
			//case rIDs.MatchString(p):
			//	match := rFields.FindStringSubmatch(p)
			//	t := match[1]
			//	params.IDs[t] = a
		}
	}

	// TODO add validation

	return params
}
