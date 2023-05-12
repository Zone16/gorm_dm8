package dm

import (
	"fmt"
	"github.com/Zone16/gorm/schema"
	"strings"
)

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

type Namer struct {
	schema.NamingStrategy
}

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func ConvertNameToFormat(n Namer, str string) string {
	//判断表前缀开头是否是小写，小写加双引号
	if n.TablePrefix != "" && n.TablePrefix[0] > 96 && n.TablePrefix[0] < 123 {
		n.TablePrefix = fmt.Sprintf(`"%s"`, n.TablePrefix)
	}

	//判断表名是否是小写
	for _, v := range str {
		if v > 96 && v < 123 {
			str = fmt.Sprintf(`"%s"`, strings.ToLower(toDBName(str)))
			break
		}
	}

	if n.TablePrefix != "" {
		return n.TablePrefix + "." + str
	}

	return str
}

func (n Namer) TableName(table string) (name string) {
	return ConvertNameToFormat(n, table)
}

func (n Namer) ColumnName(table, column string) (name string) {
	return ConvertNameToFormat(n, n.NamingStrategy.ColumnName(table, column))
}

func (n Namer) JoinTableName(table string) (name string) {
	return ConvertNameToFormat(n, n.NamingStrategy.JoinTableName(table))
}

func (n Namer) RelationshipFKName(relationship schema.Relationship) (name string) {
	return ConvertNameToFormat(n, n.NamingStrategy.RelationshipFKName(relationship))
}

func (n Namer) CheckerName(table, column string) (name string) {
	return ConvertNameToFormat(n, n.NamingStrategy.CheckerName(table, column))
}

func (n Namer) IndexName(table, column string) (name string) {
	tlc := strings.ToLower(column)

	cl := n.NamingStrategy.IndexName(table, column)
	if strings.Contains(tlc, "idx_"+strings.ToLower(table)) && strings.Contains(tlc, strings.ToLower(column)) {
		cl = column
	}

	return ConvertNameToFormat(n, cl)
}

func toDBName(name string) string {
	if name == "" {
		return ""
	}

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool // upper case == true
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase
	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}
	ret := buf.String()
	return ret
}
