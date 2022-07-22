package httptool

// 前9999的code都可以自定义
// 10000后的code将存在5位分段命名意义
// 前3位为业务类型码[注意业务码最多到655]
// 后2位为细节业务码
// 10102，其中101代表，业务类型[日志]；02代表，操作日志mysql添加操作
// 注意：每个新增code都得写上注释，不然拉出去祭天
const (
	OK               = 0   //请求成功响应码，自定义方式
	SUCCESS          = 200 //请求成功响应码，类http方式
	ERROR            = 500 //请求失败响应http状态码
	InvalidParams    = 400 //请求失败响应http状态码
	ErrNotFound      = 404 // 对应数据、资源不存在的状态码
	ErrorSSOValidate = 1001

	///////////////////分段业务码开始/////////////////
	// 100 API 相关
	ApiInfoAuthTypeAll   = 10021 //接口校验类型：需要进行默认权限校验,登录校验+接口权限校验
	ApiInfoAuthTypeLogin = 10022 //接口校验类型：仅需登录校验
	ApiInfoAuthTypeNo    = 10023 //接口校验类型：无需任何校验

	ApiInfoApiTypeParent   = 10030 //接口类型：父类型接口
	ApiInfoApiTypeMenu     = 10031 //接口类型：菜单类型接口
	ApiInfoApiTypeQuery    = 10032 //接口类型：查询子类型
	ApiInfoApiTypeMutation = 10033 //接口类型：修改子类型
	// 101 日志 相关
	LogTypeMysqlAdd    = 10100 //操作日志记录的mysql类型的实际动作：添加
	LogTypeMysqlUpdate = 10101 //操作日志记录的mysql类型的实际动作：更新
	LogTypeMysqlDelete = 10102 //操作日志记录的mysql类型的实际动作：删除
	// 200 响应校验 相关

	ErrorDb                       = 20000
	ErrorAPIAuthCheckTokenFail    = 20001
	ErrorAPIAuthCheckTokenTimeout = 20002
	ErrorAPIConfig                = 20003
	ErrorAPIAuth                  = 20004
	EmptyAPIAuth                  = 20005
	ErrorAPILimit                 = 20006

	ErrorUploadSaveImageFail = 30001
)
