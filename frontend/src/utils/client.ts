import { GetBackendBasePath } from '../main'


export class Client {
    baseUrl: string
    constructor(baseUrl) {
        this.baseUrl = baseUrl
    }

    public PostJSON(endpoint: string, reqObj: any, additionalHeaders?: Headers): Promise<Response> {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
        if (additionalHeaders != undefined && additionalHeaders != null) {
            additionalHeaders.forEach((value, name) => {
                headers.append(name, value)
            })
        }
        let reqOptions: RequestInit = {
            method: "POST",
            headers: headers,
            body: JSON.stringify(reqObj),
        }

        return fetch(this.baseUrl + endpoint, reqOptions)
    }

    public Get(endpoint: string, queryParams?: any, headers?: Headers): Promise<Response> {
        let reqOptions: RequestInit = {
            method: "GET",
            headers: headers
        }
        let query = this.getQuery(queryParams)
        return fetch(this.baseUrl + endpoint + query, reqOptions)
    }

    public Delete(endpoint: string, queryParams?: any, headers?: Headers): Promise<Response> {
        let reqOptions: RequestInit = {
            method: "DELETE",
            headers: headers
        }
        let query = this.getQuery(queryParams)
        return fetch(this.baseUrl + endpoint + query, reqOptions)
    }

    public UserAuthHeader(userId: string): Headers {
        let headers = new Headers();
        headers.set("X-USER-ID", userId)
        return headers
    }

    private getQuery(params?: any): string {
        if (params == undefined || params == null)
            return ""
        return "?" + new URLSearchParams(params).toString();
    }
}

export let DefaultClient = new Client(GetBackendBasePath()); 