export var CONSTANT: CONSTANT_INTERFACE = {
  _SERVICE_URL_DEV: 'http://localhost:8080/',
  _SERVICE_URL_PRODUCTION: 'http://172.16.90.66:8080/',

  SERVICE_URL: undefined,
  API_URL: undefined,

  API_PATH: 'api/client/v1/',
  UPLOAD_URI: 'api/client/v1/uploadSingle',

  TOKEN: undefined,

  PROFILE: undefined,

  ALL_ORGS: undefined,
  CURR_ORG_ID: undefined,

  RECENT_PROJECTS: undefined,
  CURRENT_PROJECT: {id: null, name: null},

  CUSTOM_FIELD_FOR_PROJECT: undefined,
  CASE_PROPERTY_MAP: undefined,

  TOKEN_KEY: 'com.ngtesting.token',
  TOKEN_EXPIRE: 'com.ngtesting.expire',

  ExeStatus: {'not_start': '未开始', 'in_progress': '执行中', 'end': '已结束'},
  EntityDisabled: {'false': '启用', 'true': '归档', '': '所有'},

  ScreenSize: {h: 0, w: 0},
  DebounceTime: 500,

  STATE_CHANGE_PROFILE: 'profile.refresh',
  STATE_CHANGE_ORGS: 'my.orgs.change',
  STATE_CHANGE_PROJECTS: 'recent.projects.change'
};

export interface CURRENT_PROJECT_INTERFACE {
  id: number,
  name: string
}

export interface CONSTANT_INTERFACE {
  _SERVICE_URL_DEV: string,
  _SERVICE_URL_PRODUCTION: string,

  SERVICE_URL: string,
  API_URL: string,

  API_PATH: string,
  UPLOAD_URI: string,

  TOKEN: string,

  PROFILE: any,

  ALL_ORGS: any[],
  CURR_ORG_ID: number,

  RECENT_PROJECTS: any[],
  CURRENT_PROJECT: CURRENT_PROJECT_INTERFACE,

  CUSTOM_FIELD_FOR_PROJECT: any,
  CASE_PROPERTY_MAP: any,

  TOKEN_KEY: string,
  TOKEN_EXPIRE: string,

  ExeStatus: any,
  EntityDisabled: any,

  ScreenSize: any,
  DebounceTime: number,

  STATE_CHANGE_PROFILE: string,
  STATE_CHANGE_ORGS: string,
  STATE_CHANGE_PROJECTS: string

}

