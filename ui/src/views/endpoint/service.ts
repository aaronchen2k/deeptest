import request from '@/utils/request';
import {QueryParams} from "@/views/project/data";

const apiPath = 'endpoints';

export async function query(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}
export async function get(id: number): Promise<any> {
    return request({url: `/${apiPath}/${id}`});
}
export async function getDetail(id: number): Promise<any> {
    const params = {
        detail: true,
    }
    return request({url: `/${apiPath}/${id}`, params});
}
export async function save(data: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: data.id? 'PUT': 'POST',
        data: data,
    });
}
export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}

/**
 * 接口列表
 * */
export async function getEndpointList(data: any): Promise<any> {
    return request({
        url: `/endpoint/index`,
        method: 'post',
        data: data
    });
}

/**
 * 接口详情
 * */
export async function getEndpointDetail(id: Number | String | any): Promise<any> {
    return request({
        url: `/endpoint/detail?id=${id}`,
        method: 'get',
    });
}

/**
 * 删除接口
 * */
export async function deleteEndpoint(id: Number): Promise<any> {
    return request({
        url: `/endpoint/delete?id=${id}`,
        method: 'delete',
    });
}


/**
 * 复制接口
 * */
export async function copyEndpoint(id: Number): Promise<any> {
    return request({
        url: `/endpoint/copy?id=${id}`,
        method: 'get',
    });
}


/**
 * 获取yaml展示
 * */
export async function getYaml(data: any): Promise<any> {
    return request({
        url: `/endpoint/yaml`,
        method: 'post',
        data: data
    });
}


/**
 * 接口过时
 * */
export async function expireEndpoint(id: Number): Promise<any> {
    return request({
        url: `/endpoint/expire?id=${id}`,
        method: 'put',
    });
}

/**
 * 保存接口
 * */
export async function saveEndpoint(data: any): Promise<any> {
    return request({
        url: `/endpoint/save`,
        method: 'post',
        data: data
    });
}



/**
 * 更新接口状态
 * */
export async function updateStatus(data: any): Promise<any> {
    return request({
        url: `/endpoint/updateStatus?id=${data.id}&status=${data.status}`,
        method: 'put',
    });
}
