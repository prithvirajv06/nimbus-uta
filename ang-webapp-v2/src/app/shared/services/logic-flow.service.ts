import { Injectable } from "@angular/core";
import { CommonService } from "./common.service";
import { LogicFlow } from "../types/logic-flow.type";
import { Observable, catchError } from "rxjs";
import { ApiResponse } from "../types/common.type";


@Injectable({
    providedIn: 'root'
})
export class LogicFlowService extends CommonService {

    getList(data: any): Observable<ApiResponse<LogicFlow[]>> {
        return this.httpClient.post<ApiResponse<LogicFlow[]>>(`${this.baseApiUrl}/logic-flows/list`, data)
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    get(id: string, version: number) {
        return this.httpClient.get<ApiResponse<LogicFlow>>(`${this.baseApiUrl}/logic-flow`, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    create(data: any) {
        return this.httpClient.post<ApiResponse<LogicFlow>>(`${this.baseApiUrl}/logic-flows`, data)
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    update(id: string, version: number, data: any) {
        return this.httpClient.put<ApiResponse<LogicFlow>>(`${this.baseApiUrl}/logic-flows`, data, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    delete(id: string, version: number) {
        return this.httpClient.delete<ApiResponse<LogicFlow>>(`${this.baseApiUrl}/logic-flow`, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

}