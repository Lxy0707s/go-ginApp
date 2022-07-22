已知某个结构体，需要进行文件写，然后异常情况读取回写【代码片段】
    bytes, err := json.Marshal(info)
        if err != nil {
        return
    }
    err = myfile.WriteJSON(string(bytes), "test.json")
        if err != nil {
        return
    }
    bytess, _ := myfile.ReadAll("test.json")
    if err != nil {
    fmt.Println(err)
    return
    }
    var res = ""
    err = json.Unmarshal(bytess, &res)
    if err != nil {
        fmt.Println(err)
        return
    }
    f := []byte(res)
    err = myfile.WriteJSON(string(f), "test1.json") //可行
    err = myfile.WriteJSON(bytess, "test2.json")    //直接不行
数据压缩方案gzip
    var ss = make([]*db.DemoTask, 0)
    var info, _ = task.GetTaskInfo(c, args)
    // 压缩流程
    var buf bytes.Buffer
    jsonByte, _ := json.Marshal(info)
    err = gziptool.GzipWrite(&buf, jsonByte)
    if err != nil {
        log.Fatal(err)
    }
    // 解压流程
    var buf2 bytes.Buffer
    err = gziptool.GunzipWrite(&buf2, buf.Bytes())
    if err != nil {
        log.Fatal(err)
    }
    err = json.Unmarshal(buf2.Bytes(), &ss)
    if err != nil {
        log.Fatal(err)
        return
    }