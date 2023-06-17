import request from '@/utils/request';
import {QueryParams} from "@/types/data";
import {Interface} from "@/views/component/debug/data";

const apiPath = 'testInterfaces';

export async function query(params: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}
export async function get(id: number): Promise<any> {
    return request({url: `/${apiPath}/${id}`});
}
export async function save(data: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: data.id? 'PUT': 'POST',
        data: data,
    });
}
export async function remove(id: number, type: string): Promise<any> {
    const params = {type}
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
        params,
    });
}
export async function move(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/move`,
        method: 'post',
        data: data,
    });
}

export async function clone(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}/clone`,
        method: 'post'
    })
}

export async function saveTestDebugData(interf: Interface): Promise<any> {
    return request({
        url: `/apiPath/saveTestDebugData`,
        method: 'post',
        data: interf,
    });
}