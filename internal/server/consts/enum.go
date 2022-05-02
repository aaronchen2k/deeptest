package serverConsts

type NodeCreateMode string

const (
	Brother NodeCreateMode = "brother"
	Child   NodeCreateMode = "child"
)

func (e NodeCreateMode) String() string {
	return string(e)
}

type NodeCreateType string

const (
	Dir  NodeCreateType = "dir"
	Node NodeCreateType = "node"
)

func (e NodeCreateType) String() string {
	return string(e)
}

type DropPos int

const (
	Before DropPos = -1
	Inner  DropPos = 0
	After  DropPos = 1
)

func (e DropPos) Int() int {
	return int(e)
}

type ExtractorSrc string

const (
	Header ExtractorSrc = "header"
	Body   ExtractorSrc = "body"
)

type ExtractorType string

const (
	FullText    ExtractorType = "fulltext"
	Regular     ExtractorType = "regular"
	XPath       ExtractorType = "xpath"
	JsonPath    ExtractorType = "jsonPath"
	CssSelector ExtractorType = "cssSelector"
	Boundary    ExtractorType = "boundary"
)

type CheckpointType string

const (
	ResponseStatus CheckpointType = "responseStatus"
	ResponseHeader CheckpointType = "responseHeader"
	ResponseBody   CheckpointType = "responseBody"
	Extractor      CheckpointType = "extractor"
)

type CheckpointOperator string

const (
	Contain            CheckpointOperator = "contain"
	Equal              CheckpointOperator = "="
	NotEqual           CheckpointOperator = "!="
	GreaterThan        CheckpointOperator = ">"
	LessThan           CheckpointOperator = ">"
	GreaterThanOrEqual CheckpointOperator = ">="
	LessThanOrEqual    CheckpointOperator = "<="
)
