/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/21
   Description :
-------------------------------------------------
*/

package zenvconf

import (
    "errors"
    "fmt"
    "os"
    "reflect"

    "github.com/zlyuancn/zstr"
)

const DefaultTag = "env"
const JsonTag = "json"

type EnvConf struct {
    // 标记名
    tag string
    // 在没有默认标记时使用json标记, 标记值为空也被视为没有标记
    use_json_tag bool
    // 环境变量前缀
    prefix string
}

func NewEnvConf() *EnvConf {
    return &EnvConf{
        tag:          DefaultTag,
        use_json_tag: true,
    }
}

// 设置标记
func (m *EnvConf) SetTag(tag string) *EnvConf {
    m.tag = tag
    return m
}

// 在没有默认标记时使用json标记, 标记值为空也被视为没有标记
func (m *EnvConf) UseJsonTag(b bool) *EnvConf {
    m.use_json_tag = b
    return m
}

// 设置环境变量前缀
func (m *EnvConf) SetEnvPrefix(prefix string) *EnvConf {
    m.prefix = prefix
    return m
}

// 将环境变量解析到一个结构体中, 传入的值必须是指向结构体的指针
func (m *EnvConf) Parse(a interface{}) error {
    a_type := reflect.TypeOf(a)

    if a_type.Kind() != reflect.Ptr {
        return errors.New("a必须是一个指针")
    }
    a_type = a_type.Elem()

    if a_type.Kind() != reflect.Struct {
        return errors.New("a必须是指向struct的指针")
    }

    a_value := reflect.ValueOf(a).Elem()

    field_count := a_type.NumField()
    for i := 0; i < field_count; i++ {
        field := a_type.Field(i)
        if field.PkgPath != "" {
            continue
        }

        key := field.Tag.Get(m.tag)
        if key == "" && m.use_json_tag {
            key = field.Tag.Get(JsonTag)
        }
        if key == "" {
            key = field.Name
        }

        value, ok := os.LookupEnv(m.prefix + key)
        if !ok {
            continue
        }

        new_value := reflect.New(field.Type)

        err := zstr.Scan(value, new_value.Interface())
        if err != nil {
            return fmt.Errorf("<%s>转换失败: %s", key, err)
        }

        a_value.Field(i).Set(new_value.Elem())
    }
    return nil
}
