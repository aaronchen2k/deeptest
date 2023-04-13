

/**
 * 请求方法选择项
 * */
export const requestMethodOpts = [
    {
        label: "GET",
        value: "GET",
    },
    {
        label: "POST",
        value: "POST",
    },
    {
        label: "PUT",
        value: "PUT",
    },
    {
        label: "PATCH",
        value: "PATCH",
    },
    {
        label: "DELETE",
        value: "DELETE",
    },
    {
        label: "HEAD",
        value: "HEAD",
    },
    {
        label: "OPTIONS",
        value: "OPTIONS",
    },
    {
        label: "TRACE",
        value: "TRACE",
    },
]

/**
 * 响应码枚举
 * */
export const repCodeOpts = [
    {
        label: "200",
        value: "200",
    },
    {
        label: "404",
        value: "404",
    },
    {
        label: "500",
        value: "500",
    },
    {
        label: "501",
        value: "501",
    },
    {
        label: "502",
        value: "502",
    },
    {
        label: "503",
        value: "503",
    },
]

/**
 * 所有 http 响应状态码及描述
 * */
export const responseCodes = [
    {
        "description": "100: Continue",
        "value": "100",
        "label": "100",
    },
    {
        "description": "101: Switching Protocols",
        "value": "101",
        "label": "101"
    },
    {
        "description": "200: OK",
        "value": "200",
        "enabled": true,
        "label": "200"
    },
    {
        "description": "201: Created",
        "value": "201",
        "label": "201"
    },
    {
        "description": "202: Accepted",
        "value": "202",
        "label": "202"
    },
    {
        "description": "203: Non-Authoritative Information",
        "value": "203",
        "label": "203"
    },
    {
        "description": "204: No Content",
        "value": "204",
        "label": "204"
    },
    {
        "description": "205: Reset Content",
        "value": "205",
        "label": "205"
    },
    {
        "description": "206: Partial Content",
        "value": "206",
        "label": "206"
    },
    {
        "description": "207: Multi-Status (WebDAV)",
        "value": "207",
        "label": "207"
    },
    {
        "description": "208: Already Reported (WebDAV)",
        "value": "208",
        "label": "208"
    },
    {
        "description": "226: IM Used",
        "value": "226",
        "label": "226"
    },
    {
        "description": "300: Multiple Choices",
        "value": "300",
        "label": "300"
    },
    {
        "description": "301: Moved Permanently",
        "value": "301",
        "label": "301"
    },
    {
        "description": "302: Found",
        "value": "302",
        "label": "302"
    },
    {
        "description": "303: See Other",
        "value": "303",
        "label": "303"
    },
    {
        "description": "304: Not Modified",
        "value": "304",
        "label": "304"
    },
    {
        "description": "305: Use Proxy",
        "value": "305",
        "label": "305"
    },
    {
        "description": "306: (Unused)",
        "value": "306",
        "label": "306"
    },
    {
        "description": "307: Temporary Redirect",
        "value": "307",
        "label": "307"
    },
    {
        "description": "308: Permanent Redirect (experiemental)",
        "value": "308",
        "label": "308"
    },
    {
        "description": "400: Bad Request",
        "value": "400",
        "enabled": true,
        "label": "400"
    },
    {
        "description": "401: Unauthorized",
        "value": "401",
        "label": "401"
    },
    {
        "description": "402: Payment Required",
        "value": "402",
        "label": "402"
    },
    {
        "description": "403: Forbidden",
        "value": "403",
        "enabled": true,
        "label": "403"
    },
    {
        "description": "404: Not Found",
        "value": "404",
        "label": "404"
    },
    {
        "description": "405: Method Not Allowed",
        "value": "405",
        "label": "405"
    },
    {
        "description": "406: Not Acceptable",
        "value": "406",
        "label": "406"
    },
    {
        "description": "407: Proxy Authentication Required",
        "value": "407",
        "label": "407"
    },
    {
        "description": "408: Request Timeout",
        "value": "408",
        "label": "408"
    },
    {
        "description": "409:Conflict",
        "value": "409",
        "label": "409"
    },
    {
        "description": "410: Gone",
        "value": "410",
        "label": "410"
    },
    {
        "description": "411: Length Required",
        "value": "411",
        "label": "411"
    },
    {
        "description": "412: Precondition Failed",
        "value": "412",
        "label": "412"
    },
    {
        "description": "413: Request Entity Too Large",
        "value": "413",
        "label": "413"
    },
    {
        "description": "414:Request-URI Too Long",
        "value": "414",
        "label": "414"
    },
    {
        "description": "415: Unsupported Media Type",
        "value": "415",
        "label": "415"
    },
    {
        "description": "416: Requested Range Not Satisfiable",
        "value": "416",
        "label": "416"
    },
    {
        "description": "417: Expectation Failed",
        "value": "417",
        "label": "417"
    },
    {
        "description": "418: I'm a teapot (RFC 2324)",
        "value": "418",
        "label": "418"
    },
    {
        "description": "420: Enhance Your Calm (Twitter)",
        "value": "420",
        "label": "420"
    },
    {
        "description": "422: Unprocessable Entity (WebDAV)",
        "value": "422",
        "label": "422"
    },
    {
        "description": "423: Locked (WebDAV)",
        "value": "423",
        "label": "423"
    },
    {
        "description": "424: Failed Dependency (WebDAV)",
        "value": "424",
        "label": "424"
    },
    {
        "description": "425: Too Early",
        "value": "425",
        "label": "425"
    },
    {
        "description": "426: Upgrade Required",
        "value": "426",
        "label": "426"
    },
    {
        "description": "428: Precondition Required",
        "value": "428",
        "label": "428"
    },
    {
        "description": "429: Too Many Requests",
        "value": "429",
        "label": "429"
    },
    {
        "description": "431: Request Header Fields Too Large",
        "value": "431",
        "label": "431"
    },
    {
        "description": "444: No Response (Nginx)",
        "value": "444",
        "label": "444"
    },
    {
        "description": "449: Retry With (Microsoft)",
        "value": "449",
        "label": "449"
    },
    {
        "description": "450: Blocked by Windows Parental Controls (Microsoft)",
        "value": "450",
        "label": "450"
    },
    {
        "description": "451: Unavailable For Legal Reasons",
        "value": "451",
        "label": "451"
    },
    {
        "description": "499: Client Closed Request (Nginx)",
        "value": "499",
        "label": "499"
    },
    {
        "description": "500: Internal Server Error",
        "value": "500",
        "label": "500",
        "enabled": true,
    },
    {
        "description": "501: Not Implemented",
        "value": "501",
        "label": "501",
    },
    {
        "description": "502: Bad Gateway",
        "value": "502",
        "label": "502"
    },
    {
        "description": "503: Service Unavailable",
        "value": "503",
        "label": "503",
        "enabled": true,
    },
    {
        "description": "504: Gateway Timeout",
        "value": "504",
        "label": "504",
    },
    {
        "description": "505: HTTP Version Not Supported",
        "value": "505",
        "label": "505"
    },
    {
        "description": "506: Variant Also Negotiates (Experimental)",
        "value": "506",
        "label": "506"
    },
    {
        "description": "507: Insufficient Storage (WebDAV)",
        "value": "507",
        "label": "507"
    },
    {
        "description": "508: Loop Detected (WebDAV)",
        "value": "508",
        "label": "508"
    },
    {
        "description": "509: Bandwidth Limit Exceeded (Apache)",
        "value": "509",
        "label": "509"
    },
    {
        "description": "510: Not Extended",
        "value": "510",
        "label": "510"
    },
    {
        "description": "511: Network Authentication Required",
        "value": "511",
        "label": "511"
    },
    {
        "description": "598: Network read timeout error",
        "value": "598",
        "label": "598"
    },
    {
        "description": "599: Network connect timeout error",
        "value": "599",
        "label": "599"
    }
];
/**
 * path params 数据类型
 * */
export const pathParamsDataTypesOpts = [
    {
        label: "any",
        value: "any",
    },
    {
        label: "string",
        value: "string",
    },
    {
        label: "number",
        value: "number",
    },
    {
        label: "integer",
        value: "integer",
    },
    {
        label: "boolean",
        value: "boolean",
    },
    {
        label: "array",
        value: "array",
    },
]

/**
 * openApi 数据类型
 * */
export const openApiDataTypesOpts = [
    {
        label: "any",
        value: "any",
    },
    {
        label: "string",
        value: "string",
    },
    {
        label: "number",
        value: "number",
    },
    {
        label: "interger",
        value: "interger",
    },
    {
        label: "boolean",
        value: "boolean",
    },
    {
        label: "any",
        value: "any",
    },
    {
        label: "object",
        value: "object",
    },
]

export const paramsSchemaDataTypes: any = {
    "any": {
        label: "any",
        value: "any",
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                },
            ]
        }
    },
    "string": {
        label: "string",
        value: "string",
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                },
                {
                    label: 'format',
                    name: 'format',
                    type: 'string',
                    component: 'select',
                    placeholder: 'select a value',
                    options: [
                        {
                            label: 'data-time',
                            value: 'data-time',
                        },
                        {
                            label: 'time',
                            value: 'time',
                        },
                        {
                            label: 'email',
                            value: 'email',
                        },
                        {
                            label: 'idn-email',
                            value: 'idn-email',
                        },
                        {
                            label: 'hostname',
                            value: 'hostname',
                        },
                        {
                            label: 'idn-hostname',
                            value: 'idn-hostname',
                        },
                        {
                            label: 'ipv4',
                            value: 'ipv4',
                        },
                        {
                            label: 'ipv6',
                            value: 'ipv6',
                        },
                        {
                            label: 'uri',
                            value: 'uri',
                        },
                        {
                            label: 'uri-reference',
                            value: 'uri-reference',
                        },
                        {
                            label: 'iri',
                            value: 'iri',
                        },
                        {
                            label: 'iri-reference',
                            value: 'iri-reference',
                        },
                        {
                            label: 'uri-template',
                            value: 'uri-template',
                        },
                        {
                            label: 'json-pointer',
                            value: 'json-pointer',
                        },
                        {
                            label: 'regex',
                            value: 'regex',
                        },
                        {
                            label: 'uuid',
                            value: 'uuid',
                        },
                        {
                            label: 'password',
                            value: 'password',
                        },
                        {
                            label: 'byte',
                            value: 'byte',
                        },
                    ],
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                },
                {
                    label: 'pattern',
                    name: 'pattern',
                    type: 'string',
                    component: 'input',
                    placeholder: 'pattern',
                },
                {
                    label: 'minLength',
                    name: 'minLength',
                    component: 'inputNumber',
                    placeholder: '>=0',
                    type: 'integer',
                },
                {
                    label: 'maxLength',
                    name: 'maxLength',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                },
            ]
        }
    },
    "number": {
        label: "number",
        value: "number",
        color: 'cyan',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                },
                {
                    label: 'format',
                    name: 'format',
                    type: 'string',
                    component: 'select',
                    placeholder: 'select a value',
                    options: [
                        {
                            label: 'float',
                            value: 'float',
                        },
                        {
                            label: 'double',
                            value: 'double',
                        },
                    ]
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                },
                {
                    label: 'minimum',
                    name: 'minimum',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'maximum',
                    name: 'maximum',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'maxLength',
                    name: 'maxLength',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'multipleOf',
                    name: 'multipleOf',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'exclusiveMin',
                    name: 'exclusiveMin',
                    type: 'boolean',
                    component: 'switch',
                },
                {
                    label: 'exclusiveMax',
                    name: 'exclusiveMax',
                    type: 'boolean',
                    component: 'switch',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                },
            ]
        }
    },
    "boolean": {
        label: "boolean",
        value: "boolean",
        color: 'red',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'default',
                    name: 'default',
                    type: 'boolean',
                    component: 'select',
                    placeholder: 'select a value',
                    options: [
                        {
                            label: 'true',
                            value: true,
                        },
                        {
                            label: 'false',
                            value: false,
                        },
                    ]
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                },
            ],
        }
    },
    "integer": {
        label: "integer",
        value: "integer",
        color: 'green',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'enum',
                    name: 'enum',
                    component: 'selectTag',
                    type: 'array',
                    placeholder: '输入文本后按回车添加',
                },
                {
                    label: 'format',
                    name: 'format',
                    type: 'string',
                    component: 'select',
                    placeholder: 'select a value',
                    options: [
                        {
                            label: 'int32',
                            value: 'int32',
                        },
                        {
                            label: 'int64',
                            value: 'int64',
                        },
                    ]
                },
                {
                    label: 'default',
                    name: 'default',
                    component: 'input',
                    placeholder: 'default',
                    type: 'string',
                },
                {
                    label: 'example',
                    name: 'example',
                    type: 'string',
                    component: 'input',
                    placeholder: 'example',
                },
                {
                    label: 'minimum',
                    name: 'minimum',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'maximum',
                    name: 'maximum',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'maxLength',
                    name: 'maxLength',
                    type: 'integer',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'multipleOf',
                    name: 'multipleOf',
                    type: 'number',
                    component: 'inputNumber',
                    placeholder: '>=0',
                },
                {
                    label: 'exclusiveMin',
                    name: 'exclusiveMin',
                    type: 'boolean',
                    component: 'switch',
                },
                {
                    label: 'exclusiveMax',
                    name: 'exclusiveMax',
                    type: 'boolean',
                    component: 'switch',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                },
            ]
        }
    },
    "array": {
        label: "array",
        value: "array",
        color: 'orange',
        active: false,
        props: {
            label: 'Properties',
            options: [
                {
                    label: 'minItems',
                    name: 'minItems',
                    type: 'integer',
                    placeholder: '>=0',
                    component: 'inputNumber',
                },
                {
                    label: 'maxItems',
                    name: 'maxItems',
                    type: 'integer',
                    placeholder: '>=0',
                    component: 'inputNumber',
                },
                {
                    label: 'uniqueItems',
                    name: 'additionalProperties',
                    component: 'switch',
                    type: 'boolean',
                },
                {
                    label: 'deprecated',
                    name: 'deprecated',
                    type: 'boolean',
                    component: 'switch',
                },
            ],
        },
    },
}

export const endpointStatus = new Map([[0, '未知'], [1, '设计中'], [2, '开发中'], [3, '已发布'], [4, '已过时']])

export const serveStatus = new Map([[0, '未知'], [1, '新建'], [2, '设计中'], [3, '已发布'], [4, ' 已禁用']])

export const serveStatusTagColor = new Map([[0, 'default'], [1, 'default'], [2, 'processing'], [3, 'success'], [4, 'error']])

export const endpointStatusColor = new Map([[0, 'default'], [1, 'default'], [2, 'processing'], [3, 'success'], [4, 'error']])

export const endpointStatusOpts = [
    {
        label: "未知",
        value: 0,
    },
    {
        label: "设计中",
        value: 1,
    },
    {
        label: "开发中",
        value: 2,
    },
    {
        label: "已发布",
        value: 3,
    },
    {
        label: "已过时",
        value: 4,
    }
]

export const mediaTypesOpts = [
    {
        label: "application/EDI-X12",
        value: "application/EDI-X12",
        "disabled": true,
    },
    {
        label: "application/EDIFACT",
        value: "application/EDIFACT",
        "disabled": true,
    },
    {
        label: "application/atom+xml",
        value: "application/atom+xml",
        "disabled": true,
    },
    {
        label: "application/font-woff",
        value: "application/font-woff",
    },
    {
        label: "application/gzip",
        value: "application/gzip",
    },
    {
        label: "application/javascript",
        value: "application/javascript",
    },
    {
        label: "application/json",
        value: "application/json",
    },
    {
        label: "application/octet-stream",
        value: "application/octet-stream",
    },
    {
        label: "application/ogg",
        value: "application/ogg",
    },
    {
        label: "application/pdf",
        value: "application/pdf",
    },
    {
        label: "application/postscript",
        value: "application/postscript",
    },
    {
        label: "application/soap+xml",
        value: "application/soap+xml",
    },
    {
        label: "application/x-bittorrent",
        value: "application/x-bittorrent",
    },
    {
        label: "application/x-tex",
        value: "application/x-tex",
    },
    {
        label: "application/x-www-form-urlencoded",
        value: "application/x-www-form-urlencoded",
    },
    {
        label: "application/xhtml+xml",
        value: "application/xhtml+xml",
    },
    {
        label: "application/xml",
        value: "application/xml",
    },
    {
        label: "application/xml-ditd",
        value: "application/xml-ditd",
        "disabled": true,
    },
    {
        label: "application/xop+xml",
        value: "application/xop+xml",
    },
    {
        label: "application/zip",
        value: "application/zip",
    },
    {
        label: "multipart/form-data",
        value: "multipart/form-data",
    },
    {
        label: "text/html",
        value: "text/html",
    },
    {
        label: "text/plain",
        value: "text/plain",
    },

]

export const defaultPathParams = {
    name: '',
    type: 'string',
    description: '',
    required: false,
}

export const defaultQueryParams = {
    name: '',
    type: 'string',
    description: '',
    required: false,
}

export const defaultHeaderParams = {
    name: '',
    type: 'string',
    description: '',
    required: false,
}


export const defaultCookieParams = {
    name: '',
    type: 'string',
    description: '',
    required: false,
}

export const defaultCodeResponse = {
    "code": "",
    "endpointId": "",
    "mediaType": "",
    "description": "",
    "schemaRefId": null,
    "examples": "",
    "schemaItem": {
        "id": null,
        "name": "",
        "type": "object",
        "content": "",
        "ResponseBodyId": null
    },
    "headers": []
}

export const defaultEndpointDetail = {
    "name": "",
    "projectId": "",
    "serveId": "",
    "useId": "",
    "method": "",
    "description": "",
    "operationId": "",
    "security": "",
    "requestBody": {
        "id": null,
        "endpointId": null,
        "mediaType": "",
        "description": "",
        "schemaRefId": null,
        "examples": "",
        "schemaItem": {
            "id": null,
            "name": "",
            "type": "object",
            "content": "",
            "requestBodyId": null
        }
    },
    "responseBodies": [],
    "params": [],
    "headers": [],
    "cookies": []
}
