package graphql_query

type Hello struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Context string `json:"'context'"`
}

func (hello *Hello) Query(id int, name string) (hellos []Hello, err error) {
	allHello := []Hello{
		{0, "user1", "测试值1"},
		{1, "user2", "测试值2"},
	}
	if name == "" {
		for _, v := range allHello {
			hellos = append(hellos, v)
		}
	} else {
		if id == 0 {
			id = 000
		}
		hello.Id = id
		hello.Name = name
		hello.Context = "欢迎使用graphql，用户:" + name + ""
		hellos = append(hellos, *hello)
	}
	return
}
