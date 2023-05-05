package consts

type RunType string

const (
	FromServer RunType = "server"
	FromAgent  RunType = "agent"
)

func (e RunType) String() string {
	return string(e)
}

type ExecFromType string

const (
	FromCmd    ExecFromType = "cmd"
	FromClient ExecFromType = "client"
)

func (e ExecFromType) String() string {
	return string(e)
}

type WsMsgCategory string

const (
	ProgressInProgress WsMsgCategory = "in_progress"
	ProgressEnd        WsMsgCategory = "end"
	Result             WsMsgCategory = "result"
)

func (e WsMsgCategory) String() string {
	return string(e)
}

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"

	PATCH   HttpMethod = "PATCH"
	HEAD    HttpMethod = "HEAD"
	CONNECT HttpMethod = "CONNECT"
	OPTIONS HttpMethod = "OPTIONS"
	TRACE   HttpMethod = "TRACE"
)

func (e HttpMethod) String() string {
	return string(e)
}

type FormDataType string

const (
	FormDataTypeText FormDataType = "text"
	FormDataTypeFile FormDataType = "file"
)

func (e FormDataType) String() string {
	return string(e)
}

type HttpRespCode int

const (
	Continue          HttpRespCode = 100
	SwitchingProtocol HttpRespCode = 101

	OK                          HttpRespCode = 200
	Created                     HttpRespCode = 201
	Accepted                    HttpRespCode = 202
	NonAuthoritativeInformation HttpRespCode = 203
	NoContent                   HttpRespCode = 204
	ResetContent                HttpRespCode = 205
	PartialContent              HttpRespCode = 206

	MultipleChoice    HttpRespCode = 300
	MovedPermanently  HttpRespCode = 301
	Found             HttpRespCode = 302
	SeeOther          HttpRespCode = 303
	NotModified       HttpRespCode = 304
	UseProxy          HttpRespCode = 305
	unused            HttpRespCode = 306
	TemporaryRedirect HttpRespCode = 307
	PermanentRedirect HttpRespCode = 308

	BadRequest                   HttpRespCode = 400
	Unauthorized                 HttpRespCode = 401
	PaymentRequired              HttpRespCode = 402
	Forbidden                    HttpRespCode = 403
	NotFound                     HttpRespCode = 404
	MethodNotAllowed             HttpRespCode = 405
	NotAcceptable                HttpRespCode = 406
	ProxyAuthenticationRequired  HttpRespCode = 407
	RequestTimeout               HttpRespCode = 408
	Conflict                     HttpRespCode = 409
	Gone                         HttpRespCode = 410
	LengthRequired               HttpRespCode = 411
	PreconditionFailed           HttpRespCode = 412
	RequestEntityTooLarge        HttpRespCode = 413
	RequestURITooLong            HttpRespCode = 414
	UnsupportedMediaType         HttpRespCode = 415
	RequestedRangeNotSatisfiable HttpRespCode = 416
	ExpectationFailed            HttpRespCode = 417

	InternalServerError     HttpRespCode = 500
	Implemented             HttpRespCode = 501
	BadGateway              HttpRespCode = 502
	ServiceUnavailable      HttpRespCode = 503
	GatewayTimeout          HttpRespCode = 504
	HTTPVersionNotSupported HttpRespCode = 505
)

func (e HttpRespCode) Int() int {
	return int(e)
}

type HttpContentType string

const (
	ContentTypeJSON HttpContentType = "application/json"
	ContentTypeXML  HttpContentType = "application/xml"
	ContentTypeHTML HttpContentType = "text/html"
	ContentTypeTEXT HttpContentType = "text/text"

	ContentTypeFormData       HttpContentType = "multipart/form-data"
	ContentTypeFormUrlencoded HttpContentType = "application/x-www-form-urlencoded"

	ContentTypeUnixDir HttpContentType = "httpd/unix-directory"
)

func (e HttpContentType) String() string {
	return string(e)
}

type AuthorType string

const (
	BasicAuth   AuthorType = "basicAuth"
	BearerToken AuthorType = "bearerToken"
	OAuth2      AuthorType = "oAuth2"
	ApiKey      AuthorType = "apiKey"
)

func (e AuthorType) String() string {
	return string(e)
}

type GrantType string

const (
	AuthorizationCode         GrantType = "authorizationCode"
	AuthorizationCodeWithPKCE GrantType = "authorizationCodeWithPKCE"
	Implicit                  GrantType = "implicit"
	PasswordCredential        GrantType = "passwordCredential"
	ClientCredential          GrantType = "clientCredential"
)

func (e GrantType) String() string {
	return string(e)
}

type ClientAuthenticationWay string

const (
	SendAsBasicAuthHeader       ClientAuthenticationWay = "sendAsBasicAuthHeader"
	SendClientCredentialsInBody ClientAuthenticationWay = "sendClientCredentialsInBody"
)

func (e ClientAuthenticationWay) String() string {
	return string(e)
}

type HttpRespLangType string

const (
	LangJSON HttpRespLangType = "json"
	LangXML  HttpRespLangType = "xml"
	LangHTML HttpRespLangType = "html"
	LangTEXT HttpRespLangType = "text"
)

func (e HttpRespLangType) String() string {
	return string(e)
}

type HttpRespCharset string

const (
	UTF8 HttpRespCharset = "utf-8"
)

func (e HttpRespCharset) String() string {
	return string(e)
}

type FieldSource string

const (
	System FieldSource = "requirement"
	Custom FieldSource = "task"
)

func (e FieldSource) String() string {
	return string(e)
}

type FieldType string

const (
	Input       FieldType = "input"
	TextArea    FieldType = "textarea"
	Password    FieldType = "password"
	Checkbox    FieldType = "checkbox"
	Radio       FieldType = "radio"
	File        FieldType = "file"
	image       FieldType = "image"
	Hidden      FieldType = "hidden"
	Select      FieldType = "select"
	MultiSelect FieldType = "multiselect"

	Button FieldType = "button"
)

func (e FieldType) String() string {
	return string(e)
}

type FieldFormat string

const (
	PlainText FieldFormat = "plainText"
	RichText  FieldFormat = "richText"
)

func (e FieldFormat) String() string {
	return string(e)
}

type ProductStatus string

const (
	Active ProductStatus = "active"
	Closed ProductStatus = "closed"
)

func (e ProductStatus) String() string {
	return string(e)
}

type UsedBy string

const (
	InterfaceDebug UsedBy = "interface_debug"
	ScenarioDebug  UsedBy = "scenario_debug"
	//ScenarioExec   UsedBy = "scenario_exec" // not used
)

type ExtractorSrc string

const (
	Header ExtractorSrc = "header"
	Body   ExtractorSrc = "body"
)

type ExtractorType string

const (
	Boundary  ExtractorType = "boundary"
	JsonQuery ExtractorType = "jsonquery"
	HtmlQuery ExtractorType = "htmlquery"
	XmlQuery  ExtractorType = "xmlquery"
	Regx      ExtractorType = "regx"
	//FullText  ExtractorType = "fulltext"
)

type CheckpointType string

const (
	ResponseStatus CheckpointType = "responseStatus"
	ResponseHeader CheckpointType = "responseHeader"
	ResponseBody   CheckpointType = "responseBody"
	Extractor      CheckpointType = "extractor"
	Judgement      CheckpointType = "judgement"
)

type ExtractorScope string

const (
	Private ExtractorScope = "private" // in current interface
	Public  ExtractorScope = "public"  // shared by other interfaces in serve OR scenario
)

type ComparisonOperator string

const (
	Equal              ComparisonOperator = "equal"
	NotEqual           ComparisonOperator = "notEqual"
	GreaterThan        ComparisonOperator = "greaterThan"
	GreaterThanOrEqual ComparisonOperator = "greaterThanOrEqual"
	LessThan           ComparisonOperator = "lessThan"
	LessThanOrEqual    ComparisonOperator = "lessThanOrEqual"

	Contain    ComparisonOperator = "contain"
	NotContain ComparisonOperator = "notContain"
)

func (e ComparisonOperator) String() string {
	return string(e)
}

type ValueOperator string

const (
	Get   ValueOperator = "get"
	Set   ValueOperator = "set"
	Clear ValueOperator = "clear"
)

func (e ValueOperator) String() string {
	return string(e)
}

type ProgressStatus string

const (
	Start      ProgressStatus = "start"
	InProgress ProgressStatus = "in_progress"
	End        ProgressStatus = "end"
	Cancel     ProgressStatus = "cancel"
	Error      ProgressStatus = "error"
)

func (e ProgressStatus) String() string {
	return string(e)
}

type ResultStatus string

const (
	Pass    ResultStatus = "pass"
	Fail    ResultStatus = "fail"
	Skip    ResultStatus = "skip"
	Block   ResultStatus = "block"
	Unknown ResultStatus = "unknown"
)

func (e ResultStatus) String() string {
	return string(e)
}

type ProcessorCategory string

const (
	ProcessorRoot ProcessorCategory = "processor_root"
	//ProcessorThreadGroup ProcessorCategory = "processor_thread_group"

	ProcessorInterface ProcessorCategory = "processor_interface"
	ProcessorGroup     ProcessorCategory = "processor_group"
	ProcessorLogic     ProcessorCategory = "processor_logic"
	ProcessorLoop      ProcessorCategory = "processor_loop"
	ProcessorTimer     ProcessorCategory = "processor_timer"
	ProcessorPrint     ProcessorCategory = "processor_print"
	ProcessorVariable  ProcessorCategory = "processor_variable"
	ProcessorAssertion ProcessorCategory = "processor_assertion"
	ProcessorExtractor ProcessorCategory = "processor_extractor"

	ProcessorCookie ProcessorCategory = "processor_cookie"
	ProcessorData   ProcessorCategory = "processor_data"
)

func (e ProcessorCategory) ToString() string {
	return string(e)
}

type ProcessorType string

const (
	ProcessorRootDefault ProcessorType = "processor_root_default"
	//ProcessorThreadDefault ProcessorType = "processor_thread_default"

	ProcessorInterfaceDefault ProcessorType = "processor_interface_default"
	ProcessorGroupDefault     ProcessorType = "processor_group_default"
	ProcessorTimerDefault     ProcessorType = "processor_timer_default"
	ProcessorPrintDefault     ProcessorType = "processor_print_default"

	ProcessorLogicIf   ProcessorType = "processor_logic_if"
	ProcessorLogicElse ProcessorType = "processor_logic_else"

	ProcessorLoopTime  ProcessorType = "processor_loop_time"
	ProcessorLoopUntil ProcessorType = "processor_loop_until"
	ProcessorLoopIn    ProcessorType = "processor_loop_in"
	ProcessorLoopRange ProcessorType = "processor_loop_range"
	ProcessorLoopBreak ProcessorType = "processor_loop_break"

	ProcessorVariableSet ProcessorType = "processor_variable_set"
	//ProcessorVariableGet   ProcessorType = "processor_variable_get"
	ProcessorVariableClear ProcessorType = "processor_variable_clear"

	ProcessorAssertionDefault ProcessorType = "processor_assertion_default"
	//ProcessorAssertionEqual      ProcessorType = "processor_assertion_equal"
	//ProcessorAssertionNotEqual   ProcessorType = "processor_assertion_not_equal"
	//ProcessorAssertionContain    ProcessorType = "processor_assertion_contain"
	//ProcessorAssertionNotContain ProcessorType = "processor_assertion_not_contain"

	ProcessorExtractorBoundary  ProcessorType = "processor_extractor_boundary"
	ProcessorExtractorJsonQuery ProcessorType = "processor_extractor_jsonquery"
	ProcessorExtractorHtmlQuery ProcessorType = "processor_extractor_htmlquery"
	ProcessorExtractorXmlQuery  ProcessorType = "processor_extractor_xmlquery"

	ProcessorCookieGet   ProcessorType = "processor_cookie_get"
	ProcessorCookieSet   ProcessorType = "processor_cookie_set"
	ProcessorCookieClear ProcessorType = "processor_cookie_clear"

	ProcessorDataText  ProcessorType = "processor_data_text"
	ProcessorDataExcel ProcessorType = "processor_data_excel"
	//ProcessorDataZenData ProcessorType = "processor_data_zendata"
)

func (e ProcessorType) ToString() string {
	return string(e)
}

type LogType string

const (
	LogRoot      LogType = "root"
	LogInterface LogType = "interface"
	LogProcessor LogType = "processor"
)

func (e LogType) ToString() string {
	return string(e)
}

type ErrorAction string

const (
	ActionContinue        ErrorAction = "continue"
	ActionStartNextThread ErrorAction = "start_next_thread"
	ActionLoop            ErrorAction = "loop"
	ActionStopThread      ErrorAction = "stop_thread"
	ActionStopTest        ErrorAction = "stop_test"
	ActionStopTestNow     ErrorAction = "stop_test_now"
)

func (e ErrorAction) ToString() string {
	return string(e)
}

type DataSource string

const (
	Text  DataSource = "text"
	Excel DataSource = "excel"
	//ZenData DataSource = "zendata"
)

func (e DataSource) ToString() string {
	return string(e)
}

type TimeUnit string

const (
	Second TimeUnit = "sec"
	Minute TimeUnit = "min"
	Hour   TimeUnit = "hour"
)

func (e TimeUnit) ToString() string {
	return string(e)
}

type ExecType string

const (
	ExecStart ExecType = "start"
	ExecStop  ExecType = "stop"

	ExecScenario ExecType = "execScenario"
	ExecPlan     ExecType = "execPlan"
	ExecMessage  ExecType = "execMessage"
)

func (e ExecType) String() string {
	return string(e)
}

type DataType string

const (
	Int    DataType = "int"
	Float  DataType = "float"
	String DataType = "string"
)

func (e DataType) String() string {
	return string(e)
}

type RoleType string

const (
	Admin          RoleType = "admin"
	User           RoleType = "user"
	Tester         RoleType = "tester"
	Developer      RoleType = "developer"
	ProductManager RoleType = "product_manager"
)

func (e RoleType) String() string {
	return string(e)
}

type NodeType string

const (
	NodeElem    NodeType = "elem"
	NodeProp    NodeType = "prop"
	NodeContent NodeType = "content"
	NodeText    NodeType = "text"
)

func (e NodeType) String() string {
	return string(e)
}

type PlaceholderPrefix string

const (
	PlaceholderPrefixDatapool PlaceholderPrefix = "_dp"
	PlaceholderPrefixFunction PlaceholderPrefix = "_func"
)

func (e PlaceholderPrefix) String() string {
	return string(e)
}

type PlaceholderType string

const (
	PlaceholderTypeEnvironmentVariable PlaceholderType = "environment_variable"
	PlaceholderTypeVariable            PlaceholderType = "variable"
	PlaceholderTypeDatapool            PlaceholderType = "datapool"
	PlaceholderTypeFunction            PlaceholderType = "function"
)

func (e PlaceholderType) String() string {
	return string(e)
}

type ParamType string

const (
	ParamTypeString  ParamType = "string"
	ParamTypeNumber  ParamType = "number"
	ParamTypeInteger ParamType = "integer"
)

func (e ParamType) String() string {
	return string(e)
}

type ParamIn string

const (
	ParamInPath   ParamIn = "path"
	ParamInQuery  ParamIn = "query"
	ParamInHeader ParamIn = "header"
	ParamInCookie ParamIn = "cookie"
)

func (e ParamIn) String() string {
	return string(e)
}
