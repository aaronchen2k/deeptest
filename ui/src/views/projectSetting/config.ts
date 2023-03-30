import { Schema } from "./data";

const schemaList: Schema[] = [
    {
        type: 'tooltip',
        text: '新建服务',
        title: '一个产品服务端通常对应一个或多个服务(微服务)，服务可以有多个版本并行，新的服务默认起始版本为v0.1.0。'
    },
    {
        type: 'input',
        stateName: 'name',
        placeholder: '服务名称',
        valueType: 'string'
    },
    {
        type: 'select',
        stateName: 'serveId',
        placeholder: '负责人(默认创建人)',
        options: [],
        valueType: 'string',
        mode: 'combobox',
        focusType: 'userListOptions'
    },
    {
        type: 'input',
        stateName: 'description',
        placeholder: '输入描述',
        valueType: 'string'
    },
    {
        type: 'button',
        text: '确定',
    },
];

const globalVarsColumns = [
    {
        title: '变量名',
        dataIndex: 'name',
        key: 'name',
        slots: { customRender: 'customName' },
        type: 'input',
        placeholder: ''
    },
    {
        title: '远程值',
        dataIndex: 'remoteValue',
        key: 'remoteValue',
        slots: { customRender: 'customRemoteValue' },
        type: 'input',
        placeholder: ''
    },
    {
        title: '本地值',
        dataIndex: 'localValue',
        key: 'localValue',
        slots: { customRender: 'customLocalValue' },
        type: 'input',
        placeholder: ''
    },
    {
        title: '说明',
        key: 'description',
        dataIndex: 'description',
        slots: { customRender: 'customDescription' },
        type: 'input',
        placeholder: ''
    },
    {
        title: '操作',
        key: 'action',
        slots: { customRender: 'customAction' },
        type: 'button',
        placeholder: ''
    },
];

const globalParamscolumns: any = [
    {
        title: '参数名',
        dataIndex: 'name',
        key: 'name',
        slots: { customRender: 'customName' },
    },
    {
        title: '类型',
        dataIndex: 'type',
        key: 'type',
        slots: { customRender: 'customType' },
    },
    {
        title: '必须',
        dataIndex: 'required',
        key: 'required',
        slots: { customRender: 'customRequired' },
    },
    {
        title: '默认值',
        key: 'defaultValue',
        dataIndex: 'defaultValue',
        slots: { customRender: 'customDefaultValue' },
    },
    // {
    //   title: '默认启用',
    //   key: 'description',
    //   dataIndex: 'description',
    // },
    {
        title: '说明',
        key: 'description',
        dataIndex: 'description',
        slots: { customRender: 'customDescription' },
    },
    {
        title: '操作',
        key: 'action',
        slots: { customRender: 'customAction' },
    },
];

const serveServersColumns: any = [
    {
        title: '服务名',
        dataIndex: 'serveName',
        key: 'serveName',
        slots: { customRender: 'customName' },
    },
    {
        title: '前置 URL ',
        dataIndex: 'url',
        key: 'url',
        slots: { customRender: 'customUrl' },
    },
];

// 全局参数tab切换列表
const tabPaneList = [{
    name: 'header',
    type: 'header'
}, {
    name: 'cookie',
    type: 'cookie'
}, {
    name: 'query',
    type: 'query'
}, {
    name: 'body',
    type: 'body'
}];


const serviceColumns = [
    {
        title: '服务名称',
        dataIndex: 'name',
        slots: { customRender: 'name', title: 'fdshfh' },
    },
    {
        title: '描述',
        dataIndex: 'description',
        slots: { customRender: 'description' },
    },
    {
        title: '关联服务',
        dataIndex: 'servers',
        slots: { customRender: 'customServers' },
    },
    {
        title: '状态',
        dataIndex: 'statusDesc',
        slots: { customRender: 'customStatus' },
    },
    {
        title: '创建人',
        dataIndex: 'createUser',
    },
    {
        title: '创建时间',
        dataIndex: 'createdAt',
    },
    {
        title: '最近更新时间',
        dataIndex: 'updatedAt',
    },
    {
        title: '操作',
        dataIndex: 'operation',
        slots: { customRender: 'operation' },
    },
];

const schemaColumns = [
    {
        title: '组件名称',
        dataIndex: 'name',
        width: '30%',
        slots: { customRender: 'name' },
    },
    {
        title: '标签',
        dataIndex: 'tags',
    },
    {
        title: '操作',
        dataIndex: 'operation',
        slots: { customRender: 'operation' },
    },
];

const versionColumns = [
    {
      title: '版本号',
      dataIndex: 'value',
      width: '30%',
      slots: { customRender: 'value' },
    },
    {
      title: '负责人',
      dataIndex: 'createUser',
    },
    {
      title: '描述',
      dataIndex: 'description',
    },
    {
      title: '操作',
      dataIndex: 'operation',
      slots: { customRender: 'operation' },
    },
  ];


export {
    globalParamscolumns,
    globalVarsColumns,
    serveServersColumns,
    tabPaneList,
    serviceColumns,
    schemaColumns,
    versionColumns
}